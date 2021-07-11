package msec

import (
	"crypto/aes"
	"crypto/cipher"
	"log"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
)

const (
	blockSize        = 256 / 8
	pbkdf2Iterations = 4096
	salt             = "MEN3EM-VERY-RANDOM-AND-SECRET-SALT"
)

var iv [blockSize / 2]byte

type MSecLayer struct {
	block cipher.Block
}

func NewMSecLayer(pass string) *MSecLayer {
	key := pbkdf2.Key([]byte(pass), []byte(salt), pbkdf2Iterations, blockSize, sha3.New256)
	if len(key) != blockSize {
		log.Panicf("key.len = %d and blocksize = %d, mismatch", len(key), blockSize)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic("failed to create aes cipher: ", err)
	}

	return &MSecLayer{block: block}
}

func (msec *MSecLayer) Encrypt(in []byte) []byte {
	out := make([]byte, len(in))
	cipher.NewCFBEncrypter(msec.block, iv[:]).XORKeyStream(out, in)
	return out
}

func (msec *MSecLayer) Decrypt(in []byte) []byte {
	out := make([]byte, len(in))
	cipher.NewCFBDecrypter(msec.block, iv[:]).XORKeyStream(out, in)
	return out
}
