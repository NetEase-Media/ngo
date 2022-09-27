package remote

import (
	"bytes"
	"io"

	"github.com/spf13/viper"
)

type remoteConfigProvider struct {
	datasource map[string]DataSource
}

func (rc remoteConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	ds, err := rc.getDs(rp)
	if err != nil {
		return nil, err
	}
	b, err := ds.ReadConfig()
	return bytes.NewReader(b), err
}

func (rc remoteConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	ds, err := rc.getDs(rp)
	if err != nil {
		return nil, err
	}
	b, err := ds.ReadConfig()
	return bytes.NewReader(b), err
}

func (rc remoteConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	ds, err := rc.getDs(rp)
	if err != nil {
		return nil, nil
	}

	quit := make(chan bool)
	viperResponsCh := make(chan *viper.RemoteResponse)

	go func() {
		for {
			select {
			case <-quit:
				ds.Close()
				delete(rc.datasource, rp.Provider())
				return
			case <-ds.IsConfigChanged():
				data, err := ds.ReadConfig()
				viperResponsCh <- &viper.RemoteResponse{
					Error: err,
					Value: data,
				}
			}
		}
	}()

	return viperResponsCh, quit
}

func (rc remoteConfigProvider) getDs(rp viper.RemoteProvider) (DataSource, error) {
	ds := rc.datasource[rp.Provider()]
	if ds != nil {
		return ds, nil
	}
	ds, err := NewDataSource(rp)
	if err != nil {
		return nil, err
	}
	// Note: 理论上初始化时候 ReadRemoteConfig 才会进入这里的赋值。其他情况都是只读，不存在 map 并发读写
	rc.datasource[rp.Provider()] = ds
	return ds, nil
}

func init() {
	viper.RemoteConfig = &remoteConfigProvider{datasource: make(map[string]DataSource)}
}
