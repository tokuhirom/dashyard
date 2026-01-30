package auth

import (
	"github.com/GehirnInc/crypt"
	_ "github.com/GehirnInc/crypt/sha512_crypt"
)

// VerifyPassword checks a plaintext password against a SHA-512 crypt hash ($6$ format).
func VerifyPassword(password, hash string) bool {
	c := crypt.SHA512.New()
	err := c.Verify(hash, []byte(password))
	return err == nil
}
