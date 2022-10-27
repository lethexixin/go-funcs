package producer

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

import (
	"github.com/lethexixin/go-funcs/common/logger"
)

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"
)

type KafkaProducer struct {
	Producer *kafka.Producer
}

var counterMetric *prometheus.CounterVec

func InitMetrics(appName string) {
	once := sync.Once{}
	once.Do(func() {
		counterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_kafka_producer_metric_total", strings.ReplaceAll(appName, "-", "_")),
			Help: fmt.Sprintf("kafka producer number of (topic,flag) for %s", strings.ReplaceAll(appName, "-", "_")),
		}, []string{"topic", "flag"})
		prometheus.MustRegister(counterMetric)
	})
}

type Options struct {
	bootstrapServers           string
	messageMaxBytes            int
	batchSize                  int
	lingerMs                   int
	stickyPartitioningLingerMs int
	retries                    int
	retryBackoffMs             int
	acks                       string
	compressionType            string
	securityProtocol           string
	saslMechanism              string
	saslUsername               string
	saslPassword               string

	signChan chan os.Signal
}

type MsgProducerData struct {
	Data  string
	Key   string
	Topic string
}

type Option func(*Options)

const (
	// kafka producer config refer to https://help.aliyun.com/document_detail/68165.html and https://www.lovergamer.com/posts/kafka/kafka_products

	DefaultBootstrapServers = "127.0.0.1:9092"
	// DefaultMessageMaxBytes 允许的最大记录批量大小
	DefaultMessageMaxBytes = 1048588
	// DefaultBatchSize 发往每个分区(Partition)的消息缓存量(消息内容的字节数之和, 不是条数)
	DefaultBatchSize = 16384
	// DefaultLingerMs 每条消息在缓存中的最长时间
	DefaultLingerMs = 1000
	// DefaultStickyPartitioningLingerMs 黏性分区策略每条消息在缓存中的最长时间
	DefaultStickyPartitioningLingerMs = 1000
	// DefaultRetries 重试次数
	DefaultRetries = 3
	// DefaultRetryBackoffMs 重试间隔
	DefaultRetryBackoffMs = 1000
	// DefaultAcks 确认机制
	DefaultAcks = "1"
	// DefaultCompressionType 压缩方式
	DefaultCompressionType  = "snappy"
	DefaultSecurityProtocol = "PLAINTEXT"
	DefaultSaslMechanism    = "PLAIN"
	DefaultSaslUsername     = "kafka"
	DefaultSaslPassword     = "123456"
)

func BootstrapServers(bootstrapServers string) Option {
	return func(o *Options) {
		o.bootstrapServers = bootstrapServers
	}
}

func MessageMaxBytes(messageMaxBytes int) Option {
	return func(o *Options) {
		o.messageMaxBytes = messageMaxBytes
	}
}

func BatchSize(batchSize int) Option {
	return func(o *Options) {
		o.batchSize = batchSize
	}
}

func LingerMs(lingerMs int) Option {
	return func(o *Options) {
		o.lingerMs = lingerMs
	}
}

func StickyPartitioningLingerMs(stickyPartitioningLingerMs int) Option {
	return func(o *Options) {
		o.stickyPartitioningLingerMs = stickyPartitioningLingerMs
	}
}

func Retries(retries int) Option {
	return func(o *Options) {
		o.retries = retries
	}
}

func RetryBackoffMs(retryBackoffMs int) Option {
	return func(o *Options) {
		o.retryBackoffMs = retryBackoffMs
	}
}

func Acks(acks string) Option {
	return func(o *Options) {
		o.acks = acks
	}
}

func CompressionType(compressionType string) Option {
	return func(o *Options) {
		o.compressionType = compressionType
	}
}

func SecurityProtocol(securityProtocol string) Option {
	return func(o *Options) {
		o.securityProtocol = securityProtocol
	}
}

func SaslMechanism(saslMechanism string) Option {
	return func(o *Options) {
		o.saslMechanism = saslMechanism
	}
}

func SaslUsername(saslUsername string) Option {
	return func(o *Options) {
		o.saslUsername = saslUsername
	}
}

func SaslPassword(saslPassword string) Option {
	return func(o *Options) {
		o.saslPassword = saslPassword
	}
}

func SignChan(signChan chan os.Signal) Option {
	return func(o *Options) {
		o.signChan = signChan
	}
}

