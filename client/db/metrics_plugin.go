// Copyright Ngo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"gorm.io/gorm"
)

const ngoMetricsKey = "ngo:db:metrics"

type gormMetricsPlugin struct{}

func newGormMetricsPlugin() *gormMetricsPlugin {
	return &gormMetricsPlugin{}
}

func (p *gormMetricsPlugin) Name() string {
	return "ngo:db:metrics"
}

func (p *gormMetricsPlugin) Initialize(db *gorm.DB) error {
	p.registerCallbacks(db)
	return nil
}

func (p *gormMetricsPlugin) registerCallbacks(db *gorm.DB) {

	db.Callback().Query().Before("gorm:query").Register("ngo:metrics:before_query", p.metricBefore)
	db.Callback().Query().After("gorm:query").Register("ngo:metrics:after_query", p.metricAfter)

	db.Callback().Create().Before("gorm:create").Register("ngo:metrics:before_create", p.metricBefore)
	db.Callback().Create().After("gorm:create").Register("ngo:metrics:after_create", p.metricAfter)

	db.Callback().Update().Before("gorm:update").Register("ngo:metrics:before_update", p.metricBefore)
	db.Callback().Update().After("gorm:update").Register("ngo:metrics:after_update", p.metricAfter)

	db.Callback().Delete().Before("gorm:delete").Register("ngo:metrics:before_delete", p.metricBefore)
	db.Callback().Delete().After("gorm:delete").Register("ngo:metrics:after_delete", p.metricAfter)

	db.Callback().Row().Before("gorm:row").Register("ngo:metrics:before_row", p.metricBefore)
	db.Callback().Row().After("gorm:row").Register("ngo:metrics:after_row", p.metricAfter)

	db.Callback().Raw().Before("gorm:raw").Register("ngo:metrics:before_raw", p.metricBefore)
	db.Callback().Raw().After("gorm:raw").Register("ngo:metrics:after_raw", p.metricAfter)
}

//region callbacks

func (p *gormMetricsPlugin) metricBefore(db *gorm.DB) {

}

func (p *gormMetricsPlugin) metricAfter(db *gorm.DB) {

}

//endregion
