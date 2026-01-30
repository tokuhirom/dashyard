package auth

import (
	"testing"

	"github.com/GehirnInc/crypt"
	_ "github.com/GehirnInc/crypt/sha512_crypt"
)

func generateHash(password string) string {
	c := crypt.SHA512.New()
	hash, err := c.Generate([]byte(password), nil)
	if err != nil {
		panic(err)
	}
	return hash
}

func TestVerifyPassword(t *testing.T) {
	hash := generateHash("correctpassword")

	if !VerifyPassword("correctpassword", hash) {
		t.Error("expected password to verify successfully")
	}

	if VerifyPassword("wrongpassword", hash) {
		t.Error("expected wrong password to fail verification")
	}
}

func TestVerifyPasswordInvalidHash(t *testing.T) {
	if VerifyPassword("password", "not-a-valid-hash") {
		t.Error("expected invalid hash to fail verification")
	}
}
