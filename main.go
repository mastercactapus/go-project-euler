package main

// encrypt = ignore

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/burntsushi/toml"
)

type answerKey struct {
	Answers map[string]string
}

var f = make(map[int]func() string)

func register(num int, fn func() string) {
	f[num] = fn
}

func updateKey(key, val string) {
	var a answerKey
	_, err := toml.DecodeFile("answerkey.toml", &a)
	if err != nil {
		fmt.Printf("failed to update answer key: %v\n", err)
		return
	}
	a.Answers[key] = val
	fd, err := os.Create("answerkey.toml")
	if err != nil {
		fmt.Printf("failed to update answer key: %v\n", err)
		return
	}
	defer fd.Close()
	err = toml.NewEncoder(fd).Encode(a)
	if err != nil {
		fmt.Printf("failed to update answer key: %v\n", err)
		return
	}
}

func main() {
	n := flag.Int("n", 0, "Problem number to solve.")
	g := flag.Int("gen", 0, "Generate new problem file.")
	enc := flag.Bool("clean", false, "Operate in 'clean' mode (should only be used by git)")
	dec := flag.Bool("smudge", false, "Operate in 'smudge' mode (should only be used by git)")

	flag.Parse()

	if *enc {
		clean()
		return
	} else if *dec {
		smudge()
		return
	}
	if *g > 0 {
		gen(*g)
		return
	}

	if *n < 1 {
		fmt.Println("ERROR: n must be >= 1")
		os.Exit(1)
	}

	if f[*n] == nil {
		fmt.Printf("ERROR: no solution for #%d\n", *n)
		os.Exit(1)
	}
	s := time.Now()
	result := f[*n]()
	fmt.Println("Elapsed: ", time.Since(s).String())
	fmt.Println("Solution:", result)
	updateKey(strconv.Itoa(*n), result)

}
