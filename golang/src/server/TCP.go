package server

import (
	"common"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

//TCPListener 메세지 대기
func (e *EndPoint) TCPListener(end *sync.WaitGroup) {
	defer func() {
		end.Done()
	}()

	var err error
	fmt.Println(e.Conf.TCPBind)
	listener, err := net.Listen("tcp", e.Conf.TCPBind)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tcpConn, ok := listener.(*net.TCPListener)
	if !ok {
		fmt.Println("connect change failed")
		return
	}

	e.TCPConn = tcpConn

	slowlyConnection := make(chan bool, 256)

	for {
		if !e.TCPListenWorker(slowlyConnection) {
			return
		}
	}
}

//TCPListenWorker 개별 작업자
func (e *EndPoint) TCPListenWorker(slowlyConnection chan bool) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("TCPListenWorker", "run time panic", string(common.Stack()), err)
		}
	}()

	conn, err := e.TCPConn.AcceptTCP()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	err = conn.SetNoDelay(true)
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return false
	}

	err = conn.SetKeepAlive(true)
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return false
	}

	err = conn.SetKeepAlivePeriod(time.Second * 4)
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return false
	}

	slowlyConnection <- true
	go e.ClientManager(conn, slowlyConnection)

	return true
}

//ReadTCP IO작업
func (e *EndPoint) ReadTCP(conn *net.TCPConn, buffer []byte) *common.Packet {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ReadTCP", "run time panic", string(common.Stack()), err)
		}
	}()

	var seek int
	conn.SetReadDeadline(e.GetRecvTimeOut())

	for seek < 8 {
		size, err := conn.Read(buffer[seek:8])
		if err != nil {
			return nil
		}
		seek += size
	}

	if seek != 8 {
		return nil
	}

	size := int64(binary.LittleEndian.Uint32(buffer[0:4]))

	if int(size) >= len(buffer)-8 {
		return nil
	}

	packetType := binary.LittleEndian.Uint32(buffer[4:8])

	for i := 0; i < 8; i++ {
		buffer[i] = 0
	}

	seek = 0

	for seek < int(size) {
		conn.SetReadDeadline(e.GetRecvTimeOut())
		rSize, err := conn.Read(buffer[seek:size])
		if err != nil {
			return nil
		}
		seek += rSize
	}

	if seek != int(size) {
		return nil
	}

	packet := common.GetPacket()
	packet.PacketType = packetType
	packet.Buffer.Write(buffer)

	atomic.AddInt64(&e.TCPRecv, 1)

	return packet
}

//WriteTCP 버퍼에 내용 쓰기
func (e *EndPoint) WriteTCP(conn *net.TCPConn, buffer []byte, packet *common.Packet) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("WriteTCP", "run time panic", string(common.Stack()), err)
		}
	}()

	size := packet.Buffer.Len()
	sendSize := size + 8

	binary.LittleEndian.PutUint32(buffer[0:4], uint32(size))
	binary.LittleEndian.PutUint32(buffer[4:8], uint32(packet.PacketType))

	copy(buffer[8:], packet.Buffer.Bytes())

	seek := int(0)
	for seek < sendSize {

		writeSize, err := conn.Write(buffer[0:8])
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		seek += writeSize

		writeSize, err = conn.Write(buffer[seek:sendSize])
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		seek += writeSize
	}

	if seek != sendSize {
		fmt.Println("invalid send size", seek, sendSize)
		return false
	}

	atomic.AddInt64(&e.TCPSend, 1)

	return true
}
