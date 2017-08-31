package main

import (
	"fmt"
	"io"
	"os"

	"github.com/campoy/tools/imgcat"
	"github.com/pkg/errors"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage:\n\t%s [image_path]*\n", os.Args[0])
		os.Exit(1)
	}

	for _, path := range os.Args[1:] {
		if err := cat(path); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}

func cat(path string) error {
	w, err := imgcat.New(os.Stdout, imgcat.Inline(true))
	if err != nil {
		return err
	}
	defer w.Close()

	f, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "could not open %s", path)
	}
	defer f.Close()
	_, err = io.Copy(w, f)
	return err
}
