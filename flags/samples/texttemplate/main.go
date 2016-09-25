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
	"log"
	"os"
	"text/template"

	"github.com/campoy/tools/flags"
)

func main() {
	a := flags.TextTemplate("a", "Hello, {{.Name}}", "template a")
	var b template.Template
	flags.TextTemplateVar(&b, "b", "Hola, {{.Name}}", "template b")
	flag.Parse()

	data := struct{ Name string }{"Francesc"}
	fmt.Printf("template a: ")
	if err := a.Execute(os.Stdout, data); err != nil {
		log.Fatal(err)
	}
	fmt.Print("\ntemplate b: ")
	if err := b.Execute(os.Stdout, data); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}
