package include

// 不要改变顺序，按需编译
import (
	// log
	_ "github.com/NetEase-Media/ngo/pkg/log/r"
	// service base info
	_ "github.com/NetEase-Media/ngo/pkg/service/r"

	// tracing
	_ "github.com/NetEase-Media/ngo/pkg/tracing/opentracing/optimus"
	_ "github.com/NetEase-Media/ngo/pkg/tracing/r"

	// sentinel
	_ "github.com/NetEase-Media/ngo/pkg/sentinel/r"
	// client
	_ "github.com/NetEase-Media/ngo/pkg/client/db/r"
	_ "github.com/NetEase-Media/ngo/pkg/client/httplib/r"
	_ "github.com/NetEase-Media/ngo/pkg/client/kafka/r"
	_ "github.com/NetEase-Media/ngo/pkg/client/memcache/r"
	_ "github.com/NetEase-Media/ngo/pkg/client/multicache/r"
	_ "github.com/NetEase-Media/ngo/pkg/client/redis/r"
	_ "github.com/NetEase-Media/ngo/pkg/client/zookeeper/r"

	// dlock, must load redis first
	_ "github.com/NetEase-Media/ngo/pkg/dlock/r"
	// job
	_ "github.com/NetEase-Media/ngo/pkg/job/xxljob/r"
	// http server
	_ "github.com/NetEase-Media/ngo/pkg/server/http/r"
	// http admin
	_ "github.com/NetEase-Media/ngo/pkg/server/http/admin"
)
