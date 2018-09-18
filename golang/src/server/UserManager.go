package server

import (
	"fmt"
	"net"

	cp "packet/Client"
	sp "packet/Server"

	flatbuffers "github.com/google/flatbuffers/go"
)

//WriteWorker 전송용
func (e *EndPoint) WriteWorker(conn *net.TCPConn, write Channel) {
	ok := true
	for {
		select {
		case data := <-write:

			if data == nil {
				break
			}
			//채널 마지막까지 쓰기
			if ok {
				ok = WriteTCP(conn, data.Packet)
				data.Packet = nil
			}

			PutChannelData(data)
		}
	}
}

//ReadWorker 전송용
func (e *EndPoint) ReadWorker(conn *net.TCPConn, read Channel) {
	builder := GetBuilder()
	defer PutBuilder(builder)

	for {
		packet, ok := e.ReadTCP(conn)
		if !ok {
			if packet != nil {
				PutPacket(packet)
			}
			read.Add(e.SSClose(builder))
			break
		} else {
			read.Add(packet)
		}
	}
}

//UserManagerData 데이터 처리
type UserManagerData struct {
	Packet    *Packet
	EndPoint  *EndPoint
	Builder   *flatbuffers.Builder
	Write     Channel
	AccountID int32
}

// UserManager 유저별 데이터 읽기
func (e *EndPoint) UserManager(connectChannel chan bool, conn *net.TCPConn) {
	var accountID int32
	read := make(Channel, 8)
	write := make(Channel, 8)

	builder := GetBuilder()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Run Time Panic", err)
		}

		fmt.Println("exit: ", accountID)

		if accountID != 0 {

			//게임 메니져 종료처리
			packet := e.SSLeave(builder, accountID)
			e.GameManagerChannel.Add(packet)
			PutBuilder(builder)
		}
	}()

	packet, ok := e.ReadTCP(conn)
	if ok {
		if packet != nil {
			if packet.PacketType == cp.CSPacketTypeLogin {
				csLogin := cp.GetRootAsCSLogin(packet.Buffer.Bytes(), 0)

				user := csLogin.User(nil)
				accountID = user.Num()

				ssEnter := e.SSEnter(builder, csLogin)
				e.GameManagerChannel.AddChannel(ssEnter, read)

				go e.WriteWorker(conn, write)
				go e.ReadWorker(conn, read)

				fmt.Println(fmt.Sprintf("%d login", accountID))
			}
			PutPacket(packet)
		}
	}

	<-connectChannel

	for {
		select {
		case data := <-read:

			if data.Packet.PacketType == sp.SSPacketTypeClose {
				//쓰기 채널 종료
				PutChannelData(data)
				write.Close()
				return
			}

			d := &UserManagerData{
				EndPoint:  e,
				Packet:    data.Packet,
				Builder:   builder,
				AccountID: accountID,
				Write:     write,
			}
			d.PacketParser()
		}
	}
}

//PacketParser 작업 실행
func (d *UserManagerData) PacketParser() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("run time panic", err)
		}
	}()
	packetType := d.Packet.PacketType
	switch packetType {

	case sp.SSPacketTypeEnterAck:
		d.SSLoginAck()
	case sp.SSPacketTypeEnterOther:
		d.SSEnterOther()
	case sp.SSPacketTypeLeaveOther:
		d.SSLeaveOther()
	case cp.CSPacketTypeMove:
		d.CSMove()
	case sp.SSPacketTypeMoveAck:
		d.SSMoveAck()
	case cp.CSPacketTypeShoot:
		d.SSMoveAck()
	case sp.SSPacketTypeShootAck:
		d.SSMoveAck()
	}
}

//SSLoginAck 로그인 처리.
func (d *UserManagerData) SSLoginAck() {
	ssLoginAck := sp.GetRootAsSSEnterAck(d.Packet.Buffer.Bytes(), 0)
	d.Write.Add(d.SCLogin(d.Builder, ssLoginAck))
}

//SSEnterOther 다른 유저 입장
func (d *UserManagerData) SSEnterOther() {
	ssEnterOther := sp.GetRootAsSSEnterOther(d.Packet.Buffer.Bytes(), 0)
	scEnterOther := d.SCEnterOther(d.Builder, ssEnterOther)
	d.Write.Add(scEnterOther)
}

//SSLeaveOther 다른 유저 퇴장
func (d *UserManagerData) SSLeaveOther() {
	ssLeaveOther := sp.GetRootAsSSLeaveOther(d.Packet.Buffer.Bytes(), 0)
	scLeaveOther := d.SCLeaveOther(d.Builder, ssLeaveOther.AccountID())
	d.Write.Add(scLeaveOther)
}

//CSMove 이동처리
func (d *UserManagerData) CSMove() {
	csMove := cp.GetRootAsCSMove(d.Packet.Buffer.Bytes(), 0)
	ssMove := d.SSMove(d.Builder, csMove)
	d.EndPoint.GameManagerChannel.Add(ssMove)
}

//SSMoveAck 이동처리
func (d *UserManagerData) SSMoveAck() {
	ssMoveAck := sp.GetRootAsSSMoveAck(d.Packet.Buffer.Bytes(), 0)
	scMove := d.SCMove(d.Builder, ssMoveAck)
	d.Write.Add(scMove)
}
