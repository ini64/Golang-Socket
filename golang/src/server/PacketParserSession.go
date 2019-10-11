package server

import (
	"common"
	"fmt"
	fb "packet"
)

//PacketParser 작업 실행
func (d *SessionManagerData) PacketParser() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("run time panic", string(common.Stack()), err)
		}
	}()
	packet := d.Packet
	packetType := packet.PacketType
	switch packetType {
	case fb.SessionIBegin:
		d.Begin()
	case fb.SessionIEnd:
		d.End()
	}
}

//Begin 입장
func (d *SessionManagerData) Begin() {
	packet := d.Packet
	sessionBegin := fb.GetRootAsSessionBegin(packet.Bytes(), 0)
	accountID := sessionBegin.AccountID()
	endPoint := d.EndPoint

	TraceSession(packet.PacketType, accountID)

	if accountID != 0 {

		_, ok := d.Users[accountID]
		if !ok {
			user := GetUser()
			d.Users[accountID] = user
		}

		tcpWrite, ok := d.TCPWrite[accountID]
		if ok {
			tcpWrite.Release()
		}
		d.TCPWrite[accountID] = endPoint.TCPWritePool.Get(accountID)

		udpWrite, ok := d.UDPWrite[accountID]
		if ok {
			udpWrite.Release()
		}
		d.UDPWrite[accountID] = endPoint.UDPWritePool.Get(accountID)
	}

}

//End 퇴장
func (d *SessionManagerData) End() {
	packet := d.Packet
	sessionEnd := fb.GetRootAsSessionEnd(packet.Bytes(), 0)
	accountID := sessionEnd.AccountID()

	TraceSession(packet.PacketType, accountID)

	if accountID != 0 {

		user, ok := d.Users[accountID]
		if ok {
			PutUser(user)
			delete(d.Users, accountID)
		}

		tcpWrite, ok := d.TCPWrite[accountID]
		if ok {
			tcpWrite.Release()
			delete(d.TCPWrite, accountID)
		}

		udpWrite, ok := d.UDPWrite[accountID]
		if ok {
			udpWrite.Release()
			delete(d.UDPWrite, accountID)
		}
	}
}
