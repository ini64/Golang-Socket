package common

import (
	"bytes"
	"sync"
	"sync/atomic"
)

//BufferPoolCount 갯수
var BufferPoolCount int64

//BufferPool 버퍼
var BufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

//GetBuffer 얻기
func GetBuffer() *bytes.Buffer {
	buffer := BufferPool.Get().(*bytes.Buffer)
	atomic.AddInt64(&BufferPoolCount, 1)
	return buffer
}

//PutBuffer 반납
func PutBuffer(buf *bytes.Buffer) {
	if buf != nil {
		buf.Reset()
		BufferPool.Put(buf)
		atomic.AddInt64(&BufferPoolCount, -1)
	}
}
