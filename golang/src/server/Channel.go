package server

///////////////////////////////////////////////////////////////////////////////////////
//Channel
///////////////////////////////////////////////////////////////////////////////////////

//Channel 통신용 채널
type Channel chan *ChannelData

//Add 작업 추가
func (c Channel) Add(packet *Packet) {
	data := GetChannelData()
	data.Packet = packet
	c <- data
}

//AddChannel 통신채널 넣기
func (c Channel) AddChannel(packet *Packet, channel Channel) {
	data := GetChannelData()
	data.Packet = packet
	data.Channel = channel
	c <- data
}

//Close 닫기
func (c Channel) Close() {
	c <- nil
}
