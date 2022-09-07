package etcd

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/NetEase-Media/ngo/pkg/config/remote"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/NetEase-Media/ngo/pkg/util"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/viper"
)

func init() {
	remote.Register(remote.DataSourceEtcd, func(rp viper.RemoteProvider) remote.DataSource {
		urlObj, err := url.Parse(fmt.Sprintf("%s://%s%s", rp.Provider(), rp.Endpoint(), rp.Path()))
		if err != nil {
			log.Error("parse configAddr error", err)
			return nil
		}

		req, err := http.NewRequest("GET", urlObj.String(), nil)
		util.CheckError(err)
		var option Option
		err = binding.Query.Bind(req, &option)
		util.CheckError(err)

		source, err := NewDataSource(option)
		util.CheckError(err)
		return source
	})
}
