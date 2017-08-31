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

package interceptor_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/campoy/tools/interceptor"
)

func ExampleNew() {
	// Let's create a new Interceptor without back up Transport.
	i := interceptor.New(nil)
	// All requests to google.com will be intercepted.
	i.HandleFunc("/{p}", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "requests to Google are intercepted", http.StatusUnauthorized)
	}).Host("google.com")

	// Now we can create a client and use it normally.
	client := i.Client()
	res, err := client.Get("http://google.com/hello")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("response with status", res.Status)
	// Output:
	// response with status Unauthorized
}
