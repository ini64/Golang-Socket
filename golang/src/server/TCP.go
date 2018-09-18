package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

//TCPListener 메세지 대기
func (e *EndPoint) TCPListener(wg *sync.WaitGroup) {
	//wg.Add(1)
	defer func() {
		wg.Done()
	}()

	var err error
	fmt.Println(e.TCPBind)
	listener, err := net.Listen("tcp", e.TCPBind)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tcpConn, ok := listener.(*net.TCPListener)
	if !ok {
		fmt.Println("change tcp connect not ok")
		return
	}

	e.TCPConn = tcpConn

	connect := make(chan bool, 256)

	for {
		e.TCPListenWorker(connect)
	}
}

//TCPListenWorker 개별 작업자
func (e *EndPoint) TCPListenWorker(connect chan bool) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("RunTime Panic", err)
		}
	}()

	conn, err := e.TCPConn.AcceptTCP()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = conn.SetNoDelay(true)
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return
	}

	err = conn.SetKeepAlive(true)
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return
	}

	err = conn.SetKeepAlivePeriod(time.Second * 4)
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return
	}

	connect <- true
	go e.UserManager(connect, conn)
}

//GetRecvTimeOut 타임 아웃 시간 얻기
func GetRecvTimeOut() time.Time {
	return time.Now().Add(time.Duration(360) * time.Second)
}

//GetSendTimeOut 타임 아웃 시간 얻기
func GetSendTimeOut() time.Time {
	return time.Now().Add(time.Duration(360) * time.Second)
}

//ReadTCP IO작업
func (e *EndPoint) ReadTCP(conn *net.TCPConn) (*Packet, bool) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("RunTime Panic", err)
		}
	}()

	var buffer [8]byte

	conn.SetReadDeadline(GetRecvTimeOut())
	headerSize, err := io.ReadFull(conn, buffer[:])

	if err != nil || headerSize != 8 {
		fmt.Println(err.Error())
		return nil, false
	}

	size := int64(binary.LittleEndian.Uint32(buffer[0:]))
	packetType := binary.LittleEndian.Uint32(buffer[4:])

	packet := GetPacket(packetType, nil)
	conn.SetReadDeadline(GetRecvTimeOut())
	bodySize, err := io.CopyN(&packet.Buffer, conn, size)

	if err != nil || bodySize != size {
		fmt.Println(err.Error())
		return packet, false
	}

	return packet, true
}

//WriteTCP 버퍼에 내용 쓰기
func WriteTCP(conn *net.TCPConn, packet *Packet) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("run tim panic", err)
		}
		if packet != nil {
			PutPacket(packet)
		}
	}()

	size := packet.Buffer.Len()

	var temp [8]byte

	binary.LittleEndian.PutUint32(temp[0:], uint32(size))
	binary.LittleEndian.PutUint32(temp[4:], uint32(packet.PacketType))

	buffer := GetBuffer()
	defer PutBuffer(buffer)
	buffer.Write(temp[:])
	buffer.Write(packet.Buffer.Bytes())

	sendSize := int64(size) + 8

	bodySize, err := io.CopyN(conn, buffer, sendSize)
	if bodySize != sendSize || err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
