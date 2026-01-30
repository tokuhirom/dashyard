// genhash generates a SHA-512 crypt password hash for use in Dashyard config files.
// Usage: go run ./cmd/genhash <password>
package main

import (
	"fmt"
	"os"

	"github.com/GehirnInc/crypt"
	_ "github.com/GehirnInc/crypt/sha512_crypt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <password>\n", os.Args[0])
		os.Exit(1)
	}

	c := crypt.SHA512.New()
	hash, err := c.Generate([]byte(os.Args[1]), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(hash)
}
