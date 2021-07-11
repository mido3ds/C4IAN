package flood

import (
	"log"
	"net"

	. "github.com/mido3ds/C4IAN/src/router/mac"
	. "github.com/mido3ds/C4IAN/src/router/msec"
)

type GlobalFlooder struct {
	seqNumber   uint32
	floodingTbl *floodingTable
	macConn     *MACLayerConn
	ip          net.IP
	msec        *MSecLayer
}

func NewGlobalFlooder(ip net.IP, iface *net.Interface, etherType EtherType, msec *MSecLayer) *GlobalFlooder {
	// connect to mac layer
	macConn, err := NewMACLayerConn(iface, etherType)
	if err != nil {
		log.Panic("failed to create device connection, err: ", err)
	}

	log.Println("initalized global flooder")

	return &GlobalFlooder{
		seqNumber:   0,
		floodingTbl: newFloodingTable(),
		macConn:     macConn,
		ip:          ip,
		msec:        msec,
	}
}

func (f *GlobalFlooder) Flood(encryptedPayload []byte) {
	f.seqNumber++

	floodHeader := FloodHeader{SrcIP: f.ip, SeqNum: f.seqNumber}
	encryptedFloodHeader := f.msec.Encrypt(floodHeader.MarshalBinary())

	f.macConn.Write(append(encryptedFloodHeader, encryptedPayload...), BroadcastMACAddr)
}

// ListenForFloodedMsgs inf loop that receives any flooded msgs
// calls `payloadProcessor` when it receives the message, it gives it the header and the payload
// and returns whether to continue flooding or not
func (f *GlobalFlooder) ListenForFloodedMsgs(payloadProcessor func([]byte) []byte) {
	for {
		msg := f.macConn.Read()
		go f.handleFloodedMsg(msg, payloadProcessor)
	}
}

func (f *GlobalFlooder) handleFloodedMsg(msg []byte, payloadProcessor func([]byte) []byte) {
	floodHeader, ok := UnmarshalFloodedHeader(f.msec.Decrypt(msg[:floodHeaderLen]))
	if !ok {
		return
	}

	if net.IP.Equal(floodHeader.SrcIP, f.ip) {
		return
	}

	if f.floodingTbl.isHighestSeqNum(floodHeader.SrcIP, floodHeader.SeqNum) {
		f.floodingTbl.set(floodHeader.SrcIP, floodHeader.SeqNum)
		encryptedPayload := msg[floodHeaderLen:]
		newEncryptedPayload := payloadProcessor(encryptedPayload)
		if newEncryptedPayload != nil {
			f.macConn.Write(append(msg[:floodHeaderLen], newEncryptedPayload...), BroadcastMACAddr)
		}
	}
}

func (f *GlobalFlooder) Close() {
	f.macConn.Close()
}
