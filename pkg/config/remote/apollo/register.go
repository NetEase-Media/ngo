package apollo

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/NetEase-Media/ngo/pkg/config/remote"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/philchia/agollo/v4"
	"github.com/spf13/viper"
)

func init() {
	remote.Register(remote.DataSourceApollo, func(rp viper.RemoteProvider) remote.DataSource {
		urlObj, err := url.Parse(fmt.Sprintf("%s://%s%s", rp.Provider(), rp.Endpoint(), rp.Path()))
		if err != nil {
			log.Error("parse configAddr error", err)
			return nil
		}

		ns := urlObj.Query().Get("namespaceNames")
		namespaceNames := strings.Split(ns, ",")
		apolloConf := agollo.Conf{
			AppID:              urlObj.Query().Get("appId"),
			Cluster:            urlObj.Query().Get("cluster"),
			NameSpaceNames:     namespaceNames,
			MetaAddr:           urlObj.Host + urlObj.Path,
			InsecureSkipVerify: true,
			AccesskeySecret:    urlObj.Query().Get("accesskeySecret"),
			CacheDir:           ".",
		}
		if urlObj.Query().Get("insecureSkipVerify") == "false" {
			apolloConf.InsecureSkipVerify = false
		}
		if urlObj.Query().Get("cacheDir") != "" {
			apolloConf.CacheDir = urlObj.Query().Get("cacheDir")
		}
		return NewDataSource(&apolloConf, namespaceNames)
	})
}
