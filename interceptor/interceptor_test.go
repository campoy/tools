// Copyright 2014 Google Inc. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

package interceptor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

func TestIntercept(t *testing.T) {
	tests := []struct {
		url, out string
	}{
		{"https://google.com/hello", "hello, there"},                                    // intercepted
		{"https://google.com/foo", "There was no service found for the uri requested."}, // not intercepted

		{"https://api.intercepted.com/say/hi", "hi"},                                       // intercepted
		{"https://google.com/say/hi", "There was no service found for the uri requested."}, // not intercepted
	}

	i := New(nil)
	i.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello, there")
	})
	i.HandleFunc("/say/{msg}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mux.Vars(r)["msg"])
	}).Host("api.intercepted.com")

	client := i.Client()
	for _, tt := range tests {
		res, err := client.Get(tt.url)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("could not read response: %v", err)
			continue
		}
		res.Body.Close()
		if got := string(b); got != tt.out {
			t.Errorf("expected body %s; got %s", tt.out, got)
		}
	}

}
