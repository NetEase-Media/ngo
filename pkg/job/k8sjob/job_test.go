package k8sjob

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJob(t *testing.T) {
	opt := &Options{
		CenterUrl: "",
		Namespace: "testns",
	}

	argsArr := []struct {
		hostname string
		f        Callback
		err      error
		panic    interface{}
	}{
		{
			hostname: "testhost1",
			f: func(a *Args) (string, error) {
				return strconv.Itoa(a.SharedNum), nil
			},
		},
		{
			hostname: "testhost2",
			f: func(a *Args) (string, error) {
				return strconv.Itoa(a.SharedNum), nil
			},
		},
		{
			hostname: "testhost3",
			err:      errors.New("fake error"),
			f: func(a *Args) (string, error) {
				return "", errors.New("fake error")
			},
		},
		{
			hostname: "testhost4",
			f: func(a *Args) (string, error) {
				panic("fake panic")
			},
			panic: "fake panic",
		},
	}

	for i := range argsArr {
		args := &argsArr[i]
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				assert.Equal(t, fmt.Sprintf("/api/v1/%s/podshared/report", opt.Namespace), r.URL.Path)
				b, err := ioutil.ReadAll(r.Body)
				assert.Nil(t, err)
				r.Body.Close()

				var req ReportRequest
				err = json.Unmarshal(b, &req)
				assert.Nil(t, err)
				assert.Equal(t, args.hostname, req.PodName)
				if args.panic != nil {
					assert.Equal(t, fmt.Sprintf("%v", args.panic), req.Exception)
				} else if args.err != nil {
					assert.Equal(t, args.err.Error(), req.Exception)
				} else {
					assert.Empty(t, req.Exception)
					assert.Equal(t, strconv.Itoa(i), req.Result)
				}

				b, err = json.Marshal(&ReportResponse{Status: 200})
				assert.Nil(t, err)
				_, err = w.Write(b)
				assert.Nil(t, err)

			case http.MethodGet:
				assert.Equal(t, fmt.Sprintf("/api/v1/%s/podshared/getSharedNum", opt.Namespace), r.URL.Path)
				assert.Equal(t, args.hostname, r.URL.Query().Get("podName"))

				res := &GetShareNumResponse{
					Status:  0,
					Message: "",
					Data: GetShareNumData{
						SharedNum: i,
					},
				}
				b, err := json.Marshal(res)
				assert.Nil(t, err)
				_, err = w.Write(b)
				assert.Nil(t, err)
			default:
				assert.FailNow(t, "wrong method", r.Method)
			}
		}))

		opt.CenterUrl = s.URL

		j := &job{
			opt:      opt,
			f:        args.f,
			hostname: args.hostname,
		}
		j.run()

		s.Close()
	}
}

func testRun(t *testing.T, hostname string, f Callback, err error, panic interface{}) {

}
