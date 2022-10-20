package producer

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lethexixin/go-funcs/common/logger"

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
	BootstrapServers           string
	MessageMaxBytes            int
	BatchSize                  int
	LingerMs                   int
	StickyPartitioningLingerMs int
	Retries                    int
	RetryBackoffMs             int
	Acks                       string
	CompressionType            string
	SecurityProtocol           string
	SaslMechanism              string
	SaslUsername               string
	SaslPassword               string

	SignChan chan os.Signal
}

type MsgProducerData struct {
	Data  string
	Topic string
}

type Option func(*Options)

var (
	// kafka producer config refer to https://help.aliyun.com/document_detail/68165.html and https://www.lovergamer.com/posts/kafka/kafka_products

	DefaultBootstrapServers = "127.0.0.1:9092"
	// DefaultMessageMaxBytes 允许的最大记录批量大小
	DefaultMessageMaxBytes = 1048588
	// DefaultBatchSize 发往每个分区（Partition）的消息缓存量（消息内容的字节数之和，不是条数）
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
		o.BootstrapServers = bootstrapServers
	}
}

func MessageMaxBytes(messageMaxBytes int) Option {
	return func(o *Options) {
		o.MessageMaxBytes = messageMaxBytes
	}
}

func BatchSize(batchSize int) Option {
	return func(o *Options) {
		o.BatchSize = batchSize
	}
}

func LingerMs(lingerMs int) Option {
	return func(o *Options) {
		o.LingerMs = lingerMs
	}
}

func StickyPartitioningLingerMs(stickyPartitioningLingerMs int) Option {
	return func(o *Options) {
		o.StickyPartitioningLingerMs = stickyPartitioningLingerMs
	}
}

func Retries(retries int) Option {
	return func(o *Options) {
		o.Retries = retries
	}
}

func RetryBackoffMs(retryBackoffMs int) Option {
	return func(o *Options) {
		o.RetryBackoffMs = retryBackoffMs
	}
}

func Acks(acks string) Option {
	return func(o *Options) {
		o.Acks = acks
	}
}

func CompressionType(compressionType string) Option {
	return func(o *Options) {
		o.CompressionType = compressionType
	}
}

func SecurityProtocol(securityProtocol string) Option {
	return func(o *Options) {
		o.SecurityProtocol = securityProtocol
	}
}

func SaslMechanism(saslMechanism string) Option {
	return func(o *Options) {
		o.SaslMechanism = saslMechanism
	}
}

func SaslUsername(saslUsername string) Option {
	return func(o *Options) {
		o.SaslUsername = saslUsername
	}
}

func SaslPassword(saslPassword string) Option {
	return func(o *Options) {
		o.SaslPassword = saslPassword
	}
}

func SignChan(signChan chan os.Signal) Option {
	return func(o *Options) {
		o.SignChan = signChan
	}
}

//InitProducer
//
//1. confluent-kafka-go build refer to docs/kafka.md
//
//2. (fn func()) examples:
//
// type Handler struct {
// 	 source1Chan chan MsgProducerData
//	 source2Chan chan MsgProducerData
//	 source3Chan chan MsgProducerData
// }
//
// func (h *Handler) producerFunc(p *KafkaProducer) func() {
//	 return func() {
//		 select {
//		 case v := <-h.source1Chan:
//			 _ = p.SendMsg(&v)
//		 case v := <-h.source2Chan:
//			 _ = p.SendMsg(&v)
//		 case v := <-h.source3Chan:
//			 _ = p.SendMsg(&v)
//		 }
//	 }
// }
//
func (k *KafkaProducer) InitProducer(fn func(), options ...Option) (err error) {
	opts := Options{
		BootstrapServers:           DefaultBootstrapServers,
		MessageMaxBytes:            DefaultMessageMaxBytes,
		BatchSize:                  DefaultBatchSize,
		LingerMs:                   DefaultLingerMs,
		StickyPartitioningLingerMs: DefaultStickyPartitioningLingerMs,
		Retries:                    DefaultRetries,
		RetryBackoffMs:             DefaultRetryBackoffMs,
		Acks:                       DefaultAcks,
		CompressionType:            DefaultCompressionType,
		SecurityProtocol:           DefaultSecurityProtocol,
		SaslMechanism:              DefaultSaslMechanism,
		SaslUsername:               DefaultSaslUsername,
		SaslPassword:               DefaultSaslPassword,
	}

	for _, o := range options {
		o(&opts)
	}

	if opts.SignChan == nil {
		return errors.New("opts.SignChan == nil")
	}

	kafkaConf := &kafka.ConfigMap{
		"bootstrap.servers":             opts.BootstrapServers,
		"api.version.request":           "true",
		"message.max.bytes":             opts.MessageMaxBytes,
		"batch.size":                    opts.BatchSize,
		"linger.ms":                     opts.LingerMs,
		"sticky.partitioning.linger.ms": opts.StickyPartitioningLingerMs,
		"retries":                       opts.Retries,
		"retry.backoff.ms":              opts.RetryBackoffMs,
		"acks":                          opts.Acks,
		"compression.type":              opts.CompressionType}

	switch strings.ToUpper(opts.SecurityProtocol) {
	case "PLAINTEXT":
		_ = kafkaConf.SetKey("security.protocol", "plaintext")
	case "SASL_SSL":
		_ = kafkaConf.SetKey("security.protocol", "sasl_ssl")
		_ = kafkaConf.SetKey("ssl.ca.location", "conf/mix-4096-ca-cert")
		_ = kafkaConf.SetKey("sasl.username", opts.SaslUsername)
		_ = kafkaConf.SetKey("sasl.password", opts.SaslPassword)
		_ = kafkaConf.SetKey("sasl.mechanism", opts.SaslMechanism)
		_ = kafkaConf.SetKey("enable.ssl.certificate.verification", "false")
	case "SASL_PLAINTEXT":
		_ = kafkaConf.SetKey("security.protocol", "sasl_plaintext")
		_ = kafkaConf.SetKey("sasl.username", opts.SaslUsername)
		_ = kafkaConf.SetKey("sasl.password", opts.SaslPassword)
		_ = kafkaConf.SetKey("sasl.mechanism", opts.SaslMechanism)
	default:
		return fmt.Errorf("unknown kafka protocol:%s", opts.SecurityProtocol)
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
					counterMetric.With(prometheus.Labels{"topic": "", "flag": "error"}).Inc()
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
		case <-opts.SignChan:
			logger.Info("producer get <-opts.SignChan")
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
	if err = k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &v.Topic, Partition: kafka.PartitionAny},
		Value:          []byte(v.Data),
	}, nil); err != nil {
		switch err.(kafka.Error).Code() {
		case kafka.ErrQueueFull:
			// Producer queue is full, wait 1s for messages to be delivered then try again.
			time.Sleep(time.Second)
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
