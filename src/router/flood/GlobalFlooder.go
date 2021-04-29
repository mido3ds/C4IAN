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
	f.seqNumber++

	floodHeader := FloodHeader{SrcIP: f.ip, SeqNum: f.seqNumber}
	encryptedFloodHeader := f.msec.Encrypt(floodHeader.MarshalBinary())

	encryptedPayload := f.msec.Encrypt(msg)

	f.macConn.Write(append(encryptedFloodHeader, encryptedPayload...), BroadcastMACAddr)
}

// ListenForFloodedMsgs inf loop that receives any flooded msgs
// calls `payloadProcessor` when it receives the message, it gives it the header and the payload
// and returns whether to continue flooding or not
func (f *GlobalFlooder) ListenForFloodedMsgs(payloadProcessor func(*FloodHeader, []byte) ([]byte, bool)) {
	for {
		msg := f.macConn.Read()
		go f.handleFloodedMsg(msg, payloadProcessor)
	}
}

func (f *GlobalFlooder) handleFloodedMsg(msg []byte, payloadProcessor func(*FloodHeader, []byte) ([]byte, bool)) {
	floodHeader, ok := UnmarshalFloodedHeader(f.msec.Decrypt(msg[:FloodHeaderLen]))
	if !ok {
		return
	}

	if net.IP.Equal(floodHeader.SrcIP, f.ip) {
		return
	}

	tableSeq, exist := f.fTable.Get(floodHeader.SrcIP)

	if exist && floodHeader.SeqNum <= tableSeq {
		return
	}

	f.fTable.Set(floodHeader.SrcIP, floodHeader.SeqNum)

	payload := f.msec.Decrypt(msg[FloodHeaderLen:])
	payload, ok = payloadProcessor(floodHeader, payload)
	if !ok {
		return
	}

	f.macConn.Write(append(msg[:FloodHeaderLen], f.msec.Encrypt(payload)...), BroadcastMACAddr)
}

func (f *GlobalFlooder) Close() {
	f.macConn.Close()
}
