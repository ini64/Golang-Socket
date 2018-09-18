// automatically generated by the FlatBuffers compiler, do not modify

package Client

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type CSMove struct {
	_tab flatbuffers.Table
}

func GetRootAsCSMove(buf []byte, offset flatbuffers.UOffsetT) *CSMove {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &CSMove{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *CSMove) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *CSMove) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *CSMove) User(obj *User) *User {
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

func CSMoveStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func CSMoveAddUser(builder *flatbuffers.Builder, user flatbuffers.UOffsetT) {
	builder.PrependStructSlot(0, flatbuffers.UOffsetT(user), 0)
}
func CSMoveEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
