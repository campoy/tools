package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "missing user or org name")
		os.Exit(1)
	}

	name := flag.Arg(0)
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", name)

	for url != "" {
		res, err := http.Get(url)
		if err != nil {
			log.Fatalf("could not fetch repos for user %s: %v", name, err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			log.Fatalf("could not fetch repos for user %s: %s", name, res.Status)
		}

		var repos []struct{ Name string }
		err = json.NewDecoder(res.Body).Decode(&repos)
		if err != nil {
			log.Fatalf("could not decode response: %v", err)
		}
		if len(repos) == 0 {
			return
		}

		for _, repo := range repos {
			fmt.Println(repo.Name)
		}

		url = getNext(res.Header.Get("Link"))
	}
}

func getNext(s string) string {
	for _, link := range strings.Split(s, ",") {
		ps := strings.Split(link, ";")
		if len(ps) == 2 && strings.TrimSpace(ps[1]) != `rel="next"` {
			return strings.Trim(ps[0], " <>")
		}
	}
	return ""
}
