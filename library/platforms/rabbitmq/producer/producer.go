package producer

import (
	"errors"
	"net"
	"time"
)

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

import (
	"github.com/lethexixin/go-funcs/common/logger"
)

type MQProducer struct {
	Producer *producer
}

type producer struct {
	opts            *Options
	connection      *amqp.Connection
	channel         *amqp.Channel
	done            chan bool
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	isReady         bool
}

const (
	// When reconnecting to the server after connection failure
	reconnectDelay = 5 * time.Second

	// When setting up the channel after a channel exception
	reInitDelay = 2 * time.Second

	// When resending messages the server didn't confirm
	resendDelay = 5 * time.Second
)

var (
	errNotConnected  = errors.New("not connected to a server")
	errAlreadyClosed = errors.New("already closed: not connected to the server")
	errShutdown      = errors.New("producer is shutting down")
)

type Options struct {
	addr         string                 // MQ地址
	virtualHost  string                 // 虚拟主机名称
	queueName    string                 // 队列名称
	routingKey   string                 // 队列路由键
	exchangeName string                 // 交换机名称
	exchangeType string                 // 交换机类型
	args         map[string]interface{} // 优先级队列
}

type Option func(*Options)

const (
	DefaultAddr         = "amqp://root:123456@127.0.0.1:5672"
	DefaultVirtualHost  = "rmq"
	DefaultQueueName    = "rmq-queue"
	DefaultRoutingKey   = "rmq-key"
	DefaultExchangeName = "rmq-ex"
	DefaultExchangeType = "direct"
)

func Addr(addr string) Option {
	return func(o *Options) {
		o.addr = addr
	}
}

func VirtualHost(virtualHost string) Option {
	return func(o *Options) {
		o.virtualHost = virtualHost
	}
}

func QueueName(queueName string) Option {
	return func(o *Options) {
		o.queueName = queueName
	}
}

func RoutingKey(routingKey string) Option {
	return func(o *Options) {
		o.routingKey = routingKey
	}
}

func ExchangeName(exchangeName string) Option {
	return func(o *Options) {
		o.exchangeName = exchangeName
	}
}

func ExchangeType(exchangeType string) Option {
	return func(o *Options) {
		o.exchangeType = exchangeType
	}
}

func Args(args map[string]interface{}) Option {
	return func(o *Options) {
		o.args = args
	}
}

func (p *producer) isExchangeProducer() bool {
	if len(p.opts.routingKey) > 0 && len(p.opts.exchangeName) > 0 && len(p.opts.exchangeType) > 0 {
		return true
	}
	return false
}

// Init creates a new consumer state instance, and automatically attempts to connect to the server.
func (p *MQProducer) Init(options ...Option) {
	opts := Options{
		addr:         DefaultAddr,
		virtualHost:  DefaultVirtualHost,
		queueName:    DefaultQueueName,
		routingKey:   DefaultRoutingKey,
		exchangeName: DefaultExchangeName,
		exchangeType: DefaultExchangeType,
	}

	for _, o := range options {
		o(&opts)
	}

	p.Producer = &producer{
		opts: &opts,
		done: make(chan bool),
	}

	go p.Producer.handleReconnect()
}

// handleReconnect will wait for a connection error on notifyConnClose, and then continuously attempt to reconnect.
func (p *producer) handleReconnect() {
	for {
		p.isReady = false
		logger.Infof("attempting to connect %s...", p.opts.queueName)

		conn, err := p.connect()
		if err != nil {
			logger.Errorf("failed to connect %s, err:%s. retrying...", p.opts.queueName, err.Error())
			select {
			case <-p.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}
		if done := p.handleReInit(conn); done {
			break
		}
	}
}

// connect will create a new AMQP connection
func (p *producer) connect() (*amqp.Connection, error) {
	conn, err := amqp.DialConfig(p.opts.addr, amqp.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 2*time.Second)
		},
		Vhost: p.opts.virtualHost,
	})
	if err != nil {
		logger.Errorf("connect %s is failed, err:%s", p.opts.addr, err.Error())
		return nil, err
	}

	p.changeConnection(conn)
	logger.Infof("rmq connected %s successful!", p.opts.queueName)
	return conn, nil
}

// handleReInit will wait for a channel error
// and then continuously attempt to re-initialize both channels
func (p *producer) handleReInit(conn *amqp.Connection) bool {
	for {
		p.isReady = false
		err := p.init(conn)
		if err != nil {
			logger.Errorf("failed to init channel %s, err:%. retrying...", p.opts.queueName, err.Error())
			select {
			case <-p.done:
				return true
			case <-time.After(reInitDelay):
			}
			continue
		}

		select {
		case <-p.done:
			return true
		case <-p.notifyConnClose:
			logger.Info("connection %s closed. reconnecting...", p.opts.addr)
			return false
		case <-p.notifyChanClose:
			logger.Info("channel %s closed. re-running init...", p.opts.queueName)
		}
	}
}

