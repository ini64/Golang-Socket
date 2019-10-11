package client

import (
	"common"
	fb "packet"

	flatbuffers "github.com/google/flatbuffers/go"
)

//ClientBegin 로그인 응답
func ClientBegin(builder *flatbuffers.Builder, accountID uint32) *common.Packet {
	builder.Reset()

	fb.ClientBeginStart(builder)
	fb.ClientBeginAddAccountID(builder, accountID)
	packet := fb.ClientBeginEnd(builder)
	builder.Finish(packet)

	return getPacket(builder, fb.ClientIBegin)
}
