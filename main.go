// Copyright Mike Hughes 2021 (mike AT mikehughes DOT info)
//
// Utility wrapping Google's [licensecheck](https://github.com/google/licensecheck) package.
//
// Accepts a single argument, either a local file or a URL.
//
//     lc ./PATH/TO/LICENSE
//
// or
//
//     lc https://HOSTNAME/PATH/TO/LICENSE
//
// example output:
//
//     98.8% of text covered by licenses:
//     MIT at [13:1068]
//
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/google/licensecheck"
)

func main() {
	args := os.Args
	if len(args) < 2 || len(args) > 2 {
		log.Fatal("Please provide a single filename or URL to check\n")
	}
	var (
		content []byte
		err     error
	)
	switch {
	case IsUrl(args[1]):
		resp, err := http.Get(args[1])
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		content, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
	default:
		content, err = ioutil.ReadFile(args[1])
		if err != nil {
			log.Fatalf("%v\n", err)
		}

	}
	cov := licensecheck.Scan(content)
	fmt.Printf("%.1f%% of text covered by licenses:\n", cov.Percent)
	for _, m := range cov.Match {
		fmt.Printf("%s at [%d:%d]\n", m.ID, m.Start, m.End)
	}
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
