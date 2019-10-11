package client

import (
	"common"
	"net"
)

//************//
// CRC uint32
//************//
// request seq uint32
//************//
// response seq uint32
//************//
// AcccountID
//************//
// PacketType uint32
//************//
// Packet
//************//


type RUDPData struct{
	Packet *common.Packet
	Seq uint32

}

//UDPConnect 연결
func UDPConnect(ip string, port int16) *net.UDPConn {

	ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Fatal(err)
		return nil
	}

	LocalAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return conn
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


//UDPSendWorker 전송용
func  UDPSendWorker(sendData *common.Channel, recvData *common.Channel, conn *net.UDPConn, accountID uint32) {
	sendID := sendData.ID
	recvID := recvData.ID

	defer func() {
		fmt.Println("exit UDPSendWorker", sendID, recvID)
	}()

	seq := uint32(1)

	dataIndex := 0
	var waittingList [128]*RUDPData

	for {
		if dataIndex < 128 {
			packet := sendData.GetDefault()
			if 	packet != nil {
				data := &RUDPData{
					Seq: seq,
					Packet: packet,
				}
				waittingList[dataIndex] = data
				dataIndex ++
				seq ++
				continue
			}
		}

		if dataIndex > 0 {

		}


		packet := sendData.GetDefault()
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