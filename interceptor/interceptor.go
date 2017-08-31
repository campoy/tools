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

// Package interceptor provides an easy way to intercept HTTP requests before
// being sent to remote servers.
// It benefits of the RoundTripper interface and the well known Gorilla web
// toolkit to provide an easy way to mock calls by definining routes on a
// mux.Router.
package interceptor

import (
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

// New creates a new HTTP interceptor that uses the given transport
// as a back up. If the given transport is nil the DefaultTransport
// will be used.
func New(t http.RoundTripper) *Interceptor {
	return &Interceptor{t, mux.NewRouter()}
}

// An Interceptor is both an http RoundTripper and a mux Router from Gorilla.
// When used as a transport all HTTP requests will be handled first by the
// router, then send to the back up transport if the none of the routes match
// (would return a 404).
type Interceptor struct {
	rt http.RoundTripper
	*mux.Router
}

// Client returns a new HTTP client that will intercept http requests.
func (i *Interceptor) Client() *http.Client {
	return &http.Client{Transport: i}
}

// RoundTrip handles the http request first with the local handler.
func (i *Interceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	i.ServeHTTP(w, req)
	if w.Code == http.StatusNotFound {
		if i.rt != nil {
			return i.rt.RoundTrip(req)
		}
		return http.DefaultTransport.RoundTrip(req)
	}
	return w.Result(), nil
}
