package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func main() {
	toHash := []byte("thisismystringthatneedshashing")
	h := hmac.New(sha256.New, []byte("mysecretkey"))
	h.Write(toHash)
	b := h.Sum(nil)
	fmt.Println(b)
}
