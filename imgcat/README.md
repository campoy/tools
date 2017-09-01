imgcat
======

imgcat provides a convenient way to print images into iTerm2.

[docs](http://godoc.org/github.com/campoy/tools/imgcat)

## Example

[embedmd]:# (imgcat/main.go /package main/ $)
```go
package main

import (
	"fmt"
	"os"

	"github.com/campoy/tools/imgcat"
	"github.com/pkg/errors"
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
```

### Disclaimer

This is not an official Google product (experimental or otherwise), it is just code that happens to be owned by Google.
