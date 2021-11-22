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

const ngoTracerSpanKey = "ngo:tracer:span"

type gormTracerPlugin struct{}

func newGormTracerPlugin() *gormTracerPlugin {
	return &gormTracerPlugin{}
}

func (t *gormTracerPlugin) Name() string {
	return "ngo:tracer"
}

func (t *gormTracerPlugin) Initialize(db *gorm.DB) error {
	t.registerCallbacks(db)
	return nil
}

func (t *gormTracerPlugin) registerCallbacks(db *gorm.DB) {

	db.Callback().Query().Before("gorm:query").Register("ngo:tracer:before_query", t.traceBefore)
	db.Callback().Query().After("gorm:query").Register("ngo:tracer:after_query", t.traceAfter)

	db.Callback().Create().Before("gorm:create").Register("ngo:tracer:before_create", t.traceBefore)
	db.Callback().Create().After("gorm:create").Register("ngo:tracer:after_create", t.traceAfter)

	db.Callback().Update().Before("gorm:update").Register("ngo:tracer:before_update", t.traceBefore)
	db.Callback().Update().After("gorm:update").Register("ngo:tracer:after_update", t.traceAfter)

	db.Callback().Delete().Before("gorm:delete").Register("ngo:tracer:before_delete", t.traceBefore)
	db.Callback().Delete().After("gorm:delete").Register("ngo:tracer:after_delete", t.traceAfter)

	db.Callback().Row().Before("gorm:row").Register("ngo:tracer:before_row", t.traceBefore)
	db.Callback().Row().After("gorm:row").Register("ngo:tracer:after_row", t.traceAfter)

	db.Callback().Raw().Before("gorm:raw").Register("ngo:tracer:before_raw", t.traceBefore)
	db.Callback().Raw().After("gorm:raw").Register("ngo:tracer:after_raw", t.traceAfter)
}

//region callbacks

func (p *gormTracerPlugin) traceBefore(db *gorm.DB) {

}

func (t *gormTracerPlugin) traceAfter(db *gorm.DB) {

}

//endregion
