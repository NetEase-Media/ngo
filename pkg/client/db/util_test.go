package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqlFilter(t *testing.T) {
	var sql string

	sql = sqlFilter("select * from test where id in (?,?,?,?)")
	assert.Equal(t, "select * from test where id in (?)", sql)

	sql = sqlFilter("select * from test where id IN (?,?,?,?)")
	assert.Equal(t, "select * from test where id IN (?)", sql)

	sql = sqlFilter("select * from test where id IN(? , ? ,          ?,?)")
	assert.Equal(t, "select * from test where id IN (?)", sql)

	sql = sqlFilter("select * from test where id IN(? , ? ,          ?,?) and name In( ? ,? ,?)")
	assert.Equal(t, "select * from test where id IN (?) and name In (?)", sql)

	sql = sqlFilter("select * from test where id not in (?,?,?,?)")
	assert.Equal(t, "select * from test where id not in (?)", sql)

	sql = sqlFilter("select * from test where id not iN (?,?,?,?)")
	assert.Equal(t, "select * from test where id not iN (?)", sql)

	sql = sqlFilter("         insert into test in values (?,?,?,?)         ")
	assert.Equal(t, "insert into test in values (?,?,?,?)", sql)
}
