package store

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"errors"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

type Result struct {
	URL       string
	Username  string
	Password  []byte
	TimesUsed int
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

	cbc := cipher.NewCBCDecrypter(block, []byte("                "))
	cbc.CryptBlocks(cleartext, encrypted)

	stripped, err := pkcs7strip(cleartext, 16)
	if err != nil {
		return "", fmt.Errorf("unable to strip payload: %w", err)
	}

	return string(stripped), nil
}

// pkcs7strip remove pkcs7 padding
// https://gist.github.com/nanmu42/b838acc10d393bc51cb861128ce7f89c
func pkcs7strip(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("pkcs7: Data is empty")
	}
	if length%blockSize != 0 {
		return nil, errors.New("pkcs7: Data is not block-aligned")
	}
	padLen := int(data[length-1])
	ref := bytes.Repeat([]byte{byte(padLen)}, padLen)
	if padLen > blockSize || padLen == 0 || !bytes.HasSuffix(data, ref) {
		return nil, errors.New("pkcs7: Invalid padding")
	}
	return data[:length-padLen], nil
}
