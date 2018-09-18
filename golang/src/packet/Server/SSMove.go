// automatically generated by the FlatBuffers compiler, do not modify

package Server

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type SSMove struct {
	_tab flatbuffers.Table
}

func GetRootAsSSMove(buf []byte, offset flatbuffers.UOffsetT) *SSMove {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &SSMove{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *SSMove) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *SSMove) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *SSMove) User(obj *User) *User {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := o + rcv._tab.Pos
		if obj == nil {
			obj = new(User)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func SSMoveStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func SSMoveAddUser(builder *flatbuffers.Builder, User flatbuffers.UOffsetT) {
	builder.PrependStructSlot(0, flatbuffers.UOffsetT(User), 0)
}
func SSMoveEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}