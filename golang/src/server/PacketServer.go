package server

import (
	"common"
	fb "packet"

	flatbuffers "github.com/google/flatbuffers/go"
)

//ServerBegin 로그인 응답
func (e *EndPoint) ServerBegin(builder *flatbuffers.Builder) *common.Packet {
	builder.Reset()

	fb.ServerBeginStart(builder)
	packet := fb.ServerBeginEnd(builder)
	builder.Finish(packet)

	return getPacket(builder, fb.ServerIBegin)
}
