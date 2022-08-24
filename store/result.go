package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"fmt"
	"log"

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

	// remove "v10"
	encrypted := r.Password[3:]
	cleartext := make([]byte, len(encrypted))
	log.Printf("pre-what %X len %d", cleartext, len(cleartext))
	cbc := cipher.NewCBCDecrypter(block, []byte("                "))
	cbc.CryptBlocks(cleartext, encrypted)
	log.Printf("what %X, \"%s\"", cleartext, string(cleartext))
	return string(cleartext), nil
}
