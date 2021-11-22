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

package httplib

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func ExampleGet() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello, client")
	}))
	defer ts.Close()

	Get(ts.URL).Do(context.Background())

	// Output:
	// Hello, client
}

func ExampleGet_query() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Query().Encode())
		fmt.Println(r.Header.Get("age"))
	}))
	defer ts.Close()

	Get(ts.URL).AddQuery("cat", "maimai").AddHeaderKV("age", "10").Do(context.Background())

	// Output:
	// cat=maimai
	// 10
}

func ExamplePost() {
	type testBody struct {
		I int
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		var reqBody testBody
		json.Unmarshal(b, &reqBody)
		reqBody.I += 10

		b, _ = json.Marshal(&reqBody)
		w.Header().Set("cat", "maimai")
		w.Write(b)
	}))
	defer ts.Close()

	var resBody testBody
	h := make(H)
	Post(ts.URL).SetJson(&testBody{I: 10}).BindHeader(h).BindJson(&resBody).Do(context.Background())
	fmt.Println(resBody.I)
	fmt.Println(h.Get("cat"))

	// Output:
	// 20
	// maimai
}
