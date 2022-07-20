package apollo

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/NetEase-Media/ngo/pkg/config/remote/apollo/mockserver"
	"github.com/philchia/agollo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}

func setup() {
	go func() {
		if err := mockserver.Run(); err != nil {
			log.Println(err)
		}
	}()
	// wait for mock server to run
	time.Sleep(time.Second)
}

func teardown() {
	mockserver.Close()
}

func TestReadConfig(t *testing.T) {
	t.Skip()
	testData := []string{"value1", "value2"}
	ns := []string{"application"}
	mockserver.Set("application", "key_test", testData[0])
	ds := NewDataSource(&agollo.Conf{
		AppID:          "SampleApp",
		Cluster:        "default",
		NameSpaceNames: ns,
		MetaAddr:       "localhost:16852",
		CacheDir:       ".",
	}, ns)
	value, err := ds.ReadConfig()
	assert.Nil(t, err)
	assert.Equal(t, testData[0], string(value))
	t.Logf("read: %s", value)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		mockserver.Set("application", "key_test", testData[1])
		time.Sleep(time.Second * 3)
		ds.Close()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range ds.IsConfigChanged() {
			value, err := ds.ReadConfig()
			assert.Nil(t, err)
			assert.Equal(t, testData[1], string(value))
			t.Logf("read: %s", value)
		}
	}()

	wg.Wait()
}
