package msec

import (
	"bytes"
	"crypto/cipher"
	"io"
	"log"
)

type PacketDecrypter struct {
	reader *cipher.StreamReader
	out    *bytes.Buffer
}

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

// DecryptN returns last N decrypted bytes of the packet
// advances the buffer index by N bytes, so next call will decrypt the next bytes
func (p *PacketDecrypter) DecryptN(n int) []byte {
	if n <= 0 {
		log.Panic("packet decrypter: n must be positive")
	}

	n2, err := io.CopyN(p.out, p.reader, int64(n))
	if int64(n) != n2 {
		log.Panic("packet decrypter failed to decrypt n bytes")
	}
	if err != nil {
		log.Panic("packet decrypter err:", err)
	}

	b := p.out.Bytes()
	return b[len(b)-n:]
}

// DecryptAll decrypts the rest of the packet and returns all decrypted bytes
func (p *PacketDecrypter) DecryptAll() []byte {
	_, err := io.Copy(p.out, p.reader)
	if err != nil {
		log.Panic("failed to decrypt all of the packet, err: ", err)
	}

	return p.out.Bytes()
}
