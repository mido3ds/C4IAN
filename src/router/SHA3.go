package main

import (
	"bytes"
	"log"

	"golang.org/x/crypto/sha3"
)

func Hash_SHA3(b []byte) []byte {
	h := sha3.New512()

	n, err := h.Write(b)
	if err != nil {
		log.Fatal("failed to hash, err: ", err)
	} else if n != len(b) {
		log.Fatal("failed to hash")
	}

	return h.Sum(nil)
}

func verifyHash_SHA3(data, h []byte) bool {
	return bytes.Equal(Hash_SHA3(data), h)
}
