package db

import (
	"runtime"

	"github.com/NetEase-Media/ngo/pkg/tracing"

	"gorm.io/gorm"
)

const ngoTracerSpanKey = "ngo:tracer:span:db"

type gormTracerPlugin struct{}

func newGormTracerPlugin() *gormTracerPlugin {
	return &gormTracerPlugin{}
}

func (p *gormTracerPlugin) Name() string {
	return "ngo:tracer"
}

func (p *gormTracerPlugin) Initialize(db *gorm.DB) error {
	p.registerCallbacks(db)
	return nil
}

func (p *gormTracerPlugin) registerCallbacks(db *gorm.DB) {

	db.Callback().Query().Before("gorm:query").Register("ngo:tracer:before_query", p.traceBefore)
	db.Callback().Query().After("gorm:query").Register("ngo:tracer:after_query", p.traceAfter)

	db.Callback().Create().Before("gorm:create").Register("ngo:tracer:before_create", p.traceBefore)
	db.Callback().Create().After("gorm:create").Register("ngo:tracer:after_create", p.traceAfter)

	db.Callback().Update().Before("gorm:update").Register("ngo:tracer:before_update", p.traceBefore)
	db.Callback().Update().After("gorm:update").Register("ngo:tracer:after_update", p.traceAfter)

	db.Callback().Delete().Before("gorm:delete").Register("ngo:tracer:before_delete", p.traceBefore)
	db.Callback().Delete().After("gorm:delete").Register("ngo:tracer:after_delete", p.traceAfter)

	db.Callback().Row().Before("gorm:row").Register("ngo:tracer:before_row", p.traceBefore)
	db.Callback().Row().After("gorm:row").Register("ngo:tracer:after_row", p.traceAfter)

	db.Callback().Raw().Before("gorm:raw").Register("ngo:tracer:before_raw", p.traceBefore)
	db.Callback().Raw().After("gorm:raw").Register("ngo:tracer:after_raw", p.traceAfter)
}

func (p *gormTracerPlugin) traceBefore(db *gorm.DB) {
	if !tracing.Enabled() {
		return
	}
	if db == nil || db.Statement == nil || db.Statement.Context == nil {
		return
	}

	span, _ := tracing.StartSpanFromContext(db.Statement.Context, "gorm")
	db.InstanceSet(ngoTracerSpanKey, span)
}

func (p *gormTracerPlugin) traceAfter(db *gorm.DB) {
	if !tracing.Enabled() {
		return
	}
	if db == nil || db.Statement == nil || db.Statement.Context == nil {
		return
	}
	value, ok := db.InstanceGet(ngoTracerSpanKey)
	if !ok || value == nil {
		return
	}
	span, ok := value.(tracing.Span)
	if !ok || span == nil {
		return
	}
	defer span.Finish()

	tracing.SpanKind.Set(span, "client")
	tracing.SpanType.Set(span, "JDBC")
	tracing.PluginType.Set(span, "gorm")
	dsn := getDsn(db.Dialector)
	tracing.DbUrl.Set(span, dsn)
	tracing.DbType.Set(span, db.Dialector.Name())
	tracing.DbSql.Set(span, db.Statement.SQL.String())
	tracing.RemoteAppName.Set(span, "JDBC")
	tracing.RemoteClusterName.Set(span, dsn)
	if db.Error != nil {
		stack := make([]byte, 2048)
		runtime.Stack(stack, false)
		tracing.ExceptionName.Set(span, db.Error)
		tracing.ExceptionMessage.Set(span, db.Error.Error())
		tracing.ExceptionStacktrace.Set(span, string(stack))
	}
}
