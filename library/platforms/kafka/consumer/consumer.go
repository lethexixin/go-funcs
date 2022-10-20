package consumer

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/lethexixin/go-funcs/common/logger"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"
)

type KafkaConsumer struct {
	Consumer *kafka.Consumer
}

var counterMetric *prometheus.CounterVec

func InitMetrics(appName string) {
	once := sync.Once{}
	once.Do(func() {
		counterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_kafka_consumer_metric_total", strings.ReplaceAll(appName, "-", "_")),
			Help: fmt.Sprintf("kafka consumer number of (topic,flag) for %s", strings.ReplaceAll(appName, "-", "_")),
		}, []string{"topic", "flag"})
		prometheus.MustRegister(counterMetric)
	})
}

type Options struct {
	BootstrapServers       string
	GroupId                string
	AutoOffsetReset        string
	HeartbeatIntervalMs    int
	SessionTimeoutMs       int
	MaxPollIntervalMs      int
	FetchMaxBytes          int
	MaxPartitionFetchBytes int
	SecurityProtocol       string
	SaslMechanism          string
	SaslUsername           string
	SaslPassword           string
	SubscribeTopics        []string

	SignChan         chan os.Signal
	ChanConsumerData chan MsgConsumerData
}

type MsgConsumerData struct {
	Data      []byte
	Partition int32
	Offset    int64
}

type Option func(*Options)

var (
	// kafka consumer config refer to https://help.aliyun.com/document_detail/68166.html and https://lovergamer.com/posts/kafka/kafka_consumers

	DefaultBootstrapServers = "127.0.0.1:9092"
	DefaultGroupId          = "test-group"
	// DefaultAutoOffsetReset 消费位点重置策略
	DefaultAutoOffsetReset = "latest"
	// DefaultHeartbeatIntervalMs 指定对消费者组协调器的心跳检查之间的间隔（以毫秒为单位）,以指示消费者处于活动状态并已连接
	DefaultHeartbeatIntervalMs = 3000
	// DefaultSessionTimeoutMs 指定消费者组中的消费者在被视为不活动之前可以与代理断开联系的最长时间（以毫秒为单位）
	DefaultSessionTimeoutMs = 30000
	// DefaultMaxPollIntervalMs 设置检查消费者是否继续处理消息的时间间隔
	DefaultMaxPollIntervalMs = 300000
	// DefaultFetchMaxBytes 设置一次从代理获取的数据量的最大字节数限制
	DefaultFetchMaxBytes = 52428800
	// DefaultMaxPartitionFetchBytes 设置为每个分区返回多少数据的最大字节数限制
	DefaultMaxPartitionFetchBytes = 1048576
	DefaultSecurityProtocol       = "PLAINTEXT"
	DefaultSaslMechanism          = "PLAIN"
	DefaultSaslUsername           = "kafka"
	DefaultSaslPassword           = "123456"
	DefaultSubscribeTopics        = []string{"test-topic"}
)

func BootstrapServers(bootstrapServers string) Option {
	return func(o *Options) {
		o.BootstrapServers = bootstrapServers
	}
}

func GroupId(groupId string) Option {
	return func(o *Options) {
		o.GroupId = groupId
	}
}

func AutoOffsetReset(autoOffsetReset string) Option {
	return func(o *Options) {
		o.AutoOffsetReset = autoOffsetReset
	}
}

func HeartbeatIntervalMs(heartbeatIntervalMs int) Option {
	return func(o *Options) {
		o.HeartbeatIntervalMs = heartbeatIntervalMs
	}
}

func SessionTimeoutMs(sessionTimeoutMs int) Option {
	return func(o *Options) {
		o.SessionTimeoutMs = sessionTimeoutMs
	}
}

func MaxPollIntervalMs(maxPollIntervalMs int) Option {
	return func(o *Options) {
		o.MaxPollIntervalMs = maxPollIntervalMs
	}
}

func FetchMaxBytes(fetchMaxBytes int) Option {
	return func(o *Options) {
		o.FetchMaxBytes = fetchMaxBytes
	}
}

