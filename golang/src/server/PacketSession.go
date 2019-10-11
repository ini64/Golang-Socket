package server

import (
	"common"
	fb "packet"

	flatbuffers "github.com/google/flatbuffers/go"
)

//SessionEnd 세션 종료 처리
func (e *EndPoint) SessionEnd(builder *flatbuffers.Builder, accountID uint32) *common.Packet {
	builder.Reset()

	fb.SessionEndStart(builder)
	fb.SessionEndAddAccountID(builder, accountID)
	packet := fb.SessionEndEnd(builder)
	builder.Finish(packet)

	return getPacket(builder, fb.SessionIEnd)
}

//SessionBegin 서버 입장
func (e *EndPoint) SessionBegin(builder *flatbuffers.Builder, accountID uint32) *common.Packet {
	builder.Reset()

	fb.SessionBeginStart(builder)
	fb.SessionBeginAddAccountID(builder, accountID)
	packet := fb.SessionBeginEnd(builder)
	builder.Finish(packet)

	return getPacket(builder, fb.SessionIBegin)
}
