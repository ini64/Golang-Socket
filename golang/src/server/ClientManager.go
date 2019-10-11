package server

import (
	"common"
	"fmt"
	"net"
	fb "packet"
	"sync"

	flatbuffers "github.com/google/flatbuffers/go"
)

//ClientData 클라이언트 접근 정보
type ClientData struct {
	EndPoint *EndPoint
	Packet   *common.Packet
	Builder  *flatbuffers.Builder

	AccountID uint32
	TCPConn   *net.TCPConn

	TCPWrite *common.Channel
	TCPRead  *common.Channel
	UDPWrite *common.Channel
}

//Reset 초기화
func (d *ClientData) Reset() {
	d.EndPoint = nil
	if d.Packet != nil {
		common.PutPacket(d.Packet)
		d.Packet = nil
	}

	if d.Builder != nil {
		common.PutBuilder(d.Builder)
		d.Builder = nil
	}

	d.AccountID = 0
	d.TCPConn = nil

	if d.TCPWrite != nil {
		d.TCPWrite.Close()
		d.TCPWrite.Release()
		d.TCPWrite = nil
	}

	if d.TCPRead != nil {
		d.TCPRead.Close()
		d.TCPRead.Release()
		d.TCPRead = nil
	}

	if d.UDPWrite != nil {
		d.UDPWrite.Close()
		d.UDPWrite.Release()
		d.UDPWrite = nil
	}

}

// ClientDataPool sync.Pool
var ClientDataPool = sync.Pool{
	New: func() interface{} {
		return new(ClientData)
	},
}

//TCPWriteWorker 전송용
func (e *EndPoint) TCPWriteWorker(channel *common.Channel, conn *net.TCPConn, accountID uint32) {
	id := channel.ID
	ok := true
	var buffer [65536]byte

	defer func() {
		fmt.Println("exit TCPWriteWorker", accountID, id)
	}()

	for {
		packet := channel.Get()
		if packet == nil {
			channel.Release()
			conn.Close()
			return
		}
		if ok {
			ok = e.WriteTCP(conn, buffer[:], packet)
			if !ok {
				fmt.Println("close WriteWorker", accountID)
			}

		}
		common.PutPacket(packet)
	}
}

//TCPReadWorker 전송용
func (e *EndPoint) TCPReadWorker(channel *common.Channel, conn *net.TCPConn, accountID uint32) {
	var buffer [65536]byte
	id := channel.ID

	defer func() {
		fmt.Println("exit TCPReadWorker", accountID, id)
	}()

	for {
		packet := e.ReadTCP(conn, buffer[:])
		if packet == nil {
			channel.Close()
			channel.Release()
			return
		}
		channel.Put(packet)
	}
}

//UDPWriteWorker 전송용
func (e *EndPoint) UDPWriteWorker(channel *common.Channel, accountID uint32) {
	id := channel.ID
	ok := true
	var dstPoint *net.UDPAddr
	var dstString string
	var buffer [1500]byte

	defer func() {
		fmt.Println("exit UDPWriteWorker", accountID, id)
	}()

	for {
		packet := channel.Get()
		if packet == nil {
			channel.Release()
			return
		}

		if ok {
			if len(packet.Dst) > 0 {
				if dstPoint == nil || dstString != packet.Dst {
					dst, err := net.ResolveUDPAddr("udp", packet.Dst)
					if err != nil {
						fmt.Println("ResolveUDPAddr", err.Error())
						common.PutPacket(packet)
						continue
					}
					dstPoint = dst
					dstString = packet.Dst
				}
			}
			if dstPoint != nil {
				ok = e.WriteUDP(dstPoint, buffer[:], packet)
			}
		}
		common.PutPacket(packet)
	}
}

// ClientManager 개인 유저 처리
func (e *EndPoint) ClientManager(conn *net.TCPConn, slowlyConnection chan bool) {

	data := ClientDataPool.Get().(*ClientData)
	data.Builder = common.GetBuilder()
	data.EndPoint = e

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ClientManager", "run time panic", string(common.Stack()), err)
		}

		var activeUser bool
		tcpWrite := data.TCPWrite
		udpWrite := data.UDPWrite

		if tcpWrite != nil {
			tcpWrite.Close()
			activeUser = tcpWrite.Release()
			data.TCPWrite = nil
		}

		if udpWrite != nil {
			udpWrite.Close()
			udpWrite.Release()
			data.UDPWrite = nil
		}

		if activeUser {
			if data.AccountID != 0 {
				seesionEnd := e.SessionEnd(data.Builder, data.AccountID)
				e.SessionChannel.Put(seesionEnd)
			}
		}

		fmt.Println("exit ClientManager", data.AccountID)

		data.Reset()
		ClientDataPool.Put(data)
	}()

	var buffer [65535]byte
	packet := e.ReadTCP(conn, buffer[:])

	<-slowlyConnection

	if packet != nil {
		clientLogin := fb.GetRootAsClientBegin(packet.Bytes(), 0)
		data.AccountID = clientLogin.AccountID()

		common.PutPacket(packet)

		if data.AccountID != 0 {

			tcpWrite := e.TCPWritePool.MakeNGet(data.AccountID, e.Conf.ClientChannelSize, 2)
			tcpRead := e.TCPReadPool.MakeNGet(data.AccountID, e.Conf.ClientChannelSize, 2)
			udpWrite := e.UDPWritePool.MakeNGet(data.AccountID, e.Conf.ClientChannelSize, 2)

			go e.TCPReadWorker(tcpRead, conn, data.AccountID)
			go e.TCPWriteWorker(tcpWrite, conn, data.AccountID)
			go e.UDPWriteWorker(udpWrite, data.AccountID)

			data.TCPWrite = tcpWrite
			data.UDPWrite = udpWrite

			sessionBegin := e.SessionBegin(data.Builder, data.AccountID)
			e.SessionChannel.Put(sessionBegin)

			serverBegin := e.ServerBegin(data.Builder)
			TraceServer(serverBegin.PacketType, data.AccountID)
			data.TCPWrite.Put(serverBegin)

			for {
				packet := tcpRead.Get()
				if packet == nil {
					tcpRead.Release()
					return
				}
				data.Packet = packet
				data.PacketParser()
				common.PutPacket(packet)
			}
		}
	}
}