func MaxPartitionFetchBytes(maxPartitionFetchBytes int) Option {
	return func(o *Options) {
		o.MaxPartitionFetchBytes = maxPartitionFetchBytes
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

func SubscribeTopics(subscribeTopics []string) Option {
	return func(o *Options) {
		o.SubscribeTopics = subscribeTopics
	}
}

func SignChan(signChan chan os.Signal) Option {
	return func(o *Options) {
		o.SignChan = signChan
	}
}

func ChanConsumerData(chanConsumerData chan MsgConsumerData) Option {
	return func(o *Options) {
		o.ChanConsumerData = chanConsumerData
	}
}

//InitConsumer
//
//1. confluent-kafka-go build refer to docs/kafka.md
func (k *KafkaConsumer) InitConsumer(options ...Option) (err error) {
	opts := Options{
		BootstrapServers:       DefaultBootstrapServers,
		GroupId:                DefaultGroupId,
		AutoOffsetReset:        DefaultAutoOffsetReset,
		HeartbeatIntervalMs:    DefaultHeartbeatIntervalMs,
		SessionTimeoutMs:       DefaultSessionTimeoutMs,
		MaxPollIntervalMs:      DefaultMaxPollIntervalMs,
		FetchMaxBytes:          DefaultFetchMaxBytes,
		MaxPartitionFetchBytes: DefaultMaxPartitionFetchBytes,
		SecurityProtocol:       DefaultSecurityProtocol,
		SaslMechanism:          DefaultSaslMechanism,
		SaslUsername:           DefaultSaslUsername,
		SaslPassword:           DefaultSaslPassword,
		SubscribeTopics:        DefaultSubscribeTopics,
	}

	for _, o := range options {
		o(&opts)
	}

	if opts.ChanConsumerData == nil || opts.SignChan == nil {
		return errors.New("opts.ChanConsumerData == nil || opts.SignChan == nil")
	}

	kafkaConf := &kafka.ConfigMap{
		"bootstrap.servers":         opts.BootstrapServers,
		"group.id":                  opts.GroupId,
		"api.version.request":       "true",
		"auto.offset.reset":         opts.AutoOffsetReset,
		"heartbeat.interval.ms":     opts.HeartbeatIntervalMs,
		"session.timeout.ms":        opts.SessionTimeoutMs,
		"max.poll.interval.ms":      opts.MaxPollIntervalMs,
		"fetch.max.bytes":           opts.FetchMaxBytes,
		"max.partition.fetch.bytes": opts.MaxPartitionFetchBytes}

	switch strings.ToUpper(opts.SecurityProtocol) {
	case "PLAINTEXT":
		_ = kafkaConf.SetKey("security.protocol", "plaintext")
	case "SASL_SSL":
		_ = kafkaConf.SetKey("security.protocol", "sasl_ssl")
		_ = kafkaConf.SetKey("ssl.ca.location", "./conf/mix-4096-ca-cert")
		_ = kafkaConf.SetKey("sasl.username", opts.SaslUsername)
		_ = kafkaConf.SetKey("sasl.password", opts.SaslPassword)
		_ = kafkaConf.SetKey("sasl.mechanism", opts.SaslMechanism)
	case "SASL_PLAINTEXT":
		_ = kafkaConf.SetKey("security.protocol", "sasl_plaintext")
		_ = kafkaConf.SetKey("sasl.username", opts.SaslUsername)
		_ = kafkaConf.SetKey("sasl.password", opts.SaslPassword)
		_ = kafkaConf.SetKey("sasl.mechanism", opts.SaslMechanism)
	default:
		return fmt.Errorf("unknown kafka protocol:%s", opts.SecurityProtocol)
	}

	k.Consumer, err = kafka.NewConsumer(kafkaConf)
	if err != nil {
		logger.Errorf("failed to create kafka consumer:%s", err.Error())
		return err
	}
	defer k.Consumer.Close()

	logger.Info("create kafka consumer successful")

	err = k.Consumer.SubscribeTopics(opts.SubscribeTopics, nil)
	if err != nil {
		logger.Errorf("kafka consumer subscribe topics err:%s", err.Error())
		return err
	}

	run := true
	for run {
		select {
		case <-opts.SignChan:
			logger.Info("consumer get <-opts.SignChan")
			run = false
		default:
			msg, err := k.Consumer.ReadMessage(-1)
			if err != nil {
				// The client will automatically try to recover from all errors.
				logger.Errorf("consumer read msg err:%s", err.Error())
				if counterMetric != nil {
					counterMetric.With(prometheus.Labels{"topic": strings.Join(opts.SubscribeTopics, ","), "flag": "error"}).Inc()
				}
			} else {
				opts.ChanConsumerData <- MsgConsumerData{
					Data:      msg.Value,
					Offset:    int64(msg.TopicPartition.Offset),
					Partition: msg.TopicPartition.Partition,
				}
				if counterMetric != nil {
					counterMetric.With(prometheus.Labels{"topic": strings.Join(opts.SubscribeTopics, ","), "flag": "success"}).Inc()
				}
			}
		}
	}

	return nil
}
