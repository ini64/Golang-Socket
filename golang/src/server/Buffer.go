package server

import (
	"bytes"
	"sync"
)

//BufferPool 버퍼
var BufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

//GetBuffer 얻기
func GetBuffer() *bytes.Buffer {
	buffer := BufferPool.Get().(*bytes.Buffer)
	return buffer
}

//PutBuffer 반납
func PutBuffer(buf *bytes.Buffer) {
	if buf != nil {
		buf.Reset()
		BufferPool.Put(buf)
	}
}
