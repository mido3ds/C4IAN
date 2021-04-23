package odmrp

type BufferSlot struct {
	multicast_group int
	pid             int
	packet          *Packet
}

type PacketBuffer struct {
	packets_  []BufferSlot
	curr_pos_ int // current position
}

func newPacketBuffer() PacketBuffer {
	var pb PacketBuffer
	pb.packets_ = make([]BufferSlot, PACKET_BUFFER_SIZE)
	for i := 0; i < len(pb.packets_); i++ {
		pb.packets_[i].packet = nil
	}
	pb.curr_pos_ = -1
	return pb
}

func (pb *PacketBuffer) AddPacket(p *Packet, addr int, pid int) {
	//TODO
}

func (pb *PacketBuffer) GetFirstPacket(addr int) *Packet {
	//TODO
	return nil
}

func (pb *PacketBuffer) SendPackets(addr int, agent *Agent) int {
	num_packets_sent := 0
	// TODO
	return num_packets_sent
}
