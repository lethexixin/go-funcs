package cryptos

import (
	"testing"
)

const (
	testKey = "b6c1cd0fe6e55f22fb483096822b5d1c"
)

func TestAESCBC(t *testing.T) {
	bts, err := EncodeAESCBC([]byte(`{"name" : "xin"}`),
		[]byte(testKey))
	t.Log("EncodeAESCBC,", string(bts), err)

	bts, err = DecodeAESCBC(bts, []byte(testKey))
	t.Log("DecodeAESCBC,", string(bts), err)
}

func TestAESGCM(t *testing.T) {
	bts, _, err := EncodeAESGCM([]byte(`{"name" : "xin"}`),
		[]byte(testKey))
	t.Log("EncodeAESGCM,", string(bts), err)

	bts, err = DecodeAESGCM(bts, []byte(testKey))
	t.Log("DecodeAESGCM,", string(bts), err)
}
