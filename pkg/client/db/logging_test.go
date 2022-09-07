package db

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestLogSQL(t *testing.T) {
	setupTest()
	testSqlmock.ExpectQuery("SELECT (.+) FROM `testusers`").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "gender"}).FromCSVString("1,la,male"))
	var u testuser
	WithContext(context.Background(), testClient).First(&u)
	assert.Equal(t, u.ID, "1")
	assert.Equal(t, u.Name, "la")
	assert.Equal(t, u.Gender, "male")
	tearDownTest()
}
