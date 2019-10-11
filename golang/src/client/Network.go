package client

import (
	"common"
	"encoding/binary"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

//TCPReadCount 수신 갯수
var TCPReadCount int64

//TCPReadSize 바이트 사이즈
var TCPReadSize int64

//UDPReadCount 수신 갯수
var UDPReadCount int64

//UDPReadSize 바이트 사이즈
var UDPReadSize int64

//TCPConnect 연결
func TCPConnect(conf *Conf) (*net.TCPConn, func() *common.Packet, func(packet *common.Packet) bool) {

	iconn, err := net.Dial("tcp", conf.TCPConnect)
	if err != nil {
		fmt.Println(err)
		return nil, nil, nil
	}
	conn, ok := iconn.(*net.TCPConn)
	if !ok {
		fmt.Println("변환 실패")
		return nil, nil, nil
	}
	conn.SetNoDelay(true)

	var readBuffer [65536]byte
	tcpRead := func() *common.Packet {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("run time panic", err, string(common.Stack()))
			}
		}()

		if conn == nil {
			return nil
		}

		var seek int

		for seek < 8 {
			size, err := conn.Read(readBuffer[seek:8])
			if err != nil {
				return nil
			}
			seek += size
		}

		if seek != 8 {
			return nil
		}

		size := int64(binary.LittleEndian.Uint32(readBuffer[0:4]))
		packetType := binary.LittleEndian.Uint32(readBuffer[4:8])

		for i := 0; i < 8; i++ {
			readBuffer[i] = 0
		}

		seek = 0

		for seek < int(size) {
			rSize, err := conn.Read(readBuffer[seek:size])
			if err != nil {
				return nil
			}
			seek += rSize
		}

		if seek != int(size) {
			return nil
		}

		atomic.AddInt64(&TCPReadCount, 1)
		atomic.AddInt64(&TCPReadSize, int64(size))

		packet := common.GetPacket()
		packet.PacketType = packetType
		packet.Buffer.Write(readBuffer[0:size])

		return packet
	}

	var writeBuffer [65536]byte
	tcpWrite := func(packet *common.Packet) bool {
		defer common.PutPacket(packet)

		packetBuffer := packet.Bytes()
		size := len(packetBuffer)
		sendSize := size + 8

		binary.LittleEndian.PutUint32(writeBuffer[0:4], uint32(size))
		binary.LittleEndian.PutUint32(writeBuffer[4:8], uint32(packet.PacketType))

		copy(writeBuffer[8:], packetBuffer)

		seek := int(0)
		for seek < sendSize {
			size, err := conn.Write(writeBuffer[0:8])
			if err != nil {
				return false
			}
			seek += size

			size, err = conn.Write(writeBuffer[seek:sendSize])
			if err != nil {
				return false
			}
			seek += size
		}

		if seek != sendSize {
			return false
		}

		return true
	}

	return conn, tcpRead, tcpWrite
}

// //TCPReader tcp에서 읽기
// func TCPRead(conn *net.TCPConn, buffer []byte) *common.Packet {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			fmt.Println("run time panic", err, string(common.Stack()))
// 		}
// 	}()

// 	if conn == nil {
// 		return nil
// 	}

// 	var seek int

// 	for seek < 8 {
// 		size, err := conn.Read(buffer[seek:8])
// 		if err != nil {
// 			//fmt.Println("ERROR", err.Error())
// 			return nil
// 		}
// 		seek += size
// 	}

// 	if seek != 8 {
// 		//fmt.Println("ERROR", "seek size", seek)
// 		return nil
// 	}

// 	size := int64(binary.LittleEndian.Uint32(buffer[0:4]))
// 	packetType := binary.LittleEndian.Uint32(buffer[4:8])

// 	for i := 0; i < 8; i++ {
// 		buffer[i] = 0
// 	}

// 	seek = 0

// 	for seek < int(size) {
// 		rSize, err := conn.Read(buffer[seek:size])
// 		if err != nil {
// 			return nil
// 		}
// 		seek += rSize
// 	}

// 	if seek != int(size) {
// 		return nil
// 	}

// 	atomic.AddInt64(&TCPReadCount, 1)
// 	atomic.AddInt64(&TCPReadSize, int64(size))

// 	packet := common.GetPacket()
// 	packet.PacketType = packetType
// 	packet.Buffer.Write(buffer)

// 	return packet
// }

// //TCPWrite 네트워크에 데이터 적기
// func TCPWrite(conn *net.TCPConn, packet *common.Packet) bool {
// 	defer common.PutPacket(packet)

// 	sendBuffer := packet.Bytes()
// 	len := len(sendBuffer)
// 	seek := 0

// 	for seek < len {
// 		size, err := conn.Write(sendBuffer[0:8])
// 		if err != nil {
// 			return false
// 		}
// 		seek += size

// 		size, err = conn.Write(sendBuffer[seek:len])
// 		if err != nil {
// 			return false
// 		}
// 		seek += size
// 	}

// 	if seek != len {
// 		return false
// 	}

// 	return true
// }

//TCPWait 네트워크 전송 대기
func TCPWait(packetType uint32, timeOut int, tcpRead func() *common.Packet) *common.Packet {

	channel := make(chan *common.Packet, 2)

	for {
		channel <- tcpRead()

		select {
		case <-time.After(time.Duration(timeOut) * time.Second):
			return nil
		case packet := <-channel:
			if packet.PacketType == packetType {
				return packet
			}
			common.PutPacket(packet)
		}
	}

}
