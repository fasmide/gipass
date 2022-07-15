package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

type Result struct {
	URL      string
	Username string
	Password []byte
}

func (r Result) CleartextPassword() (string, error) {
	dk := pbkdf2.Key([]byte("peanuts"), []byte("saltysalt"), 1, 16, sha1.New)
	block, err := aes.NewCipher(dk)
	if err != nil {
		return "", fmt.Errorf("cannot create new aes cipher: %w", err)
	}

	encrypted := r.Password[3:]
	cleartext := make([]byte, len(encrypted))

	cbc := cipher.NewCBCDecrypter(block, []byte("                "))
	cbc.CryptBlocks(cleartext, encrypted)

	return string(cleartext), nil
}
