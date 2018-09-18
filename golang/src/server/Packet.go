package server

import (
	"bytes"
	"sync"
	"sync/atomic"

	flatbuffers "github.com/google/flatbuffers/go"
)

//Packet 통신패킷
type Packet struct {
	PacketType uint32
	Buffer     bytes.Buffer
}

//Reset 리셋
func (p *Packet) Reset() {
	p.PacketType = 0
	p.Buffer.Reset()
}

//Copy 복사
func (p *Packet) Copy() *Packet {
	packet := GetPacket(p.PacketType, nil)
	packet.Buffer.Write(p.Buffer.Bytes())
	return packet
}

//PacketPoolCount 카운트 측정
var PacketPoolCount int32

// PacketMemoryPool 메모리 풀
var PacketMemoryPool = sync.Pool{
	New: func() interface{} {
		return new(Packet)
	},
}

//GetPacket 얻기
func GetPacket(packetType uint32, builder *flatbuffers.Builder) *Packet {
	packet := PacketMemoryPool.Get().(*Packet)

	packet.PacketType = packetType
	if builder != nil {
		packet.Buffer.Write(builder.FinishedBytes())
	}

	atomic.AddInt32(&PacketPoolCount, 1)
	return packet
}

//PutPacket 반납
func PutPacket(packet *Packet) {
	if packet != nil {
		packet.Reset()
		atomic.AddInt32(&PacketPoolCount, -1)
	}
}
