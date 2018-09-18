package server

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
)

//EndPoint 전체 정보
type EndPoint struct {
	TCPConn *net.TCPListener
	Conf
	GameManagerChannel Channel
}

//NewEndPoint 서버 기본 정보 가져오기
func NewEndPoint(fileName string) *EndPoint {

	endPoint := &EndPoint{
		GameManagerChannel: make(Channel, 1024),
	}

	file, _ := os.Open(fileName)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&endPoint.Conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, "not found .conf file:", err.Error())
		return nil
	}

	return endPoint
}

//Main 서버 몸체
func Main(conf string, version string) {
	var wg sync.WaitGroup

	e := NewEndPoint(conf)

	wg.Add(1)
	go e.TCPListener(&wg)

	wg.Add(1)
	go e.GameManager(&wg)

	wg.Wait()

}
