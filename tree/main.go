// Copyright 2014 Google Inc. All rights reserved.
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

// tree is a very simple implementation of the tree unix command.
// This implementation doesn't provide any options as flags.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type caseInsensitive struct {
	values []string
}

func (ci caseInsensitive) Len() int {
	return len(ci.values)
}

func (ci caseInsensitive) Less(i, j int) bool {
	return strings.ToLower(ci.values[i]) < strings.ToLower(ci.values[j])
}

func (ci caseInsensitive) Swap(i, j int) {
	ci.values[i], ci.values[j] = ci.values[j], ci.values[i]
}

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("absolute %s: %v", path, err)
	}
	dirs, files, err := visit(path, "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n%v directories, %v files\n", dirs, files)
}

func visit(path, indent string) (dirs, files int, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, 0, fmt.Errorf("stat %s: %v", path, err)
	}
	fmt.Println(fi.Name())
	if !fi.IsDir() {
		return 0, 1, nil
	}

	dir, err := os.Open(path)
	if err != nil {
		return 1, 0, fmt.Errorf("open %s: %v", path, err)
	}
	names, err := dir.Readdirnames(-1)
	dir.Close()
	if err != nil {
		return 1, 0, fmt.Errorf("read dir names %s: %v", path, err)
	}
	names = removeHidden(names)
	sort.Sort(caseInsensitive{names})
	add := "│   "
	for i, name := range names {
		if i == len(names)-1 {
			fmt.Printf(indent + "└── ")
			add = "    "
		} else {
			fmt.Printf(indent + "├── ")
		}
		d, f, err := visit(filepath.Join(path, name), indent+add)
		if err != nil {
			log.Println(err)
		}
		dirs, files = dirs+d, files+f
	}
	return dirs + 1, files, nil
}

func removeHidden(files []string) []string {
	var clean []string
	for _, f := range files {
		if f[0] != '.' {
			clean = append(clean, f)
		}
	}
	return clean
}