// init will initialize channel & declare queue
func (p *producer) init(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	if err = ch.Confirm(false); err != nil {
		return err
	}

	// 用于检查队列是否存在,已经存在不需要重复声明
	if _, err = ch.QueueDeclarePassive(
		p.opts.queueName, true, false, false, true, nil); err != nil {
		// 队列不存在,声明队列
		if _, err = ch.QueueDeclare(
			p.opts.queueName, // 队列名称
			true,             // 是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能
			false,            // 是否自动删除
			false,            // 是否设置排他
			false,            // 是否非阻塞,true为是,不等待RMQ返回信息
			p.opts.args,      // 参数,也可以传入nil
		); err != nil {
			logger.Errorf("registration queue:%s failed, err:%s", p.opts.queueName, err.Error())
			return err
		}
	}

	if p.isExchangeProducer() {
		// 队列绑定
		if err = ch.QueueBind(
			p.opts.queueName, p.opts.routingKey, p.opts.exchangeName, true, nil); err != nil {
			logger.Errorf("bind queue:%s failed, err:%s", p.opts.queueName, err.Error())
			return err
		}

		// 用于检查交换机是否存在,已经存在不需要重复声明
		if err = ch.ExchangeDeclarePassive(
			p.opts.exchangeName, p.opts.exchangeType, true, false, false, true, nil); err != nil {
			// 交换机不存在, 注册交换机
			if err = ch.ExchangeDeclare(
				p.opts.exchangeName, //交换机名称
				p.opts.exchangeType, //交换机类型
				true,                //是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能
				false,               //是否自动删除
				false,               // 是否为内部
				true,                // 是否非阻塞, true为是,不等待RMQ返回信息
				p.opts.args,         // 参数,传nil即可
			); err != nil {
				logger.Errorf("registration exchange:%s failed, err:%s", p.opts.exchangeName, err.Error())
				return err
			}
		}
	}

	p.changeChannel(ch)
	p.isReady = true
	logger.Infof("rmq %s setup successful!", p.opts.queueName)
	return nil
}

// changeConnection takes a new connection to the queue,
// and updates the close listener to reflect this.
func (p *producer) changeConnection(connection *amqp.Connection) {
	p.connection = connection
	p.notifyConnClose = make(chan *amqp.Error)
	p.connection.NotifyClose(p.notifyConnClose)
}

// changeChannel takes a new channel to the queue,
// and updates the channel listeners to reflect this.
func (p *producer) changeChannel(channel *amqp.Channel) {
	p.channel = channel
	p.notifyChanClose = make(chan *amqp.Error)
	p.notifyConfirm = make(chan amqp.Confirmation, 1)
	p.channel.NotifyClose(p.notifyChanClose)
	p.channel.NotifyPublish(p.notifyConfirm)
}

// Push will push data onto the queue, and wait for confirmation.
// If no confirms are received until within the resendTimeout,
// it continuously re-sends messages until a confirmation is received.
// This will block until the server sends a confirmation. Errors are
// only returned if the push action itself fails, see UnsafePush.
func (p *producer) Push(data []byte, priority uint8) error {
	if !p.isReady {
		return errors.New("failed to push data: not connected")
	}
	for {
		err := p.UnsafePush(data, priority)
		if err != nil {
			logger.Errorf("push data failed, err:%s. retrying...", err.Error())
			select {
			case <-p.done:
				return errShutdown
			case <-time.After(resendDelay):
			}
			continue
		}
		select {
		case confirm := <-p.notifyConfirm:
			if confirm.Ack {
				return nil
			}
		case <-time.After(resendDelay):
		}
		logger.Info("push didn't confirm. retrying...")
	}
}

// UnsafePush will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// receive the message.
func (p *producer) UnsafePush(data []byte, priority uint8) error {
	if !p.isReady {
		return errNotConnected
	}
	return p.channel.Publish(
		p.opts.exchangeName, // Exchange
		p.opts.routingKey,   // Routing key
		false,               // Mandatory
		false,               // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
			Priority:    priority,
		},
	)
}

// Stream will continuously put queue items on the channel.
// It is required to call delivery.Ack when it has been
// successfully processed, or delivery.Nack when it fails.
// Ignoring this will cause data to build up on the server.
func (p *producer) Stream() (<-chan amqp.Delivery, error) {
	if !p.isReady {
		return nil, errNotConnected
	}
	return p.channel.Consume(
		p.opts.queueName,
		"",    // Consumer
		false, // Auto-Ack
		false, // Exclusive
		false, // No-local
		false, // No-Wait
		nil,   // Args
	)
}

// Close will cleanly shut down the channel and connection.
func (p *producer) Close() error {
	if !p.isReady {
		return errAlreadyClosed
	}
	err := p.channel.Close()
	if err != nil {
		return err
	}
	err = p.connection.Close()
	if err != nil {
		return err
	}
	close(p.done)
	p.isReady = false
	return nil
}
