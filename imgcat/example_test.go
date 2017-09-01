// Copyright 2017 Google Inc. All rights reserved.
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

package imgcat_test

import (
	"log"
	"os"

	"github.com/campoy/tools/imgcat"
)

func ExampleNewEncoder() {
	enc, err := imgcat.NewEncoder(os.Stdout, imgcat.Width(imgcat.Pixels(100)), imgcat.Inline(true), imgcat.Name("smiley.png"))
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open("testdata/icon.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = f.Close() }()

	// Display the image in the terminal.
	if err := enc.Encode(f); err != nil {
		log.Fatal(err)
	}
}
