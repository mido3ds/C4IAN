package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
)

type MSecLayer struct {
	key []byte
}

// TODO: increase speed by lazy readers/writers
// TODO: add DecryptDest

func NewMSecLayer(pass string) *MSecLayer {
	return &MSecLayer{
		key: pbkdf2.Key([]byte(pass), []byte("MEN3EM-VERY-RANDOM-AND-SECRET-SALT"), 4096, 32, sha3.New256),
	}
}

func (msec *MSecLayer) Encrypt(in []byte) ([]byte, error) {
	bReader := bytes.NewBuffer(in)
	block, err := aes.NewCipher(msec.key)
	if err != nil {
		return nil, err
	}

	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	var out bytes.Buffer

	writer := &cipher.StreamWriter{S: stream, W: &out}
	if _, err := io.Copy(writer, bReader); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func (msec *MSecLayer) DecryptAll(in []byte) ([]byte, error) {
	bReader := bytes.NewBuffer(in)
	block, err := aes.NewCipher(msec.key)
	if err != nil {
		return nil, err
	}

	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	var out bytes.Buffer

	reader := &cipher.StreamReader{S: stream, R: bReader}
	if _, err := io.Copy(&out, reader); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
