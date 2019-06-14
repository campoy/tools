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

package httplog_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/campoy/tools/v2/httplog"
)

var noLog = func(format string, vs ...interface{}) {}

func TestRequestDoesntChange(t *testing.T) {
	const reqBody = "foo bar"

	errc := make(chan error, 1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errc <- func() error {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return err
			}
			if reqBody != string(b) {
				return fmt.Errorf("expected request body %q, got %q", reqBody, b)
			}
			return nil
		}()
	}))
	defer ts.Close()

	c := httplog.NewTransport(nil, false, noLog).Client()
	if _, err := c.Post(ts.URL, "text", strings.NewReader(reqBody)); err != nil {
		t.Fatalf("get: %v", err)
	}
	if err := <-errc; err != nil {
		t.Fatal(err)
	}
}

func TestResponseDoesntChange(t *testing.T) {
	const resBody = "foo bar"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, resBody)
	}))
	defer ts.Close()

	c := httplog.NewTransport(nil, false, noLog).Client()
	res, err := c.Get(ts.URL)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	if resBody != string(b) {
		t.Fatalf("expected request response %q, got %q", resBody, b)
	}
}

func TestMessagesAreLogged(t *testing.T) {
	const (
		start         = "^httplog: "
		end           = "(\r\n)*$"
		get           = "GET / HTTP/1.1\r\n"
		post          = "POST / HTTP/1.1\r\n"
		host          = "Host: 127.0.0.1:[0-9]+\r\n"
		contentType   = "Content-Type: text/plain(; charset=utf-8)?\r\n"
		contentLength = "Content-Length: [0-9]+\r\n"
		date          = "Date: [A-Za-z]+, [0-9]+ [A-Za-z]+ [0-9]+ ([0-9]+:?)+ GMT\r\n"
		ok            = "HTTP/1.1 200 OK\r\n"
	)

	var tests = []struct {
		desc                 string
		logBody              bool
		reqBody, resBody     io.Reader
		reqRegexp, resRegexp string
	}{
		{
			"request and response empty and body not logged",
			false, nil, nil,
			start + get + host + end,
			start + ok + "(" + contentType + date + contentLength + ")?" + end,
		},
		{
			"request not empty, response empty and body not logged",
			false, strings.NewReader("foo"), nil,
			start + post + host + contentType + end,
			start + ok + "(" + contentLength + contentType + date + ")?" + end,
		},
		{
			"request empty, response not empty, and body not logged",
			false, nil, strings.NewReader("bar"),
			start + get + host + end,
			start + ok + contentLength + contentType + date + end,
		},
		{
			"request and response not empty and body not logged",
			false, strings.NewReader("foo"), strings.NewReader("bar"),
			start + post + host + contentType + end,
			start + ok + contentLength + contentType + date + end,
		},
		{
			"request and response empty and body logged",
			true, nil, nil,
			start + get + host + end,
			start + ok + contentType + date + contentLength + end,
		},
		{
			"request not empty, response empty and body logged",
			true, strings.NewReader("foo"), nil,
			start + post + host + contentType + "\r\nfoo" + end,
			start + ok + contentLength + contentType + date + end,
		},
		{
			"request empty, response not empty, and body logged",
			true, nil, strings.NewReader("bar"),
			start + get + host + end,
			start + ok + contentLength + contentType + date + "\r\nbar" + end,
		},
		{
			"request and response not empty and body logged",
			true, strings.NewReader("foo"), strings.NewReader("bar"),
			start + post + host + contentType + "\r\nfoo" + end,
			start + ok + contentLength + contentType + date + "\r\nbar" + end,
		},
	}

	for i, test := range tests {
		var logs []string
		log := func(format string, vs ...interface{}) {
			logs = append(logs, fmt.Sprintf(format, vs...))
		}

		err := func() error {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if test.resBody != nil {
					_, err := io.Copy(w, test.resBody)
					checkError(t, err)
				}
			}))
			defer ts.Close()

			c := httplog.NewTransport(nil, test.logBody, log).Client()
			if test.reqBody == nil {
				_, err := c.Get(ts.URL)
				checkError(t, err)
			} else {
				_, err := c.Post(ts.URL, "text/plain", test.reqBody)
				checkError(t, err)
			}

			if len(logs) != 2 {
				return fmt.Errorf("expected two logs, got %v: %#v", len(logs), logs)
			}
			if !regexp.MustCompile(test.reqRegexp).MatchString(logs[0]) {
				return fmt.Errorf("bad request log:\n%q\ndoesn't match\n%q", logs[0], test.reqRegexp)
			}
			if !regexp.MustCompile(test.resRegexp).MatchString(logs[1]) {
				return fmt.Errorf("bad response log:\n%q\ndoesn't match\n%q", logs[1], test.resRegexp)
			}
			return nil
		}()

		if err != nil {
			t.Errorf("test case %v %s failed: %v", i, test.desc, err)
		}
	}
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
