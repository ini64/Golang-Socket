package server

import (
	"common"
	"fmt"
)

//PacketParser 작업 실행
func (d *UDPClientData) PacketParser() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("UDPClientData", "PacketParser", "run time panic", string(common.Stack()), err)
		}
	}()

	packet := d.Packet

	switch packet.PacketType {
	}
}

//Login 로그인 패킷
func (d *UDPClientData) Login() {

	// packet := d.Packet
	// clientSyncUDP := fb.GetRootAsClientSyncUDP(packet.Buffer.Bytes(), 0)
	// sTime := clientSyncUDP.Stime()

	// dst := fmt.Sprintf("%s:%d", d.UDPAddr.IP.String(), d.UDPAddr.Port)

	// TraceUDPClient(packet.PacketType, d.AccountID, sTime, dst)

	// if sTime < 1 {
	// 	return
	// }

	// sessionSyncUDP := d.EndPoint.SessionSyncUDP(d.Builder, d.AccountID, clientSyncUDP, dst)
	// d.Channel.Put(sessionSyncUDP)
}
