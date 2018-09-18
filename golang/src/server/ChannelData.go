package server

import (
	"sync"
)

//ChannelData 통신용 데이터
type ChannelData struct {
	Packet  *Packet
	Channel Channel
}

//Reset 초기화
func (d *ChannelData) Reset() {
	if d.Packet != nil {
		PutPacket(d.Packet)
		d.Packet = nil
	}
}

//ChannelDataPool 네트워크 전송용 풀
var ChannelDataPool = sync.Pool{
	New: func() interface{} {
		return new(ChannelData)
	},
}

//GetChannelData 얻기
func GetChannelData() *ChannelData {
	data := ChannelDataPool.Get().(*ChannelData)
	return data
}

//PutChannelData 반납
func PutChannelData(d *ChannelData) {
	if d != nil {
		d.Reset()
		ChannelDataPool.Put(d)
	}
}
