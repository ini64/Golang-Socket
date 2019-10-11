package client

import (
	"common"

	flatbuffers "github.com/google/flatbuffers/go"
)

func getPacket(builder *flatbuffers.Builder, packetType uint32) *common.Packet {
	packet := common.GetPacket()
	packet.PacketType = packetType
	packet.Buffer.Write(builder.FinishedBytes())
	return packet
}
