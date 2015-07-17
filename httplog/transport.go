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

// httplog provides an implementation of http.RoundTripper that logs every
// single request and response using a given logging function.
package httplog

import (
	"log"
	"net/http"
	"net/http/httputil"
)

// Transport satisfies http.RoundTripper
type Transport struct {
	transport http.RoundTripper
	// Should the body of the requests and responses be logged.
	logBody bool
	// If logf is nil log.Printf will be used.
	logf func(format string, vs ...interface{})
}

// NewTransport returns a new Transport that uses the given RoundTripper, or
// http.DefaultTransport if nil, and logs all requests and responses using
// logf, or log.Printf if nil.
// The body of the requests and responses are logged too only if logBody is
// true.
func NewTransport(rt http.RoundTripper, logBody bool, logf func(string, ...interface{})) Transport {
	if rt == nil {
		rt = http.DefaultTransport
	}
	if logf == nil {
		logf = log.Printf
	}
	return Transport{rt, logBody, logf}
}

// Client returns a new http.Client using the given transport.
func (t Transport) Client() http.Client { return http.Client{Transport: t} }

// Transport satifies http.RoundTripper
func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if b, err := httputil.DumpRequest(req, t.logBody); err != nil {
		t.logf("httplog: dump request: %v", err)
		return nil, err
	} else {
		t.logf("httplog: %s", b)
	}

	res, err := t.transport.RoundTrip(req)
	if err != nil {
		t.logf("httplog: roundtrip error: %v", err)
		return res, err
	}

	if b, err := httputil.DumpResponse(res, t.logBody); err != nil {
		t.logf("httplog: dump response: %v", err)
	} else {
		t.logf("httplog: %s", b)
	}
	return res, err
}
