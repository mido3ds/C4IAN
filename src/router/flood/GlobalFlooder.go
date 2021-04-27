package flood

import (
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/msec"
	. "github.com/mido3ds/C4IAN/src/router/tables"
)

type GlobalFlooder struct {
	seqNumber uint32
	fTable    *FloodingTable
	macConn   *MACLayerConn
	ip        net.IP
	msec      *MSecLayer
}

func NewGlobalFlooder(ip net.IP, iface *net.Interface, etherType EtherType, msec *MSecLayer) (*GlobalFlooder, error) {
	// connect to mac layer
	macConn, err := NewMACLayerConn(iface, etherType)
	if err != nil {
		return nil, err
	}

	fTable := NewFloodingTable()

	log.Println("initalized global flooder")

	return &GlobalFlooder{
		seqNumber: 0,
		fTable:    fTable,
		macConn:   macConn,
		ip:        ip,
		msec:      msec,
	}, nil
}

func (f *GlobalFlooder) Flood(msg []byte) {
	hdr := FloodHeader{SrcIP: f.ip, SeqNum: f.seqNumber}
	msg = append(hdr.MarshalBinary(), msg...)

	f.seqNumber++

	f.macConn.Write(f.msec.Encrypt(msg), BroadcastMACAddr)
}

// ReceiveFloodedMsgs inf loop that receives any flooded msgs
// calls `payloadProcessor` when it receives the message, it gives it the header and the payload
// and returns whether to continue flooding or not
func (f *GlobalFlooder) ReceiveFloodedMsgs(payloadProcessor func(*FloodHeader, []byte) bool) {
	for {
		msg := f.macConn.Read()

		pd := f.msec.NewPacketDecrypter(msg)
		decryptedHDR := pd.DecryptN(FloodHeaderLen)

		hdr, _, ok := UnmarshalFloodedHeader(decryptedHDR)
		if !ok {
			return
		}

		if net.IP.Equal(hdr.SrcIP, f.ip) {
			return
		}

		tableSeq, exist := f.fTable.Get(hdr.SrcIP)

		if exist && hdr.SeqNum <= tableSeq {
			return
		}

		f.fTable.Set(hdr.SrcIP, hdr.SeqNum)

		go func() {
			payload := pd.DecryptAll()[FloodHeaderLen:]

			if !payloadProcessor(hdr, payload) {
				return
			}

			f.macConn.Write(msg, BroadcastMACAddr)
		}()
	}
}

func (f *GlobalFlooder) Close() {
	f.macConn.Close()
}
