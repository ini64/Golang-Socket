package client

import (
	"common"
	"fmt"
	fb "packet"
	"server"
	"sync"
	"time"
)

//PerfomanceTest 퍼포먼스 체크용
func PerfomanceTest(filePath string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("run tim panic", err, string(common.Stack()))
		}
	}()

	conf := &Conf{}
	if !conf.Load(filePath) {
		return
	}

	endPoint := server.Main("conf/server.json", "")
	time.Sleep(3 * time.Second)

	channel := make(chan bool, conf.ConcurrentAccess)

	var run int32
	run = 1

	go Monitor(run)

	var wg sync.WaitGroup

	for i := 1; i < (conf.Total + 1); i++ {
		input := &Inputs{
			AccountID: uint32(i),
			LoopCount: 10,
			Conf:      conf,
			wg:        &wg,
		}
		wg.Add(1)

		channel <- true

		go PerfomanceTestJob(input)
	}
	wg.Wait()

	endPoint.TCPConn.Close()
	endPoint.UDPConn.Close()
	endPoint.SessionChannel.Close()
	endPoint.WG.Wait()

	time.Sleep(3 * time.Second)
}

//PerfomanceTestJob 서브 함수
func PerfomanceTestJob(input *Inputs) {
	defer input.wg.Done()

	conn, tcpRead, tcpWrite := TCPConnect(input.Conf)
	if conn == nil {
		return
	}

	builder := common.GetBuilder()
	defer common.PutBuilder(builder)

	clientBegin := ClientBegin(builder, input.AccountID)
	if !tcpWrite(clientBegin) {
		return
	}

	packet := TCPWait(fb.ServerIBegin, 10, tcpRead)
	if packet == nil {
		return
	}
	common.PutPacket(packet)

	conn.Close()
}
