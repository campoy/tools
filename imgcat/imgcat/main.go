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

package main

import (
	"fmt"
	"github.com/campoy/tools/imgcat"
	"github.com/pkg/errors"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage:\n\t%s [image_path]*\n", os.Args[0])
		os.Exit(1)
	}

	enc, err := imgcat.NewEncoder(os.Stdout,
		imgcat.Inline(true),
		imgcat.Width(imgcat.Percent(100)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	for _, path := range os.Args[1:] {
		if err := cat(enc, path); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}

func cat(enc *imgcat.Encoder, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "could not open %s", path)
	}
	if err := enc.Encode(f); err != nil {
		return err
	}
	return f.Close()
}
