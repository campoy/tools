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

package main

import (
	"flag"
	"fmt"
	"image/color"

	"github.com/campoy/tools/v2/flags"
)

func main() {
	a := flags.HexColor("a", color.White, "color a")
	var b color.Color
	flags.HexColorVar(&b, "b", color.White, "color b")
	flag.Parse()

	fmt.Printf("a is %v; b is %v", a, b)
}
