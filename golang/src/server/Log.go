package server

import (
	"bytes"
	"common"
	"fmt"
	"io"
	"os"
	fb "packet"
	"time"
)

const (
	//LogError 심각한 오류
	LogError = 1
	//LogInfo 기본적인 정보
	LogInfo = 2
	//LogDebug 나중에 참고 하면 좋은거
	LogDebug = 3
	//LogTrace 기본 정보
	LogTrace = 4
)

//LogLevel 출력 레벨
var LogLevel int32

func writeBuffer(buffer *bytes.Buffer, level string, manager string, packetType string, v []interface{}) {
	t := time.Now()
	buffer.WriteString(t.Format(time.RFC3339))
	buffer.WriteString(" ")
	buffer.WriteString(level)
	buffer.WriteString(" ")

	buffer.WriteString(manager)
	buffer.WriteString(",")
	buffer.WriteString(packetType)

	for _, value := range v {
		switch t := value.(type) {
		default:
			buffer.WriteString(",")
			buffer.WriteString(fmt.Sprint(t))
		}
	}
}

func writeLog(w io.Writer, level string, manager string, packetType string, v []interface{}) {
	buffer := common.GetBuffer()
	defer common.PutBuffer(buffer)

	writeBuffer(buffer, level, manager, packetType, v)

	fmt.Fprintln(w, buffer.String())
}

//TraceSession TraceSession
func TraceSession(type1 uint32, v ...interface{}) {
	if LogLevel < LogTrace {
		return
	}
	writeLog(os.Stdout, "TRACE", "Session", fb.EnumNamesSessionI[int(type1)], v)
}

//DebugSession DebugSession
func DebugSession(type1 uint32, v ...interface{}) {
	if LogLevel < LogDebug {
		return
	}
	writeLog(os.Stdout, "Debug", "Session", fb.EnumNamesSessionI[int(type1)], v)
}

//TraceClient TraceClient
func TraceClient(type1 uint32, v ...interface{}) {
	if LogLevel < LogTrace {
		return
	}
	writeLog(os.Stdout, "TRACE", "Client", fb.EnumNamesClientI[int(type1)], v)
}

//DebugClient DebugClient
func DebugClient(type1 uint32, v ...interface{}) {
	if LogLevel < LogDebug {
		return
	}
	writeLog(os.Stdout, "Debug", "Client", fb.EnumNamesClientI[int(type1)], v)
}

//TraceServer TraceServer
func TraceServer(type1 uint32, v ...interface{}) {
	if LogLevel < LogTrace {
		return
	}
	writeLog(os.Stdout, "TRACE", "Server", fb.EnumNamesServerI[int(type1)], v)
}

//DebugServer DebugServer
func DebugServer(type1 uint32, v ...interface{}) {
	if LogLevel < LogDebug {
		return
	}
	writeLog(os.Stdout, "Debug", "Server", fb.EnumNamesServerI[int(type1)], v)
}
