# Project Euler

This repo has my solutions for [Project Euler](https://projecteuler.net/) in Go.

## Viewing Code

To view the solutions you must perform the following steps:
- run `./setup.sh` to setup the git filters
- create `answerkey.toml` (see Answer Key Format below)
- re-checkout the Go files (run `rm *.go` and then `git checkout .`)

## Answer Key Format

The `answerkey.toml` file should be in the root of the repo:

```
go-project-euler/
  main.go
  encrypt.go
  p1.go
  answerkey.toml
```

The contents should look like this:

```toml
[Answers]
1 = "answer-here"
2 = "answer-here"
```

