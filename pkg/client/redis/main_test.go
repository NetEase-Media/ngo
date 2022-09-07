package redis

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setupTest()
	ret := m.Run()
	tearDownTest()
	os.Exit(ret)
}

func setupTest() {
}

func tearDownTest() {

}
