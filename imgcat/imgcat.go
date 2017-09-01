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
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// An Option modifies how an image is displayed.
type Option string

// Length is used by the Width and Height options.
type Length string

// Cells gives a length in character cells.
func Cells(x int) Length { return Length(fmt.Sprint(x)) }

// Pixels gives a length in pixels.
func Pixels(x int) Length { return Length(fmt.Sprintf("%dpx", x)) }

// Percent gives relative length to the session's width or height.
func Percent(x int) Length { return Length(fmt.Sprintf("%d%%", x)) }

// Auto keeps the image to its inherent size.
func Auto() Length { return Length("auto") }

// Name sents the filename for the image. Defaults to "Unnamed file".
func Name(name string) Option {
	buf := new(bytes.Buffer)
	enc := base64.NewEncoder(base64.StdEncoding, buf)
	fmt.Fprint(enc, name)
	enc.Close()
	return Option(fmt.Sprintf("name=%s", buf))
}

// Size sets the file size in bytes. It's only used by the progress indicator.
func Size(size int) Option {
	return Option(fmt.Sprintf("size=%d", size))
}

// Width to render, it can be in cells, pixels, percentage, or auto.
func Width(l Length) Option {
	return Option(fmt.Sprintf("width=%s", l))
}

// Height to render, it can be in cells, pixels, percentage, or auto.
func Height(l Length) Option {
	return Option(fmt.Sprintf("height=%s", l))
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// PreserveAspectRatio set to false causes the image's inherent
// aspect ratio will not be respected; otherwise, it will fill the
// specified width and height as much as possible without stretching.
// Defaults to true.
func PreserveAspectRatio(b bool) Option {
	return Option(fmt.Sprintf("preserveAspectRatio=%d", boolToInt(b)))
}

// Inline set to true causes the to be displayed inline.
// Otherwise, it will be downloaded with no visual
// representation in the terminal session.
// Defaults to false.
func Inline(b bool) Option {
	return Option(fmt.Sprintf("inline=%d", boolToInt(b)))
}

// IsSupported check whether imgcat works in the current terminal.
func IsSupported() bool { return isSupported() }

// Can be swapped for testing.
var isSupported = func() bool {
	return os.Getenv("TERM_PROGRAM") == "iTerm.app"
}

// New returns a writer that encodes an image for iterm2.
func New(w io.Writer, options ...Option) (io.WriteCloser, error) {
	if !IsSupported() {
		return nil, fmt.Errorf("imgcat is only supported with iTerm2")
	}

	pr, pw := io.Pipe()
	res := &writer{
		out:        w,
		encoder:    base64.NewEncoder(base64.StdEncoding, pw),
		pipeWriter: pw,
		done:       make(chan struct{}),
	}
	go res.writeFrom(pr, options...)

	return res, nil
}

type writer struct {
	out io.Writer

	encoder    io.WriteCloser
	pipeWriter *io.PipeWriter

	done chan struct{}
}

func (w *writer) Write(p []byte) (int, error) { return w.encoder.Write(p) }

func (w *writer) Close() error {
	if err := w.encoder.Close(); err != nil {
		return fmt.Errorf("could not close base64 encoder: %v", err)
	}
	if err := w.pipeWriter.Close(); err != nil {
		return fmt.Errorf("could not close pipe writer: %v", err)
	}
	<-w.done
	return nil
}

func (w *writer) writeFrom(pr *io.PipeReader, options ...Option) {
	defer close(w.done)

	var buf bytes.Buffer
	fmt.Fprint(&buf, "\x1b]1337;File=")
	for i, option := range options {
		fmt.Fprintf(&buf, "%s", option)
		if i < len(options)-1 {
			fmt.Fprintf(&buf, ";")
		}
	}
	fmt.Fprintf(&buf, ":")

	_, err := io.Copy(w.out, io.MultiReader(
		&buf, pr, bytes.NewBufferString("\a\n"),
	))
	pr.CloseWithError(err)
}
