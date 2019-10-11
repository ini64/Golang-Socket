package server

import (
	"common"
	"conf"
	"fmt"
	"net"
	"sync"
	"time"
)

//EndPoint 전체 정보
type EndPoint struct {
	conf.Conf
	TCPConn        *net.TCPListener
	UDPConn        *net.UDPConn
	SessionChannel *common.Channel

	TCPWritePool *common.ChannelPool
	TCPReadPool  *common.ChannelPool
	UDPWritePool *common.ChannelPool

	WG sync.WaitGroup

	UDPRecv   int64
	TCPRecv   int64
	TCPSend   int64
	UDPSend   int64
	UserCount int64
}

//GetRecvTimeOut 설정
func (e *EndPoint) GetRecvTimeOut() time.Time {
	return time.Now().Add(time.Duration(e.Conf.ReadTimeOutSeconds) * time.Second)
}

//GetSendTimeOut 설정
func (e *EndPoint) GetSendTimeOut() time.Time {
	return time.Now().Add(time.Duration(e.Conf.WriteTimeOutSeconds) * time.Second)
}

//NewEndPoint 서버 기본 정보 가져오기
func NewEndPoint(fileName string) *EndPoint {

	endPoint := &EndPoint{
		SessionChannel: common.MakeChannel(nil, 0, 1024),
		TCPWritePool:   common.MakeChannelPool(),
		TCPReadPool:    common.MakeChannelPool(),
		UDPWritePool:   common.MakeChannelPool(),
	}

	if !endPoint.Load(fileName) {
		fmt.Println(fileName + "not found")
	}

	switch endPoint.LogLevel {
	case "ERROR":
		LogLevel = LogError
	case "INFO":
		LogLevel = LogInfo
	case "DEBUG":
		LogLevel = LogDebug
	case "TRACE":
		LogLevel = LogTrace
	}

	return endPoint
}

//Main 서버 몸체
func Main(conf string, version string) *EndPoint {
	e := NewEndPoint(conf)

	e.WG.Add(3)
	go e.TCPListener(&e.WG)
	go e.UDPListener(&e.WG)
	go e.SessionManager(&e.WG)

	return e

}
