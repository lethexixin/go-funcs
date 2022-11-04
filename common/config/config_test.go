package config

import (
	"encoding/json"
	"testing"
	"time"
)

func TestLoadFileConfig(t *testing.T) {
	conf := make(map[string]interface{})
	c := &FileConfig{Path: "config_test.toml"}
	err := c.LoadFileConfig(&conf)
	t.Log(conf, err)
}

func TestLoadNacosConfig(t *testing.T) {
	conf := make(map[string]interface{})
	c := &NacosConfig{
		Addr:      "127.0.0.1:8848",
		Namespace: "",
		DataId:    "test",
		Group:     "DEFAULT_GROUP",
	}
	err := c.LoadNacosConfig(&conf)
	t.Log(conf, err)

	// 测试配置监听
	go func() {
		preConf, _ := json.Marshal(conf)
		for {
			time.Sleep(time.Second)
			currentConf, _ := json.Marshal(conf)
			if string(currentConf) != string(preConf) {
				t.Log("currentConf:", conf)
				preConf = currentConf
			}
		}
	}()

	t.Log("start c.ListenConfig(&conf)")
	c.ListenConfig(&conf)

	go func() {
		// 模拟60秒之后退出配置监听
		time.Sleep(time.Second * 60)
		err = c.CancelListenConfig()
		t.Log("c.CancelListenConfig():", err)
	}()

	time.Sleep(time.Second * 70)
}

func TestFileConsulConfig(t *testing.T) {
	conf := make(map[string]interface{})
	c := &ConsulConfig{
		Addr:      "127.0.0.1:8500",
		ServiceId: "test",
	}
	err := c.LoadConsulConfig(&conf)
	t.Log(conf, err)
}
