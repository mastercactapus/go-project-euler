package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/burntsushi/toml"
	"golang.org/x/crypto/scrypt"
)

// encrypt = ignore

const pemType = "ENCRYPTED CODE"

func getKey(n string) string {
	var a answerKey
	_, err := toml.DecodeFile("answerkey.toml", &a)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read answer key: %v\n", err)
		return ""
	}

	return a.Answers[n]
}

var encRx = regexp.MustCompile(`^encrypt\s*=\s*(\d+|ignore)$`)

func errCheck(err error, msg string) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, "ERROR:", msg+":", err)
	os.Exit(1)
}

func clean() {
	data, err := ioutil.ReadAll(os.Stdin)
	errCheck(err, "read stdin")

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "", string(data), parser.ParseComments)
	errCheck(err, "parse Go")

	var encKeyName string
	for _, c := range f.Comments {
		s := strings.TrimSpace(c.Text())
		m := encRx.FindStringSubmatch(s)
		if m != nil {
			encKeyName = m[1]
			continue
		}
	}
	if encKeyName == "ignore" {
		os.Stdout.Write(data)
		return
	}
	encKey := getKey(encKeyName)
	if encKey == "" {
		errCheck(errors.New("no answer found for #"+encKeyName), "answer key lookup")
	}

	salt := make([]byte, 32)
	_, err = rand.Read(salt)
	errCheck(err, "get random data")
	nonce := make([]byte, 12)
	_, err = rand.Read(nonce)
	errCheck(err, "get random data")

	key, err := scrypt.Key([]byte(encKey), salt, 16384, 8, 1, 32)
	errCheck(err, "calculate encryption key")

	c, err := aes.NewCipher(key)
	errCheck(err, "create new cipher")

	aesgcm, err := cipher.NewGCM(c)
	errCheck(err, "initialize cipher")

	var b pem.Block
	b.Headers = make(map[string]string)
	b.Type = pemType
	b.Headers["Problem Number"] = encKeyName
	b.Headers["Salt"] = base64.StdEncoding.EncodeToString(salt)
	b.Headers["Nonce"] = base64.StdEncoding.EncodeToString(nonce)
	b.Bytes = aesgcm.Seal(nil, nonce, data, nil)

	fmt.Println(`package main

// encrypt = ignore

/*

This file is encrypted.
The key is the answer to the indicated problem.

`)

	err = pem.Encode(os.Stdout, &b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: encode output: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("*/")
}

func smudge() {
	data, err := ioutil.ReadAll(os.Stdin)
	errCheck(err, "read stdin")

	b, _ := pem.Decode(data)
	if b == nil || b.Type != pemType {
		os.Stdout.Write(data)
		return
	}

	salt, err := base64.StdEncoding.DecodeString(b.Headers["Salt"])
	errCheck(err, "decode salt")
	nonce, err := base64.StdEncoding.DecodeString(b.Headers["Nonce"])
	errCheck(err, "decode nonce")

	encKey := getKey(b.Headers["Problem Number"])

	key, err := scrypt.Key([]byte(encKey), salt, 16384, 8, 1, 32)
	errCheck(err, "calculate encryption key")

	c, err := aes.NewCipher(key)
	errCheck(err, "create new cipher")

	aesgcm, err := cipher.NewGCM(c)
	errCheck(err, "initialize cipher")

	data, err = aesgcm.Open(nil, nonce, b.Bytes, nil)
	errCheck(err, "decrypt data")

	os.Stdout.Write(data)
}
