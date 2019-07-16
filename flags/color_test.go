// Copyright 2016 Google Inc. All rights reserved.
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

package flags

import (
	"fmt"
	"image/color"
	"testing"
)

type rgba struct{ r, g, b, a uint32 }

func (c rgba) String() string {
	return fmt.Sprintf("rgba(%d, %d, %d, %d)", c.r, c.g, c.b, c.a)
}

func TestHexColor(t *testing.T) {
	tests := []struct {
		text    string
		parsed  color.Color
		invalid bool
	}{
		{text: "ffffff", parsed: color.White},
		{text: "FFFFFF", parsed: color.White},
		{text: "000000", parsed: color.Black},
		{text: "#000000", parsed: color.Black},
		{text: "#111", parsed: color.RGBA{16, 16, 16, 255}},
		{text: "eeeeee", parsed: color.RGBA{238, 238, 238, 255}},
		{text: "010203", parsed: color.RGBA{1, 2, 3, 255}},
		{text: "not a number", invalid: true},
		{text: "deadbeef", invalid: true}, // too long
		{text: "abcd", invalid: true},     // too short
	}

	for _, tt := range tests {
		var c hexColor
		if err := c.Set(tt.text); err != nil {
			if !tt.invalid {
				t.Errorf("parsing %s failed unexpectedly: %v", tt.text, err)
			}
			continue
		}
		if tt.invalid {
			t.Errorf("parsing %s should have failed", tt.text)
			continue
		}

		r, g, b, a := c.RGBA()
		got := rgba{r, g, b, a}
		r, g, b, a = tt.parsed.RGBA()
		wants := rgba{r, g, b, a}
		if got != wants {
			t.Errorf("%s should be parsed as %v; got %v", tt.text, wants, got)
		}
	}

}
