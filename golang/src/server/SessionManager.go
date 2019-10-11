package server

import (
	"common"
	"fmt"
	"sync"

	flatbuffers "github.com/google/flatbuffers/go"
)

//SessionManagerData 데이터 처리
type SessionManagerData struct {
	Packet   *common.Packet
	EndPoint *EndPoint
	Builder  *flatbuffers.Builder
	TCPWrite map[uint32]*common.Channel
	UDPWrite map[uint32]*common.Channel
	Users    map[uint32]*User
}

//SessionManager 유저의 상태 저장
func (e *EndPoint) SessionManager(end *sync.WaitGroup) {
	defer func() {
		fmt.Println("exit SessionManager")
		end.Done()
	}()

	builder := flatbuffers.NewBuilder(2048)
	tcpWrite := make(map[uint32]*common.Channel)
	udpWrite := make(map[uint32]*common.Channel)
	users := make(map[uint32]*User)

	d := &SessionManagerData{
		EndPoint: e,
		Builder:  builder,
		TCPWrite: tcpWrite,
		UDPWrite: udpWrite,
		Users:    users,
	}

	for {
		packet := e.SessionChannel.Get()
		if packet == nil {
			e.SessionChannel.Release()
			return
		}

		d.Packet = packet
		d.PacketParser()

		common.PutPacket(packet)
	}
}
