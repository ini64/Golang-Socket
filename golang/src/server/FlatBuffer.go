package server

import (
	cp "packet/Client"
	sp "packet/Server"

	flatbuffers "github.com/google/flatbuffers/go"
)

//SSEnter 서버 입장
func (e *EndPoint) SSEnter(builder *flatbuffers.Builder, csLogin *cp.CSLogin) *Packet {
	user := csLogin.User(nil)
	pos := user.Pos(nil)
	rot := user.Rot(nil)

	builder.Reset()

	sp.SSEnterStart(builder)
	sp.SSEnterAddAccountID(builder, user.Num())

	vec3 := sp.CreateVec3(builder, pos.X(), pos.Y(), pos.Z())
	sp.SSEnterAddPos(builder, vec3)

	qua := sp.CreateQua(builder, rot.X(), rot.Y(), rot.Z(), rot.W())
	sp.SSEnterAddRot(builder, qua)

	sp.SSEnterAddTime(builder, user.Time())
	packet := sp.SSEnterEnd(builder)
	builder.Finish(packet)

	return GetPacket(sp.SSPacketTypeEnter, builder)
}

//SSEnterAck 서버 입장
func (d *GameManagerData) SSEnterAck(builder *flatbuffers.Builder) *Packet {
	builder.Reset()

	length := len(d.Users)
	sp.SSEnterAckStartUsersVector(builder, length)
	for _, user := range d.Users {
		sp.CreateUser(builder, user.AccountID, user.Pos.x, user.Pos.y, user.Pos.z, user.Rot.x, user.Rot.y, user.Rot.z, user.Rot.w, user.Time)
	}
	offset := builder.EndVector(length)

	sp.SSEnterAckStart(builder)
	sp.SSEnterAckAddUsers(builder, offset)
	packet := sp.SSEnterAckEnd(builder)
	builder.Finish(packet)

	return GetPacket(sp.SSPacketTypeEnterAck, builder)
}

//SSEnterNak 서버 입장
func (d *GameManagerData) SSEnterNak(builder *flatbuffers.Builder, error int32) *Packet {
	builder.Reset()

	sp.SSEnterAckStart(builder)
	sp.SSEnterAckAddError(builder, error)
	packet := sp.SSEnterAckEnd(builder)
	builder.Finish(packet)

	return GetPacket(sp.SSPacketTypeEnterAck, builder)
}

//SSLeave 서버 퇴장
func (e *EndPoint) SSLeave(builder *flatbuffers.Builder, accountID int32) *Packet {
	builder.Reset()

	sp.SSLeaveStart(builder)
	sp.SSLeaveAddAccountID(builder, accountID)
	packet := sp.SSLeaveEnd(builder)
	builder.Finish(packet)

	return GetPacket(sp.SSPacketTypeLeave, builder)
}

//SSClose 연결 끊어짐
func (e *EndPoint) SSClose(builder *flatbuffers.Builder) *Packet {
	builder.Reset()

	sp.SSCloseStart(builder)
	packet := sp.SSCloseEnd(builder)
	builder.Finish(packet)

	return GetPacket(sp.SSPacketTypeClose, builder)
}

//SCLogin 로그인 응답
func (d *UserManagerData) SCLogin(builder *flatbuffers.Builder, ssEnterAck *sp.SSEnterAck) *Packet {
	builder.Reset()

	if ssEnterAck.Error() != 0 {
		cp.SCLoginStart(builder)
		cp.SCLoginAddError(builder, ssEnterAck.Error())
	} else {

		length := ssEnterAck.UsersLength()

		var user sp.User
		var vec3 sp.Vec3
		var qua sp.Qua

		cp.SCLoginStartUsersVector(builder, length)
		for i := 0; i < length; i++ {

			ssEnterAck.Users(&user, i)
			user.Pos(&vec3)
			user.Rot(&qua)
			cp.CreateUser(builder, user.Num(), vec3.X(), vec3.Y(), vec3.Z(), qua.X(), qua.Y(), qua.Z(), qua.W(), user.Time())
		}
		offset := builder.EndVector(length)

		cp.SCLoginStart(builder)
		cp.SCLoginAddUsers(builder, offset)

	}
	packet := cp.SCLoginEnd(builder)
	builder.Finish(packet)

	return GetPacket(cp.SCPacketTypeLogin, builder)
}

