package consumer

import (
	"errors"
	"fmt"
	"net"
	"time"
)

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

import (
	"github.com/lethexixin/go-funcs/common/logger"
)

type MQConsumer struct {
	Consumer *consumer
}

type Options struct {
	addr        string // MQ地址
	virtualHost string // 虚拟主机名称
	queueName   string // 队列名称
	tagConsumer string // 消费者标签
}

type Option func(*Options)

var (
	DefaultAddr        = "amqp://root:123456@127.0.0.1:5672"
	DefaultVirtualHost = "rmq"
	DefaultQueueName   = "rmq-queue"
	DefaultTagConsumer = "rmq-tag"
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

func TagConsumer(tagConsumer string) Option {
	return func(o *Options) {
		o.tagConsumer = tagConsumer
	}
}

// Consumer holds all information
// about the RabbitMQ connection
// This setup does limit a consumer
// to one exchange. This should not be
// an issue. Having to connect to multiple
// exchanges means something else is
// structured improperly.
type consumer struct {
	opts    *Options
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan error
}

// NewConsumer returns a Consumer struct
// that has been initialized properly
// essentially don't touch conn, channel, or
// done and you can create Consumer manually
func (c *MQConsumer) Init(options ...Option) {
	opts := Options{
		addr:        DefaultAddr,
		virtualHost: DefaultVirtualHost,
		queueName:   DefaultQueueName,
		tagConsumer: DefaultTagConsumer,
	}

	for _, o := range options {
		o(&opts)
	}

	c.Consumer = &consumer{
		opts: &opts,
		done: make(chan error),
	}
}

// ReConnect is called in places where NotifyClose() channel is called
// wait 30 seconds before trying to reconnect. Any shorter amount of time
// will  likely destroy the error log while waiting for servers to come
// back online. This requires two parameters which is just to satisfy
// the AnnounceQueue call and allows greater flexibility
func (c *consumer) reConnect(queueName string, ack bool) (<-chan amqp.Delivery, error) {
	time.Sleep(30 * time.Second)

	if err := c.Connect(); err != nil {
		logger.Errorf("connect %s err:%s", queueName, err.Error())
	}
	deliveries, err := c.AnnounceQueue(queueName, ack)
	if err != nil {
		return deliveries, errors.New("couldn't connect")
	}
	return deliveries, nil
}

// Connect to RabbitMQ server
func (c *consumer) Connect() (err error) {
	logger.Infof("dialing rmq %s...", c.opts.addr)
	c.conn, err = amqp.DialConfig(c.opts.addr, amqp.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 2*time.Second)
		},
		Vhost: c.opts.virtualHost,
	})
	if err != nil {
		return fmt.Errorf("dial err: %s", err.Error())
	}

	go func() {
		// Waits here for the channel to be closed
		logger.Infof("closing %s:%s", c.opts.addr, <-c.conn.NotifyClose(make(chan *amqp.Error)))
		// Let Handle know it's not time to reconnect
		c.done <- errors.New("channel Closed")
	}()

	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel err: %s", err.Error())
	}
	logger.Infof("rmq %s connected successful!", c.opts.queueName)
	return nil
}

// AnnounceQueue sets the queue that will be listened to for this connection...
func (c *consumer) AnnounceQueue(queueName string, ack bool) (<-chan amqp.Delivery, error) {
	// Qos determines the amount of messages that the queue will pass to you before
	// it waits for you to ack them. This will slow down queue consumption but
	// give you more certainty that all messages are being processed. As load increases
	// I would recommend upping the about of Threads and Processors the go process
	// uses before changing this although you will eventually need to reach some
	// balance between threads, process, and Qos.
	err := c.channel.Qos(50, 0, false)
	if err != nil {
		return nil, fmt.Errorf("error setting qos: %s", err.Error())
	}

	logger.Infof("starting consumer... (queue:%s, tag:%s)", queueName, c.opts.tagConsumer)
	deliveries, err := c.channel.Consume(
		queueName,          // name
		c.opts.tagConsumer, // tagConsumer,
		ack,                // autoAck
		false,              // exclusive
		false,              // noLocal
		false,              // noWait
		nil,                // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue consume err: %s", err.Error())
	}

	return deliveries, nil
}

// Handle has all the logic to make sure your program keeps running
// d should be a delivery channel as created when you call AnnounceQueue
// fn should be a function that handles the processing of deliveries
// this should be the last thing called in main as code under it will
// become unreachable unless put int a goroutine. The q and rk params
// are redundant but allow you to have multiple queue listeners in main
// without them you would be tied into only using one queue per connection
func (c *consumer) Handle(delivery <-chan amqp.Delivery, fn func(<-chan amqp.Delivery), threads int, queue string, ack bool) {
	var err error
	for {
		for i := 0; i < threads; i++ {
			go fn(delivery)
		}

		// Go into reconnect loop when
		// c.done is passed non nil values
		if <-c.done != nil {
			delivery, err = c.reConnect(queue, ack)
			if err != nil {
				// Very likely chance of failing
				// should not cause worker to terminate
				logger.Errorf("%s reconnecting err:%s", queue, err.Error())
			}
		}
		logger.Infof("%s reconnected successful...", queue)
	}
}
