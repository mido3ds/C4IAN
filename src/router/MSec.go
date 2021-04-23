package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"log"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
)

type MSecLayer struct {
	key []byte
}

func NewMSecLayer(pass string) *MSecLayer {
	return &MSecLayer{
		key: pbkdf2.Key([]byte(pass), []byte("MEN3EM-VERY-RANDOM-AND-SECRET-SALT"), 4096, 32, sha3.New256),
	}
}

func (msec *MSecLayer) Encrypt(in []byte) []byte {
	bReader := bytes.NewBuffer(in)
	block, err := aes.NewCipher(msec.key)
	if err != nil {
		log.Panic("failed to encrypt, err: ", err)
	}

	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	var out bytes.Buffer

	writer := &cipher.StreamWriter{S: stream, W: &out}
	if _, err := io.Copy(writer, bReader); err != nil {
		log.Panic("failed to encrypt, err: ", err)
	}

	return out.Bytes()
}

func (msec *MSecLayer) decryptStream() (cipher.Stream, error) {
	block, err := aes.NewCipher(msec.key)
	if err != nil {
		return nil, err
	}

	var iv [aes.BlockSize]byte
	return cipher.NewOFB(block, iv[:]), nil
}

func (msec *MSecLayer) Decrypt(in []byte) []byte {
	stream, err := msec.decryptStream()
	if err != nil {
		log.Panic("failed to decrypt msg, err: ", err)
	}

	var out bytes.Buffer

	reader := &cipher.StreamReader{S: stream, R: bytes.NewBuffer(in)}
	if _, err := io.Copy(&out, reader); err != nil {
		log.Panic("failed to decrypt msg, err: ", err)
	}

	return out.Bytes()
}

type PacketDecrypter struct {
	reader *cipher.StreamReader
	out    *bytes.Buffer
}

// TODO: add DecryptN, and remove DecryptAndVerify

func (msec *MSecLayer) NewPacketDecrypter(in []byte) *PacketDecrypter {
	stream, err := msec.decryptStream()
	if err != nil {
		log.Panic("failed to build packet decrypter, err: ", err)
	}

	out := new(bytes.Buffer)
	reader := &cipher.StreamReader{S: stream, R: bytes.NewBuffer(in)}

	return &PacketDecrypter{
		reader: reader,
		out:    out,
	}
}

func (p *PacketDecrypter) DecryptAndVerifyZID() (*ZIDHeader, bool) {
	n, err := io.CopyN(p.out, p.reader, ZIDHeaderLen)
	if err != nil || n != ZIDHeaderLen {
		return nil, false
	}
	b := p.out.Bytes()
	return UnmarshalZIDHeader(b[len(b)-ZIDHeaderLen:])
}

func (p *PacketDecrypter) DecryptAndVerifyIP() (*IPHeader, bool) {
	n, err := io.CopyN(p.out, p.reader, IPv4HeaderLen)
	if err != nil || n != IPv4HeaderLen {
		return nil, false
	}
	b := p.out.Bytes()
	return UnmarshalIPHeader(b[len(b)-IPv4HeaderLen:])
}

func (p *PacketDecrypter) DecryptAll() ([]byte, error) {
	_, err := io.Copy(p.out, p.reader)
	if err != nil {
		return nil, err
	}

	return p.out.Bytes(), nil
}
