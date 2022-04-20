package db
//
//import (
//	"strings"
//	"time"
//
//	"github.com/NetEase-Media/ngo/pkg/metrics"
//	collectors "github.com/NetEase-Media/ngo/pkg/metrics/colloctors"
//	"gorm.io/gorm"
//)
//
//const ngoMetricsKey = "ngo:db:metrics"
//
//type gormMetricsPlugin struct{}
//
//func newGormMetricsPlugin() *gormMetricsPlugin {
//	return &gormMetricsPlugin{}
//}
//
//func (p *gormMetricsPlugin) Name() string {
//	return "ngo:db:metrics"
//}
//
//func (p *gormMetricsPlugin) Initialize(db *gorm.DB) error {
//	p.registerCallbacks(db)
//	return nil
//}
//
//func (p *gormMetricsPlugin) registerCallbacks(db *gorm.DB) {
//
//	db.Callback().Query().Before("gorm:query").Register("ngo:metrics:before_query", p.metricBefore)
//	db.Callback().Query().After("gorm:query").Register("ngo:metrics:after_query", p.metricAfter)
//
//	db.Callback().Create().Before("gorm:create").Register("ngo:metrics:before_create", p.metricBefore)
//	db.Callback().Create().After("gorm:create").Register("ngo:metrics:after_create", p.metricAfter)
//
//	db.Callback().Update().Before("gorm:update").Register("ngo:metrics:before_update", p.metricBefore)
//	db.Callback().Update().After("gorm:update").Register("ngo:metrics:after_update", p.metricAfter)
//
//	db.Callback().Delete().Before("gorm:delete").Register("ngo:metrics:before_delete", p.metricBefore)
//	db.Callback().Delete().After("gorm:delete").Register("ngo:metrics:after_delete", p.metricAfter)
//
//	db.Callback().Row().Before("gorm:row").Register("ngo:metrics:before_row", p.metricBefore)
//	db.Callback().Row().After("gorm:row").Register("ngo:metrics:after_row", p.metricAfter)
//
//	db.Callback().Raw().Before("gorm:raw").Register("ngo:metrics:before_raw", p.metricBefore)
//	db.Callback().Raw().After("gorm:raw").Register("ngo:metrics:after_raw", p.metricAfter)
//}
//
////region callbacks
//
//func (p *gormMetricsPlugin) metricBefore(db *gorm.DB) {
//	if !metrics.IsMetricsEnabled() {
//		return
//	}
//	if db == nil || db.Statement == nil || db.Statement.Context == nil {
//		return
//	}
//	now := time.Now()
//	db.InstanceSet(ngoMetricsKey, now)
//}
//
//func (p *gormMetricsPlugin) metricAfter(db *gorm.DB) {
//	if !metrics.IsMetricsEnabled() {
//		return
//	}
//	if db == nil || db.Statement == nil || db.Statement.Context == nil {
//		return
//	}
//	value, ok := db.InstanceGet(ngoMetricsKey)
//	if !ok || value == nil {
//		return
//	}
//	startTime, ok := value.(time.Time)
//	if !ok {
//		return
//	}
//	dsn := getDsn(db.Dialector)
//	sql := sqlFilter(db.Statement.SQL.String())
//	if db.Error != nil {
//		collectors.MysqlCollector().OneStepComplete(startTime, dsn, sql, 0, 0, 0, db.Error)
//		return
//	}
//	var updatedRowCount int
//	var readRowCount int
//	if strings.HasPrefix(strings.ToLower(sql), "select") {
//		readRowCount = int(db.RowsAffected)
//	} else {
//		updatedRowCount = int(db.RowsAffected)
//	}
//	collectors.MysqlCollector().OneStepComplete(startTime, dsn, sql, getConnectionCount(db), updatedRowCount, readRowCount, nil)
//}
//
////endregion
