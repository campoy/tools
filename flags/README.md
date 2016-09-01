flags
=====

Package flags implements command-line flag parsing for types
that are not part of the standard library flag package, but I
consider to be useful.

[docs](http://godoc.org/github.com/campoy/tools/flags)

The available flags are:

## HexColor flags

You can use them like this:

[embedmd]:# (samples/color/main.go /package main/ $)
```go
package main

import (
	"flag"
	"fmt"
	"image/color"

	"github.com/campoy/tools/flags"
)

func main() {
	a := flags.HexColor("a", color.White, "color a")
	var b color.Color
	flags.HexColorVar(&b, "b", color.White, "color b")
	flag.Parse()

	fmt.Printf("a is %v; b is %v", a, b)
}
```

## TextTemplate flags

You can use them like this:

[embedmd]:# (samples/texttemplate/main.go /package main/ $)
```go
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
```

### Disclaimer

This is not an official Google product (experimental or otherwise), it is just code that happens to be owned by Google.
