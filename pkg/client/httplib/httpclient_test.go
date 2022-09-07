package httplib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/djimenez/iconv-go"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	var body []byte
	c, err := New(&Options{})
	_, err = c.Get("http://www.163.com").BindBytes(&body).doInternal()
	assert.Nil(t, err)
}

func TestHttpClientGet(t *testing.T) {
	body := &testJsonBody{
		A: "fdstt",
		B: 5323,
		C: 43.54,
	}
	b, err := json.Marshal(body)
	assert.Nil(t, err)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err = w.Write(b)
		assert.Nil(t, err)
	}))
	defer s.Close()

	c, err := New(&Options{})

	var obj testJsonBody
	statusCode, err := c.Get(s.URL).BindJson(&obj).doInternal()
	assert.Nil(t, err)
	assert.EqualValues(t, body, &obj)
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestHttpClientPost(t *testing.T) {
	body := &testJsonBody{
		A: "fdstt",
		B: 5323,
		C: 43.54,
	}
	b, err := json.Marshal(body)
	assert.Nil(t, err)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		_, err = w.Write(b)
		assert.Nil(t, err)
	}))
	defer s.Close()

	c, err := New(&Options{})

	var obj testJsonBody
	statusCode, err := c.Post(s.URL).BindJson(&obj).doInternal()
	assert.Nil(t, err)
	assert.EqualValues(t, body, &obj)
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestInit(t *testing.T) {
	body := []byte("ok")
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(body)
		assert.Nil(t, err)
	}))
	defer s.Close()

	defaultHttpClient = nil
	var res1 []byte
	client, err := New(&Options{})
	assert.Nil(t, err)
	statusCode, err := client.Get(s.URL).BindBytes(&res1).doInternal()
	assert.EqualValues(t, body, res1)
	assert.Equal(t, http.StatusOK, statusCode)

	var res2 []byte
	c, err := New(&Options{})
	statusCode, err = c.Get(s.URL).BindBytes(&res2).doInternal()
	assert.Nil(t, err)
	assert.EqualValues(t, body, res2)
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestCharset(t *testing.T) {
	body := []byte("{\"a\": \"成功\"}")
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.URL.Query()["charset"]
		if c != nil || len(c) > 0 {
			w.Header().Set("Content-Type", "application/json;charset="+c[0])
			output, err := iconv.ConvertString(string(body), "utf-8", c[0])
			assert.Nil(t, err)
			_, err = w.Write([]byte(output))
			assert.Nil(t, err)
		} else {
			_, err := w.Write(body)
			assert.Nil(t, err)
		}
	}))
	defer s.Close()

	c, err := New(&Options{})

	var res0 string
	statusCode, err := c.Get(s.URL).BindString(&res0).doInternal()
	assert.NoError(t, err)
	assert.EqualValues(t, string(body), res0)
	assert.Equal(t, http.StatusOK, statusCode)

	var res1 string
	statusCode, err = c.Get(s.URL + "?charset=gbk").BindString(&res1).doInternal()
	assert.NoError(t, err)
	assert.EqualValues(t, string(body), res1)
	assert.Equal(t, http.StatusOK, statusCode)

	var res2 testJsonBody
	statusCode, err = c.Get(s.URL + "?charset=gbk").BindJson(&res2).doInternal()
	assert.NoError(t, err)
	assert.EqualValues(t, "成功", res2.A)
	assert.Equal(t, http.StatusOK, statusCode)

}

func TestCharset_gbk(t *testing.T) {
	var res string
	url := "http://help.3g.163.com/special/supervision_level/"
	c, err := New(&Options{})
	statusCode, err := c.Get(url).BindString(&res).doInternal()
	assert.NoError(t, err)
	assert.Contains(t, res, "习近平")
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestCharset_utf8(t *testing.T) {
	var res string
	url := "http://test-elephant.hz.service.163.org/api/v1/elephant/doc/retrieve.do?idc=dy&needMobileBody=1&needPcBody=1&docid=F863DLQD000198G3&productId=OPEN&needOrigimg=1"
	c, err := New(&Options{})
	statusCode, err := c.Get(url).BindString(&res).doInternal()
	assert.NoError(t, err)
	assert.Contains(t, res, "这位同学我要测试延伸阅读了哈")
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestCharset_String2JsonErr(t *testing.T) {
	url := "http://test-elephant.hz.service.163.org/api/v1/elephant/doc/retrieve.do?idc=dy&needMobileBody=1&needPcBody=1&docid=F863DLQD000198G3&productId=OPEN&needOrigimg=1"
	c, err := New(&Options{})
	var res1 string
	_, err = c.Get(url).BindJson(&res1).doInternal()
	assert.Error(t, err)
}
