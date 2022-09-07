package db

import "fmt"

var (
	_ error = NoSuchDBError{}
)

// NoSuchDBError 错误表示找不到对应name的db client
type NoSuchDBError struct {
	DBName string
}

func (err NoSuchDBError) Error() string {
	return fmt.Sprintf("can't find db named %s", err.DBName)
}
