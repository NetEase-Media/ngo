package k8sjob

import (
	"os"

	"github.com/NetEase-Media/ngo/pkg/config"
	"github.com/NetEase-Media/ngo/pkg/util"
)

// Run 运行函数并退出进程
func Run(f Callback) {
	var opt Options
	err := config.UnmarshalKey("job", &opt)
	util.CheckError(err)
	opt.check()

	hostname, err := os.Hostname()
	util.CheckError(err)

	j := &job{
		opt:      &opt,
		f:        f,
		hostname: hostname,
	}
	j.run()

	os.Exit(0)
}
