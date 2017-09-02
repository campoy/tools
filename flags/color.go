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
	"flag"
	"fmt"
	"image/color"
	"strconv"
)

type hexColor struct {
	color.Color
}

func (c *hexColor) String() string {
	r, g, b, a := c.Color.RGBA()
	return fmt.Sprintf("rgba(%d, %d, %d, %d)", r, g, b, a)
}

func (c *hexColor) Set(s string) error {
	if len(s) != 6 {
		return fmt.Errorf("color should be 6 digits")
	}
	n, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return fmt.Errorf("not hexadecimal: %v", err)
	}
	c.Color = &color.RGBA{
		R: uint8(n >> 16),
		G: uint8(n >> 8),
		B: uint8(n),
		A: 0xff,
	}
	return nil
}

// HexColor defines a hex color flag with specified name, default value, and usage string.
// The return value is the address of an RGBA color variable that stores the value of the flag.
func HexColor(name string, value color.Color, usage string) color.Color {
	c := &hexColor{value}
	flag.Var(c, name, usage)
	return c
}

// HexColorVar defines a hex color flag with specified name, default value, and usage string.
// The argument c points to an RGBA color variable in which to store the value of the flag.
func HexColorVar(c *color.Color, name string, value color.Color, usage string) {
	p := &hexColor{value}
	*c = p
	flag.Var(p, name, usage)
}
