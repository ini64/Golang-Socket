package server

import "sync"

//Vec3 위치
type Vec3 struct {
	x float32
	y float32
	z float32
}

//Qua 회전
type Qua struct {
	x float32
	y float32
	z float32
	w float32
}

//User 유저 정보
type User struct {
	AccountID int32
	Pos       Vec3
	Rot       Qua
}

//Reset 초기화
func (u *User) Reset() {
	u.AccountID = 0

	u.Pos.x = 0
	u.Pos.y = 0
	u.Pos.z = 0

	u.Rot.x = 0
	u.Rot.y = 0
	u.Rot.z = 0
	u.Rot.w = 0
}

//UserPool 네트워크 전송용 풀
var UserPool = sync.Pool{
	New: func() interface{} {
		return new(User)
	},
}

//GetUser 얻기
func GetUser() *User {
	return UserPool.Get().(*User)
}

//PutUser 반납
func PutUser(u *User) {
	if u != nil {
		u.Reset()
		UserPool.Put(u)
	}
}
