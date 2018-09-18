// automatically generated by the FlatBuffers compiler, do not modify

package Client

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type SCMove struct {
	_tab flatbuffers.Table
}

func GetRootAsSCMove(buf []byte, offset flatbuffers.UOffsetT) *SCMove {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &SCMove{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *SCMove) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *SCMove) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *SCMove) User(obj *User) *User {
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

func SCMoveStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func SCMoveAddUser(builder *flatbuffers.Builder, user flatbuffers.UOffsetT) {
	builder.PrependStructSlot(0, flatbuffers.UOffsetT(user), 0)
}
func SCMoveEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}