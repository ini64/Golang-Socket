package common

import (
	"bytes"
	"sync"
	"sync/atomic"
)

//Packet 통신패킷
type Packet struct {
	PacketType uint32
	Buffer     bytes.Buffer
	ID         uint64
	Dst        string
}

//Reset 리셋
func (p *Packet) Reset() {
	p.PacketType = 0
	p.ID = 0
	p.Buffer.Reset()
}

//Copy 복사
func (p *Packet) Copy() *Packet {
	packet := GetPacket()
	packet.PacketType = p.PacketType
	packet.Buffer.Write(p.Buffer.Bytes())
	return packet
}

//Bytes 포인터 리턴
func (p *Packet) Bytes() []byte {
	return p.Buffer.Bytes()
}

//PacketPoolCount 카운트 측정
var PacketPoolCount int32

// PacketMemoryPool 메모리 풀
var PacketMemoryPool = sync.Pool{
	New: func() interface{} {
		return new(Packet)
	},
}

//GetPacket 패킷 얻기
func GetPacket() *Packet {
	packet := PacketMemoryPool.Get().(*Packet)
	atomic.AddInt32(&PacketPoolCount, 1)
	return packet
}

//PutPacket 패킷 반납
func PutPacket(packet *Packet) {
	if packet != nil {
		atomic.AddInt32(&PacketPoolCount, -1)
		packet.Reset()
		PacketMemoryPool.Put(packet)
	}
}
