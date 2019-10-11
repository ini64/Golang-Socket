package client

import (
	"sync"
	"testing"
)

//Inputs 테스트에 필요한 변수
type Inputs struct {
	*testing.T
	wg        *sync.WaitGroup
	AccountID uint32
	LoopCount int
	Conf      *Conf
}
