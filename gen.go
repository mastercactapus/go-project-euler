package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

// encrypt = ignore

func gen(n int) {

	resp, err := http.Get(fmt.Sprintf("https://projecteuler.net/problem=%d", n))
	errCheck(err, "download problem text")

	root, err := html.Parse(resp.Body)
	errCheck(err, "parse response")

	node, ok := scrape.Find(root, scrape.ByClass("problem_content"))
	if !ok {
		fmt.Println("could not find problem text")
		os.Exit(1)
	}

	prob := strings.TrimSpace(scrape.Text(node))

	fd, err := os.Create(fmt.Sprintf("p%d.go", n))
	errCheck(err, "create file")
	defer fd.Close()

	fmt.Fprintf(fd,
		`package main

// encrypt = %d

func init() {
	register(%d, p%d)
}

/*

%s

*/

func p%d() string {

}

`, n, n, n, prob, n)

}
