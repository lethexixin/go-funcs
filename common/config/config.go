package config

import (
	"bytes"
	"errors"
	"net"
	"strconv"
)

import (
	"github.com/lethexixin/go-funcs/common/logger"
)

import (
	"github.com/BurntSushi/toml"
	"github.com/hashicorp/consul/api"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type FileConfig struct {
	Path string `json:"path"`
}

type NacosConfig struct {
	Client config_client.IConfigClient

	Addr      string `json:"addr"`
	Namespace string `json:"namespace"`
	DataId    string `json:"dataId"`
	Group     string `json:"group"`
}

type ConsulConfig struct {
	Client *api.Client

	Addr      string `json:"addr"`
	Namespace string `json:"namespace"`
	ServiceId string `json:"serviceId"`
}

var (
	ErrAddrIsEmpty    = errors.New("config center addr is empty")
	ErrContentIsEmpty = errors.New("config content is empty")
)

// LoadFileConfig 引导 file 配置数据给 conf
func (c *FileConfig) LoadFileConfig(conf interface{}) (err error) {
	if _, err = toml.DecodeFile(c.Path, conf); err != nil {
		logger.Errorf("file, failed to decode content from local file:%s, err:%s", c.Path, err.Error())
		return err
	}
	logger.Infof("file, load config successful from local file:%s", c.Path)
	return nil
}

// LoadNacosConfig 引导 nacos 配置数据给 conf
func (c *NacosConfig) LoadNacosConfig(conf interface{}) (err error) {
	if len(c.Addr) == 0 {
		return ErrAddrIsEmpty
	}

	// nacos 相关参数配置,具体配置可参考 https://github.com/nacos-group/nacos-sdk-go

	ipAddr, hPort, _ := net.SplitHostPort(c.Addr)
	port, _ := strconv.Atoi(hPort)
	//create ServerConfig
	sc := []constant.ServerConfig{*constant.NewServerConfig(ipAddr, uint64(port), constant.WithContextPath("/nacos"))}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(c.Namespace),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("tmp/nacos/log"),
		constant.WithCacheDir("tmp/nacos/cache"),
		constant.WithLogLevel("error"),
	)

	// create config client
	c.Client, err = clients.NewConfigClient(vo.NacosClientParam{ClientConfig: &cc, ServerConfigs: sc})
	if err != nil {
		logger.Errorf("failed create nacos:%s client, err:%s", c.Addr, err.Error())
		return err
	}

	// get nacos config
	content, err := c.Client.GetConfig(vo.ConfigParam{
		DataId: c.DataId,
		Group:  c.Group,
	})
	if err != nil {
		logger.Errorf("nacos, failed to get config content from addr:%s, namespaceId:%s, dataId:%s, group:%s, err:%s", c.Addr, c.Namespace, c.DataId, c.Group, err.Error())
		return err
	}

	if len(content) == 0 {
		logger.Errorf("nacos, config content is empty from addr:%s, namespaceId:%s, dataId:%s, group:%s", c.Addr, c.Namespace, c.DataId, c.Group)
		return ErrContentIsEmpty
	}

	if _, err = toml.NewDecoder(bytes.NewBuffer([]byte(content))).Decode(conf); err != nil {
		logger.Errorf("nacos, failed to decode content from addr:%s, namespaceId:%s, dataId:%s, group:%s, err:%s", c.Addr, c.Namespace, c.DataId, c.Group, err.Error())
		return err
	}

	logger.Infof("nacos, load config successful from addr:%s, namespaceId:%s, dataId:%s, group:%s", c.Addr, c.Namespace, c.DataId, c.Group)
	return nil
}

// ListenConfig 配置监听
func (c *NacosConfig) ListenConfig(conf interface{}) {
	_ = c.Client.ListenConfig(vo.ConfigParam{
		DataId: c.DataId, Group: c.Group, OnChange: func(namespace, group, dataId, data string) {
			if len(data) != 0 {
				if _, err := toml.NewDecoder(bytes.NewBuffer([]byte(data))).Decode(conf); err != nil {
					logger.Errorf("nacos, failed to decode content from addr:%s, namespaceId:%s, dataId:%s, group:%s, err:%s", c.Addr, c.Namespace, c.DataId, c.Group, err.Error())
				}
			}
		}})
}

// CancelListenConfig 取消配置监听
func (c *NacosConfig) CancelListenConfig() (err error) {
	if err = c.Client.CancelListenConfig(vo.ConfigParam{DataId: c.DataId, Group: c.Group}); err != nil {
		logger.Errorf("nacos, failed to cancel config listen from addr:%s, namespaceId:%s, dataId:%s, group:%s, err:%s", c.Addr, c.Namespace, c.DataId, c.Group, err.Error())
		return err
	}
	return nil
}

// LoadConsulConfig 引导 consul 配置数据给 conf
func (c *ConsulConfig) LoadConsulConfig(conf interface{}) (err error) {
	if len(c.Addr) == 0 {
		return ErrAddrIsEmpty
	}

	c.Client, err = api.NewClient(&api.Config{
		Address:   c.Addr,
		Namespace: c.Namespace,
	})
	if err != nil {
		logger.Errorf("failed create consul:%s client, err:%s", c.Addr, err.Error())
		return err
	}

	content, _, err := c.Client.KV().Get(c.ServiceId, nil)
	if err != nil {
		logger.Errorf("consul, failed to get config content from addr:%s, serviceId:%s, err:%s", c.Addr, c.ServiceId, err.Error())
		return err
	}

	if content == nil {
		logger.Errorf("consul, config content is empty from addr:%s, serviceId:%s", c.Addr, c.ServiceId)
		return ErrContentIsEmpty
	}

	if _, err = toml.NewDecoder(bytes.NewBuffer(content.Value)).Decode(conf); err != nil {
		logger.Errorf("consul, failed to decode content from addr:%s, serviceId:%s, err:%s", c.Addr, c.ServiceId, err.Error())
		return err
	}

	logger.Infof("consul, read config successful from addr:%s, serviceId:%s", c.Addr, c.ServiceId)
	return nil
}
