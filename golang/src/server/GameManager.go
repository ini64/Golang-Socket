package server

import (
	"fmt"
	"sync"

	flatbuffers "github.com/google/flatbuffers/go"

	sp "packet/Server"
)

//GameManagerData 데이터 처리
type GameManagerData struct {
	Data     *ChannelData
	EndPoint *EndPoint
	Builder  *flatbuffers.Builder
	Write    map[int32]Channel
	Users    map[int32]*User
}

//GameManager 전송 내역 관리
func (e *EndPoint) GameManager(wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("exit game manager")
		wg.Done()
	}()

	builder := flatbuffers.NewBuilder(2048)
	write := make(map[int32]Channel)
	users := make(map[int32]*User)

	//now := time.Now()

	var counter int
	for {
		select {
		case data := <-e.GameManagerChannel:

			d := &GameManagerData{
				Data:     data,
				EndPoint: e,
				Builder:  builder,
				Write:    write,
				Users:    users,
			}
			d.PacketParser()

			PutChannelData(data)
			counter++
		}
	}
}

//PacketParser 작업 실행
func (d *GameManagerData) PacketParser() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("run time panic", err)
		}
	}()
	packetType := d.Data.Packet.PacketType
	switch packetType {
	case sp.SSPacketTypeEnter:
		d.Enter()
	case sp.SSPacketTypeLeave:
		d.Leave()
	case sp.SSPacketTypeMove:
		d.Move()
	}
}

//Enter 입장
func (d *GameManagerData) Enter() {
	ssEnter := sp.GetRootAsSSEnter(d.Data.Packet.Buffer.Bytes(), 0)
	accountID := ssEnter.AccountID()

	_, ok := d.Write[accountID]
	if ok {
		ssEnterAck := d.SSEnterNak(d.Builder, 1)
		d.Data.Channel.Add(ssEnterAck)
	} else {
		d.Write[accountID] = d.Data.Channel

		user := GetUser()
		user.AccountID = ssEnter.AccountID()

		vec3 := ssEnter.Pos(nil)
		user.Pos.x = vec3.X()
		user.Pos.y = vec3.Y()
		user.Pos.z = vec3.Z()

		Qua := ssEnter.Rot(nil)
		user.Rot.x = Qua.X()
		user.Rot.y = Qua.Y()
		user.Rot.z = Qua.Z()
		user.Rot.w = Qua.W()

		user.Time = ssEnter.Time()

		d.Users[accountID] = user

		ssEnterOther := d.SSEnterOther(d.Builder, user)
		for key, channel := range d.Write {
			if key != accountID {
				other := ssEnterOther.Copy()
				channel.Add(other)
			}
		}
		PutPacket(ssEnterOther)

		ssEnterAck := d.SSEnterAck(d.Builder)
		d.Data.Channel.Add(ssEnterAck)
	}
}

//Leave 퇴장
func (d *GameManagerData) Leave() {
	ssEnter := sp.GetRootAsSSLeave(d.Data.Packet.Buffer.Bytes(), 0)
	accountID := ssEnter.AccountID()

	delete(d.Write, accountID)
	delete(d.Users, accountID)

	ssLeaveOther := d.SSLeaveOther(d.Builder, accountID)
	for key, channel := range d.Write {
		if key != accountID {
			other := ssLeaveOther.Copy()
			channel.Add(other)
		}
	}
	PutPacket(ssLeaveOther)
}

//Move 입장
func (d *GameManagerData) Move() {
	ssMove := sp.GetRootAsSSMove(d.Data.Packet.Buffer.Bytes(), 0)

	user := &sp.User{}
	ssMove.User(user)

	vec3 := &sp.Vec3{}
	qua := &sp.Qua{}

	user.Pos(vec3)
	user.Rot(qua)

	accountID := user.Num()

	memory, ok := d.Users[accountID]
	if ok {
		memory.Pos.x = vec3.X()
		memory.Pos.y = vec3.Y()
		memory.Pos.z = vec3.Z()

		fmt.Println(memory.Pos.x, memory.Pos.y, memory.Pos.z)

		memory.Rot.x = qua.X()
		memory.Rot.y = qua.Y()
		memory.Rot.z = qua.Z()
		memory.Rot.w = qua.W()

		memory.Time = user.Time()

		ssMoveAck := d.SSMoveAck(d.Builder, memory)
		for _, channel := range d.Write {
			//if key != accountID {
			other := ssMoveAck.Copy()
			channel.Add(other)
			//}
		}
		PutPacket(ssMoveAck)
	}
}
