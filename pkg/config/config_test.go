package config

import (
	"testing"

	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testConfigAddrYaml = "./testdata/app.yaml"
)

func TestConfiguration_Sub(t *testing.T) {
	configAddr := testConfigAddrYaml
	watch := true
	cfg, err := New(configAddr, watch)
	assert.NoError(t, err)
	sub := cfg.Sub("service")
	assert.Equal(t, "ngo", sub.GetString("appName"))
}

func TestConfiguration_SubSlice(t *testing.T) {
	configAddr := testConfigAddrYaml
	watch := true
	cfg, err := New(configAddr, watch)
	assert.NoError(t, err)
	subSlice := cfg.SubSlice("redis")
	assert.Equal(t, "client1", subSlice[0].GetString("name"))
}

func TestConfiguration_Watch(t *testing.T) {
	configAddr := testConfigAddrYaml
	watch := true
	cfg, err := New(configAddr, watch)
	assert.NoError(t, err)
	cfg.OnChange(func(_ *Configuration) {
		log.Info("change")
	})
}

func TestEmpty(t *testing.T) {

}

func TestNewWithFile(t *testing.T) {
	configAddr := testConfigAddrYaml
	watch := true
	cfg, err := New(configAddr, watch)
	assert.NoError(t, err)
	assert.Equal(t, "ngo", cfg.GetString("service.appName"))
	assert.Equal(t, "value", cfg.GetString("propkey"))
	assert.Equal(t, []string{"1", "2"}, cfg.GetStringSlice("slice"))
}

func TestNewWithApollo(t *testing.T) {
	t.Skip()
	configAddr := "apollo://106.54.227.205:8080?appId=ngo&cluster=ngo-demo&namespaceNames=app.yaml,httpServer.yaml"
	watch := true
	cfg, err := New(configAddr, watch)
	assert.NoError(t, err)
	assert.Equal(t, "ngo", cfg.GetString("service.appName"))
	assert.Equal(t, 8080, cfg.GetInt("httpServer.port"))
	cfg.OnChange(func(_ *Configuration) {
		log.Info("change")
	})
}

func TestNewWithEtcd(t *testing.T) {
	configAddr := "etcd://any.com?endpoints=10.201.209.134:2379&endpoints=10.201.209.134:2379&user_name=&password=&read_timeout_sec=3&key=brain/dev/ci_test_key&with_prefix=false"
	cfg, err := New(configAddr, true)
	require.NoError(t, err)
	assert.Equal(t, "test_value", cfg.GetString("test_key"))
	cfg.OnChange(func(_ *Configuration) {
		log.Info("change")
		log.Info(cfg.AllSettings())
	})
	// time.Sleep(time.Second * 15)
}
