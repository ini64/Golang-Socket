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
	Time      float32
}

//Reset 초기화
func (u *User) Reset() {

}

//UserPool 네트워크 전송용 풀
var UserPool = sync.Pool{
	New: func() interface{} {
		return new(User)
	},
}

//GetUser 얻기
func GetUser() *User {
	data := UserPool.Get().(*User)
	return data
}

//PutUser 반납
func PutUser(d *User) {
	if d != nil {
		d.Reset()
		UserPool.Put(d)
	}
}