//SSEnterOther 다른유저 퇴장
func (d *GameManagerData) SSEnterOther(builder *flatbuffers.Builder, user *User) *Packet {
	builder.Reset()
	sp.SSEnterOtherStart(builder)
	offset := sp.CreateUser(builder, user.AccountID, user.Pos.x, user.Pos.y, user.Pos.z, user.Rot.x, user.Rot.y, user.Rot.z, user.Rot.w, user.Time)
	sp.SSEnterOtherAddUser(builder, offset)
	packet := sp.SSEnterOtherEnd(builder)
	builder.Finish(packet)

	return GetPacket(sp.SSPacketTypeEnterOther, builder)
}

//SSLeaveOther 다른유저 퇴장
func (d *GameManagerData) SSLeaveOther(builder *flatbuffers.Builder, accountID int32) *Packet {
	builder.Reset()
	sp.SSLeaveOtherStart(builder)
	sp.SSLeaveOtherAddAccountID(builder, accountID)
	packet := sp.SSLeaveOtherEnd(builder)
	builder.Finish(packet)

	return GetPacket(sp.SSPacketTypeLeaveOther, builder)
}

//SCEnterOther 다른유저 입장
func (d *UserManagerData) SCEnterOther(builder *flatbuffers.Builder, other *sp.SSEnterOther) *Packet {
	builder.Reset()

	user := &sp.User{}
	other.User(user)

	vec3 := &sp.Vec3{}
	qua := &sp.Qua{}

	user.Pos(vec3)
	user.Rot(qua)

	cp.SCEnterOtherStart(builder)
	offset := cp.CreateUser(builder, user.Num(), vec3.X(), vec3.Y(), vec3.Z(), qua.X(), qua.Y(), qua.Z(), qua.W(), user.Time())
	cp.SCEnterAddUser(builder, offset)
	packet := cp.SCEnterOtherEnd(builder)
	builder.Finish(packet)

	return GetPacket(cp.SCPacketTypeEnterOther, builder)
}

//SCLeaveOther 다른유저 입장
func (d *UserManagerData) SCLeaveOther(builder *flatbuffers.Builder, accountID int32) *Packet {
	builder.Reset()

	cp.SCLeaveOtherStart(builder)
	cp.SCLeaveOtherAddAccountID(builder, accountID)
	packet := cp.SCLeaveOtherEnd(builder)
	builder.Finish(packet)

	return GetPacket(cp.SCPacketTypeLeaveOther, builder)
}

//SSMove 이동
func (d *UserManagerData) SSMove(builder *flatbuffers.Builder, move *cp.CSMove) *Packet {
	builder.Reset()

	user := &cp.User{}
	move.User(user)

	vec3 := &cp.Vec3{}
	qua := &cp.Qua{}

	user.Pos(vec3)
	user.Rot(qua)

	sp.SSMoveStart(builder)
	offset := sp.CreateUser(builder, user.Num(), vec3.X(), vec3.Y(), vec3.Z(), qua.X(), qua.Y(), qua.Z(), qua.W(), user.Time())
	sp.SSMoveAddUser(builder, offset)
	packet := sp.SSMoveEnd(builder)
	builder.Finish(packet)

	return GetPacket(sp.SSPacketTypeMove, builder)
}

//SSMoveAck 이동 처리
func (d *GameManagerData) SSMoveAck(builder *flatbuffers.Builder, user *User) *Packet {
	builder.Reset()
	sp.SSMoveAckStart(builder)
	offset := sp.CreateUser(builder, user.AccountID, user.Pos.x, user.Pos.y, user.Pos.z, user.Rot.x, user.Rot.y, user.Rot.z, user.Rot.w, user.Time)
	sp.SSMoveAckAddUser(builder, offset)
	packet := sp.SSMoveAckEnd(builder)
	builder.Finish(packet)

	return GetPacket(sp.SSPacketTypeMoveAck, builder)
}

//SCMove 이동 처리
func (d *UserManagerData) SCMove(builder *flatbuffers.Builder, move *sp.SSMoveAck) *Packet {
	builder.Reset()

	user := &sp.User{}
	move.User(user)

	vec3 := &sp.Vec3{}
	qua := &sp.Qua{}

	user.Pos(vec3)
	user.Rot(qua)

	cp.SCMoveStart(builder)
	offset := cp.CreateUser(builder, user.Num(), vec3.X(), vec3.Y(), vec3.Z(), qua.X(), qua.Y(), qua.Z(), qua.W(), user.Time())
	cp.SCMoveAddUser(builder, offset)
	packet := cp.SCMoveEnd(builder)
	builder.Finish(packet)

	return GetPacket(cp.SCPacketTypeMove, builder)
}
