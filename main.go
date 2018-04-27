package main

import (
	"fmt"
	"github.com/kemics/rsa/pkg/rsa"
)

func main() {
	r := rsa.New()
	encoded := r.Encode("hello, world")
	decoded := r.Decode(encoded...)
	fmt.Println(decoded)
}
