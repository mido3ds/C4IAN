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

func (msec *MSecLayer) decryptStream() (cipher.Stream, error) {
	block, err := aes.NewCipher(msec.key)
	if err != nil {
		return nil, err
	}

	var iv [aes.BlockSize]byte
	return cipher.NewOFB(block, iv[:]), nil
}

func (msec *MSecLayer) decrypt(in []byte) ([]byte, error) {
	stream, err := msec.decryptStream()
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer

	reader := &cipher.StreamReader{S: stream, R: bytes.NewBuffer(in)}
	if _, err := io.Copy(&out, reader); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

type PacketDecrypter struct {
	reader *cipher.StreamReader
	out    *bytes.Buffer
}

// TODO: add DecryptN, and remove DecryptAndVerify

func (msec *MSecLayer) NewPacketDecrypter(in []byte) (*PacketDecrypter, error) {
	stream, err := msec.decryptStream()
	if err != nil {
		return nil, err
	}

	out := new(bytes.Buffer)
	reader := &cipher.StreamReader{S: stream, R: bytes.NewBuffer(in)}

	return &PacketDecrypter{
		reader: reader,
		out:    out,
	}, nil
}

func (p *PacketDecrypter) DecryptAndVerifyZID() (*ZIDHeader, bool) {
	n, err := io.CopyN(p.out, p.reader, ZIDHeaderLen)
	if err != nil || n != ZIDHeaderLen {
		return nil, false
	}
	return UnmarshalZIDHeader(p.out.Bytes())
}

func (p *PacketDecrypter) DecryptAndVerifyIP() (*IPHeader, bool) {
	n, err := io.CopyN(p.out, p.reader, IPv4HeaderLen)
	if err != nil || n != IPv4HeaderLen {
		return nil, false
	}

	return UnmarshalIPHeader(p.out.Bytes()[ZIDHeaderLen:])
}

func (p *PacketDecrypter) DecryptAll() ([]byte, error) {
	_, err := io.Copy(p.out, p.reader)
	if err != nil {
		return nil, err
	}

	return p.out.Bytes(), nil
}
