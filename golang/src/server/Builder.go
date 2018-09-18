package server

import (
	"sync"

	flatbuffers "github.com/google/flatbuffers/go"
)

//BuilderPool 빌더 빌더
var BuilderPool = sync.Pool{
	New: func() interface{} {
		return flatbuffers.NewBuilder(1024)
	},
}

//GetBuilder 얻기
func GetBuilder() *flatbuffers.Builder {
	builder := BuilderPool.Get().(*flatbuffers.Builder)
	return builder
}

//PutBuilder 반납
func PutBuilder(builder *flatbuffers.Builder) {
	if builder != nil {
		builder.Reset()
		BuilderPool.Put(builder)
	}
}
