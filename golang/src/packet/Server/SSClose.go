// automatically generated by the FlatBuffers compiler, do not modify

package Server

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type SSClose struct {
	_tab flatbuffers.Table
}

func GetRootAsSSClose(buf []byte, offset flatbuffers.UOffsetT) *SSClose {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &SSClose{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *SSClose) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *SSClose) Table() flatbuffers.Table {
	return rcv._tab
}

func SSCloseStart(builder *flatbuffers.Builder) {
	builder.StartObject(0)
}
func SSCloseEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}