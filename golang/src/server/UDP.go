package server

import (
	"common"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"net"
	"runtime"
	"sync"
	"sync/atomic"

	flatbuffers "github.com/google/flatbuffers/go"
)

//UDPListener 메세지 대기
func (e *EndPoint) UDPListener(end *sync.WaitGroup) {
	defer func() {
		end.Done()
	}()

	var waitGroup sync.WaitGroup

	var err error
	fmt.Println(e.Conf.UDPBind)
	ServerAddr, err := net.ResolveUDPAddr("udp", e.Conf.UDPBind)
	if err != nil {
		fmt.Println("ResolveUDPAddr failed", err.Error())
		return
	}

	conn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		fmt.Println("ListenUDP failed", err.Error())
		return
	}

	e.UDPConn = conn

	for i := 0; i < runtime.NumCPU(); i++ {
		waitGroup.Add(1)
		go e.UDPListenWorker(&waitGroup)
	}
	waitGroup.Wait()
}

//UDPClientData 데이터 처리용 객체
type UDPClientData struct {
	EndPoint  *EndPoint
	Packet    *common.Packet
	AccountID uint32
	UDPAddr   *net.UDPAddr
	Builder   *flatbuffers.Builder
	Channel   *common.Channel
}

//************//
// RoomID(proxy)
//************//
// CRC
//************//
// AcccountID
//************//
// PacketType
//************//
// Packet
//************//

//UDPListenWorker udp 수신작업
func (e *EndPoint) UDPListenWorker(wg *sync.WaitGroup) {

	data := &UDPClientData{
		Builder:  common.GetBuilder(),
		EndPoint: e,
		Channel:  e.SessionChannel,
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("UdpListenWorker", "run time panic", string(common.Stack()), err)
		}
		common.PutBuilder(data.Builder)

		wg.Done()
	}()

	//Ethernet v2 1500[8] 이더넷 구현에서 거의 대부분의 IP는 이더넷 V2 프레임 형식을 사용한다.
	buf := make([]byte, 1500)

	for {
		n, addr, err := e.UDPConn.ReadFromUDP(buf)

		if err != nil {
			fmt.Println("ReadFromUDP", err.Error())
			return
		}

		if n < 13 {
			continue
		}

		pos := 0

		checkSum := binary.LittleEndian.Uint32(buf[pos:])
		pos += 4

		//udp 데이터 검증
		newChecksum := crc32.ChecksumIEEE(buf[pos:n])
		if checkSum != newChecksum {
			fmt.Println("invalid check sum", checkSum, newChecksum)
			continue
		}

		accoundID := binary.LittleEndian.Uint32(buf[pos:])
		pos += 4

		packet := common.GetPacket()

		packet.PacketType = binary.LittleEndian.Uint32(buf[pos:])
		pos += 4

		packet.Buffer.Write(buf[pos:n])

		data.AccountID = accoundID
		data.Packet = packet
		data.UDPAddr = addr

		data.PacketParser()

		atomic.AddInt64(&e.UDPRecv, 1)

		common.PutPacket(packet)
	}
}

//WriteUDP udp전송 함수
func (e *EndPoint) WriteUDP(addr *net.UDPAddr, buffer []byte, packet *common.Packet) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("WriteUDP", "run time panic", string(common.Stack()), err)
		}
	}()

	//var temp [1500]byte

	pos := 4
	binary.LittleEndian.PutUint32(buffer[pos:], packet.PacketType)
	pos += 4

	copy(buffer[pos:], packet.Buffer.Bytes())
	pos += packet.Buffer.Len()

	checksum := crc32.ChecksumIEEE(buffer[4:pos])
	binary.LittleEndian.PutUint32(buffer[0:], checksum)

	seek := 0
	for seek < pos {
		sendSize, err := e.UDPConn.WriteToUDP(buffer[seek:pos], addr)

		if err != nil {
			//fmt.Println("WriteToUDP failed", err.Error())
			return false
		}

		seek += sendSize
	}
	atomic.AddInt64(&e.UDPSend, 1)

	return true
}