//InitProducer
//
//1. confluent-kafka-go build refer to docs/kafka.md
//
//2. (fn func()) examples:
//
// type Handler struct {
// 	 sourceChan chan MsgProducerData
// }
//
// func (h *Handler) producerFunc(p *KafkaProducer) func() {
//	 return func() {
//		 select {
//		 case v := <-h.sourceChan:
//           go func() {
//               _ = p.SendMsg(&v)
//           }()
//		 }
//	 }
// }
//
func (k *KafkaProducer) InitProducer(fn func(), options ...Option) (err error) {
	opts := Options{
		bootstrapServers:           DefaultBootstrapServers,
		messageMaxBytes:            DefaultMessageMaxBytes,
		batchSize:                  DefaultBatchSize,
		lingerMs:                   DefaultLingerMs,
		stickyPartitioningLingerMs: DefaultStickyPartitioningLingerMs,
		retries:                    DefaultRetries,
		retryBackoffMs:             DefaultRetryBackoffMs,
		acks:                       DefaultAcks,
		compressionType:            DefaultCompressionType,
		securityProtocol:           DefaultSecurityProtocol,
		saslMechanism:              DefaultSaslMechanism,
		saslUsername:               DefaultSaslUsername,
		saslPassword:               DefaultSaslPassword,
	}

	for _, o := range options {
		o(&opts)
	}

	if opts.signChan == nil {
		return errors.New("opts.signChan == nil")
	}

	kafkaConf := &kafka.ConfigMap{
		"bootstrap.servers":             opts.bootstrapServers,
		"api.version.request":           "true",
		"message.max.bytes":             opts.messageMaxBytes,
		"batch.size":                    opts.batchSize,
		"linger.ms":                     opts.lingerMs,
		"sticky.partitioning.linger.ms": opts.stickyPartitioningLingerMs,
		"retries":                       opts.retries,
		"retry.backoff.ms":              opts.retryBackoffMs,
		"acks":                          opts.acks,
		"compression.type":              opts.compressionType}

	switch strings.ToUpper(opts.securityProtocol) {
	case "PLAINTEXT":
		_ = kafkaConf.SetKey("security.protocol", "plaintext")
	case "SASL_SSL":
		_ = kafkaConf.SetKey("security.protocol", "sasl_ssl")
		_ = kafkaConf.SetKey("ssl.ca.location", "conf/mix-4096-ca-cert")
		_ = kafkaConf.SetKey("sasl.username", opts.saslUsername)
		_ = kafkaConf.SetKey("sasl.password", opts.saslPassword)
		_ = kafkaConf.SetKey("sasl.mechanism", opts.saslMechanism)
		_ = kafkaConf.SetKey("enable.ssl.certificate.verification", "false")
	case "SASL_PLAINTEXT":
		_ = kafkaConf.SetKey("security.protocol", "sasl_plaintext")
		_ = kafkaConf.SetKey("sasl.username", opts.saslUsername)
		_ = kafkaConf.SetKey("sasl.password", opts.saslPassword)
		_ = kafkaConf.SetKey("sasl.mechanism", opts.saslMechanism)
	default:
		return fmt.Errorf("unknown kafka protocol:%s", opts.securityProtocol)
	}

	k.Producer, err = kafka.NewProducer(kafkaConf)
	if err != nil {
		logger.Errorf("failed to create kafka producer:%s", err.Error())
		return err
	}
	defer k.Producer.Close()

	logger.Info("create kafka producer successful")

	// Listen to all the events on the default events channel
	go func() {
		for e := range k.Producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				// The message delivery report, indicating success or permanent failure after retries have been exhausted.
				// Application level retries won't help since the client is already configured to do that.
				if ev.TopicPartition.Error != nil {
					logger.Errorf("delivery failed:%s", ev.TopicPartition.Error.Error())
					if counterMetric != nil {
						counterMetric.With(prometheus.Labels{"topic": *ev.TopicPartition.Topic, "flag": "error"}).Inc()
					}
				} else {
					if counterMetric != nil {
						counterMetric.With(prometheus.Labels{"topic": *ev.TopicPartition.Topic, "flag": "success"}).Inc()
					}
				}
			case kafka.Error:
				// Generic client instance-level errors, such as broker connection failures, authentication issues, etc.
				// These errors should generally be considered informational as the underlying client will automatically try to recover from any errors encountered, the application does not need to take action on them.
				logger.Errorf("p.Events is err:%s", ev.Error())
				if counterMetric != nil {
					counterMetric.With(prometheus.Labels{"topic": "producer", "flag": "error"}).Inc()
				}
			default:
				logger.Infof("ignored event:%v", ev)
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	run := true
	for run {
		select {
		case <-opts.signChan:
			logger.Info("producer get <-opts.signChan")
			run = false
		default:
			fn()
		}
	}

	// Wait for message deliveries before shutting down
	k.Producer.Flush(15 * 1000)

	return nil
}

func (k *KafkaProducer) SendMsg(v *MsgProducerData) (err error) {
	if len(v.Key) == 0 {
		v.Key = strconv.Itoa(int(time.Now().UnixNano()))
	}
	if err = k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &v.Topic, Partition: kafka.PartitionAny},
		Key:            []byte(v.Key),
		Value:          []byte(v.Data),
	}, nil); err != nil {
		switch err.(kafka.Error).Code() {
		case kafka.ErrQueueFull:
			// Producer queue is full, wait 1s for messages to be delivered then try again.
			// time.Sleep(time.Second)
			logger.Errorf("failed to produce message, err:%s", err.Error())
		case kafka.ErrMsgSizeTooLarge:
			logger.Errorf("failed to produce message, err:%s, len:%d", err.Error(), len(v.Data))
		default:
			logger.Errorf("failed to produce message:%s, err:%s", v.Data, err.Error())
		}
		if counterMetric != nil {
			counterMetric.With(prometheus.Labels{"topic": v.Topic, "flag": "sink_err"}).Inc()
		}
		return err
	}
	if counterMetric != nil {
		counterMetric.With(prometheus.Labels{"topic": v.Topic, "flag": "sink_ok"}).Inc()
	}
	return nil
}
