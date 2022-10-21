package consumer

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

import (
	"github.com/lethexixin/go-funcs/common/logger"
)

import (
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
	bootstrapServers       string
	groupId                string
	autoOffsetReset        string
	heartbeatIntervalMs    int
	sessionTimeoutMs       int
	maxPollIntervalMs      int
	fetchMaxBytes          int
	maxPartitionFetchBytes int
	securityProtocol       string
	saslMechanism          string
	saslUsername           string
	saslPassword           string
	subscribeTopics        []string

	signChan         chan os.Signal
	chanConsumerData chan MsgConsumerData
}

type MsgConsumerData struct {
	Data      []byte
	Partition int32
	Offset    int64
}

type Option func(*Options)

const (
	// kafka consumer config refer to https://help.aliyun.com/document_detail/68166.html and https://lovergamer.com/posts/kafka/kafka_consumers

	DefaultBootstrapServers = "127.0.0.1:9092"
	DefaultGroupId          = "test-group"
	// DefaultAutoOffsetReset 消费位点重置策略
	DefaultAutoOffsetReset = "latest"
	// DefaultHeartbeatIntervalMs 指定对消费者组协调器的心跳检查之间的间隔(以毫秒为单位),以指示消费者处于活动状态并已连接
	DefaultHeartbeatIntervalMs = 3000
	// DefaultSessionTimeoutMs 指定消费者组中的消费者在被视为不活动之前可以与代理断开联系的最长时间(以毫秒为单位)
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
	DefaultSubscribeTopics        = "test-topic"
)

func BootstrapServers(bootstrapServers string) Option {
	return func(o *Options) {
		o.bootstrapServers = bootstrapServers
	}
}

func GroupId(groupId string) Option {
	return func(o *Options) {
		o.groupId = groupId
	}
}

func AutoOffsetReset(autoOffsetReset string) Option {
	return func(o *Options) {
		o.autoOffsetReset = autoOffsetReset
	}
}

func HeartbeatIntervalMs(heartbeatIntervalMs int) Option {
	return func(o *Options) {
		o.heartbeatIntervalMs = heartbeatIntervalMs
	}
}

func SessionTimeoutMs(sessionTimeoutMs int) Option {
	return func(o *Options) {
		o.sessionTimeoutMs = sessionTimeoutMs
	}
}

func MaxPollIntervalMs(maxPollIntervalMs int) Option {
	return func(o *Options) {
		o.maxPollIntervalMs = maxPollIntervalMs
	}
}

func FetchMaxBytes(fetchMaxBytes int) Option {
	return func(o *Options) {
		o.fetchMaxBytes = fetchMaxBytes
	}
}

func MaxPartitionFetchBytes(maxPartitionFetchBytes int) Option {
	return func(o *Options) {
		o.maxPartitionFetchBytes = maxPartitionFetchBytes
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

func SubscribeTopics(subscribeTopics []string) Option {
	return func(o *Options) {
		o.subscribeTopics = subscribeTopics
	}
}

func SignChan(signChan chan os.Signal) Option {
	return func(o *Options) {
		o.signChan = signChan
	}
}

func ChanConsumerData(chanConsumerData chan MsgConsumerData) Option {
	return func(o *Options) {
		o.chanConsumerData = chanConsumerData
	}
}

//InitConsumer
//
//1. confluent-kafka-go build refer to docs/kafka.md
func (k *KafkaConsumer) InitConsumer(options ...Option) (err error) {
	opts := Options{
		bootstrapServers:       DefaultBootstrapServers,
		groupId:                DefaultGroupId,
		autoOffsetReset:        DefaultAutoOffsetReset,
		heartbeatIntervalMs:    DefaultHeartbeatIntervalMs,
		sessionTimeoutMs:       DefaultSessionTimeoutMs,
		maxPollIntervalMs:      DefaultMaxPollIntervalMs,
		fetchMaxBytes:          DefaultFetchMaxBytes,
		maxPartitionFetchBytes: DefaultMaxPartitionFetchBytes,
		securityProtocol:       DefaultSecurityProtocol,
		saslMechanism:          DefaultSaslMechanism,
		saslUsername:           DefaultSaslUsername,
		saslPassword:           DefaultSaslPassword,
		subscribeTopics:        []string{DefaultSubscribeTopics},
	}

	for _, o := range options {
		o(&opts)
	}

	if opts.chanConsumerData == nil || opts.signChan == nil {
		return errors.New("opts.chanConsumerData == nil || opts.signChan == nil")
	}

	kafkaConf := &kafka.ConfigMap{
		"bootstrap.servers":         opts.bootstrapServers,
		"group.id":                  opts.groupId,
		"api.version.request":       "true",
		"auto.offset.reset":         opts.autoOffsetReset,
		"heartbeat.interval.ms":     opts.heartbeatIntervalMs,
		"session.timeout.ms":        opts.sessionTimeoutMs,
		"max.poll.interval.ms":      opts.maxPollIntervalMs,
		"fetch.max.bytes":           opts.fetchMaxBytes,
		"max.partition.fetch.bytes": opts.maxPartitionFetchBytes}

	switch strings.ToUpper(opts.securityProtocol) {
	case "PLAINTEXT":
		_ = kafkaConf.SetKey("security.protocol", "plaintext")
	case "SASL_SSL":
		_ = kafkaConf.SetKey("security.protocol", "sasl_ssl")
		_ = kafkaConf.SetKey("ssl.ca.location", "./conf/mix-4096-ca-cert")
		_ = kafkaConf.SetKey("sasl.username", opts.saslUsername)
		_ = kafkaConf.SetKey("sasl.password", opts.saslPassword)
		_ = kafkaConf.SetKey("sasl.mechanism", opts.saslMechanism)
	case "SASL_PLAINTEXT":
		_ = kafkaConf.SetKey("security.protocol", "sasl_plaintext")
		_ = kafkaConf.SetKey("sasl.username", opts.saslUsername)
		_ = kafkaConf.SetKey("sasl.password", opts.saslPassword)
		_ = kafkaConf.SetKey("sasl.mechanism", opts.saslMechanism)
	default:
		return fmt.Errorf("unknown kafka protocol:%s", opts.securityProtocol)
	}

	k.Consumer, err = kafka.NewConsumer(kafkaConf)
	if err != nil {
		logger.Errorf("failed to create kafka consumer:%s", err.Error())
		return err
	}
	defer k.Consumer.Close()

	logger.Info("create kafka consumer successful")

	err = k.Consumer.SubscribeTopics(opts.subscribeTopics, nil)
	if err != nil {
		logger.Errorf("kafka consumer subscribe topics err:%s", err.Error())
		return err
	}

	run := true
	for run {
		select {
		case <-opts.signChan:
			logger.Info("consumer get <-opts.signChan")
			run = false
		default:
			msg, err := k.Consumer.ReadMessage(-1)
			if err != nil {
				// The client will automatically try to recover from all errors.
				logger.Errorf("consumer read msg err:%s", err.Error())
				if counterMetric != nil {
					counterMetric.With(prometheus.Labels{"topic": strings.Join(opts.subscribeTopics, ","), "flag": "error"}).Inc()
				}
			} else {
				opts.chanConsumerData <- MsgConsumerData{
					Data:      msg.Value,
					Offset:    int64(msg.TopicPartition.Offset),
					Partition: msg.TopicPartition.Partition,
				}
				if counterMetric != nil {
					counterMetric.With(prometheus.Labels{"topic": strings.Join(opts.subscribeTopics, ","), "flag": "success"}).Inc()
				}
			}
		}
	}

	return nil
}
