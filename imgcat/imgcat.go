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

// Package imgcat provides a writer useful to show images directly into iterm2.
package imgcat

import (
	"encoding/base64"
	"fmt"
	"io"
)

// New returns a writer that encodes images and for iterm2.
func New(w io.Writer, width string) io.WriteCloser {
	pr, pw := io.Pipe()
	enc := base64.NewEncoder(base64.StdEncoding, pw)

	res := writer{enc, pr, make(chan struct{})}
	go func() {
		fmt.Fprintf(w, "\x1b]1337;File=inline=1;width=%s:", width)
		io.Copy(w, pr)
		fmt.Fprintf(w, "\a\n")
		close(res.done)
	}()

	return res
}

type writer struct {
	io.WriteCloser
	c    io.Closer
	done chan struct{}
}

func (w writer) Close() error {
	w.c.Close()
	<-w.done
	return w.WriteCloser.Close()
}
