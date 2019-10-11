package server

import (
	"common"
	"fmt"
)

//PacketParser 작업 실행
func (d *ClientData) PacketParser() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ClientData", "PacketParser", "run time panic", string(common.Stack()), err)
		}
	}()

	packet := d.Packet

	switch packet.PacketType {
	// case fb.ClientLogin:
	// 	d.Login()
	}
}

// //Login 로긴
// func (d *ClientData) Login() {
// 	// packet := d.Packet()
// 	// login := fb.GetRootAsClientPing(packet.Bytes(), 0)
// }
