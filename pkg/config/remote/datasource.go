package remote

import (
	"errors"
	"io"

	"github.com/spf13/viper"
)

const (
	DataSourceApollo        = "apollo"
	DataSourceEtcd          = "etcd"
	DataSourceNeteaseConfig = "ncm"
)

var (
	// ErrConfigAddr not config
	ErrConfigAddr = errors.New("no config... ")
	// ErrInvalidDataSource defines an error that the scheme has been registered
	ErrInvalidDataSource = errors.New("invalid data source, please make sure the scheme has been registered")
	datasourceBuilders   = make(map[string]DataSourceCreatorFunc)
)

// DataSourceCreatorFunc represents a dataSource creator function
type DataSourceCreatorFunc func(viper.RemoteProvider) DataSource

// DataSource ...
type DataSource interface {
	ReadConfig() ([]byte, error)
	IsConfigChanged() <-chan struct{}
	io.Closer
}

// Register registers a dataSource creator function to the registry
func Register(scheme string, creator DataSourceCreatorFunc) {
	datasourceBuilders[scheme] = creator
}

// NewDataSource create datasource or use existing instance
func NewDataSource(rp viper.RemoteProvider) (DataSource, error) {
	creatorFunc, exist := datasourceBuilders[rp.Provider()]
	if !exist {
		return nil, ErrInvalidDataSource
	}
	return creatorFunc(rp), nil
}
