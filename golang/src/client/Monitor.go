package client

import (
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

//UserCount 현 유저
var UserCount int32

//Monitor 모니터
func Monitor(run int32) {

	//var TCPCount, TCPSize, UDPCount, UDPSize int64
	var TCPCount, UDPCount int64
	var m runtime.MemStats

	for {
		select {
		case <-time.After(1 * time.Second):

			tcpRecvCount := atomic.LoadInt64(&TCPReadCount)
			udpRecvCount := atomic.LoadInt64(&UDPReadCount)

			tcount := tcpRecvCount - TCPCount
			TCPCount = tcpRecvCount

			ucount := udpRecvCount - UDPCount
			UDPCount = udpRecvCount

			t := time.Now()

			runtime.ReadMemStats(&m)

			fmt.Fprintln(os.Stdout, t.Format(time.RFC3339), ",", runtime.NumGoroutine(), ",", tcount, ",", ucount)
		}
	}
}
