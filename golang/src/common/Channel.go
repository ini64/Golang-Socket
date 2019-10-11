package common

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

//PoolCounter 카운트 체크
type PoolCounter struct {
	ID uint64
}

//Add 스택추가
func (c *PoolCounter) Add() uint64 {
	id := atomic.AddUint64(&c.ID, 1)

	return id
}

//ChannelCounter 카운트 측정
var ChannelCounter PoolCounter

//ChannelPoolCount 카운트 측정
var ChannelPoolCount int32

//Channel 통신용 채널
type Channel struct {
	channel chan *Packet
	//게임과 상관없는 추적용 일련번호
	ID uint64
	//비지니스 로직의 아이디
	KeyID uint32
	//참조 횟수
	ReferenceCount int32
	//부모를 포인터
	ChannelPool *ChannelPool
}

// ChannelMemoryPool google 제공 pool
var ChannelMemoryPool = sync.Pool{
	New: func() interface{} {
		return new(Channel)
	},
}

//ChannelPool 채널 관리
type ChannelPool struct {
	sync.RWMutex
	Channels map[uint32]*Channel
}

//MakeChannelPool 채널풀 생성
func MakeChannelPool() *ChannelPool {
	return &ChannelPool{
		Channels: make(map[uint32]*Channel),
	}
}

//MakeChannel 채널생성
func MakeChannel(cp *ChannelPool, keyID uint32, size int) *Channel {
	atomic.AddInt32(&ChannelPoolCount, 1)

	channel := ChannelMemoryPool.Get().(*Channel)
	channel.channel = make(chan *Packet, size)
	channel.ID = ChannelCounter.Add()
	channel.KeyID = keyID
	channel.ChannelPool = cp

	return channel
}

//MakeNGet 채널생성
func (cp *ChannelPool) MakeNGet(keyID uint32, size int, referenceCount int32) *Channel {
	cp.Lock()
	defer cp.Unlock()

	channel, ok := cp.Channels[keyID]
	//기존 채널이 있으면 닫는다
	if ok {
		channel.Close()
	}
	//새 채널 할당
	channel = MakeChannel(cp, keyID, size)
	cp.Channels[keyID] = channel

	channel.ReferenceCount += referenceCount
	return channel
}

//Get Channel 얻기
func (cp *ChannelPool) Get(keyID uint32) *Channel {
	cp.Lock()
	defer cp.Unlock()

	channel, ok := cp.Channels[keyID]
	if !ok {
		return nil
	}
	channel.ReferenceCount++
	return channel
}

//Release 메모리 삭제, 풀에 해당 포인터가 들어있는지 리턴
func (c *Channel) Release() bool {
	var activeUser bool
	//fmt.Println("call Release channel", c.ID, c.ReferenceCount, string(stack()))

	if c.ChannelPool != nil {
		c.ChannelPool.Lock()
		defer c.ChannelPool.Unlock()

		channel, ok := c.ChannelPool.Channels[c.KeyID]
		if !ok {
			fmt.Println("not found channel map", c.ID, string(Stack()))
			return activeUser
		}

		if channel.ID == c.ID {
			activeUser = true
		}

		c.ReferenceCount--

		if c.ReferenceCount < 0 {
			fmt.Println("invalid Put")
			return activeUser
		}

		if c.ReferenceCount == 0 {

			//남은 버퍼 비우기
			close(c.channel)

			for data := range c.channel {
				if data != nil {
					PutPacket(data)
				}
			}

			if activeUser {
				delete(c.ChannelPool.Channels, c.KeyID)
			}

			atomic.AddInt32(&ChannelPoolCount, -1)
			ChannelMemoryPool.Put(c)

		}
	} else {
		close(c.channel)

		for packet := range c.channel {
			if packet != nil {
				PutPacket(packet)
			}
		}
		atomic.AddInt32(&ChannelPoolCount, -1)
		ChannelMemoryPool.Put(c)
	}

	return activeUser
}

//Close 닫기
func (c *Channel) Close() {
	c.channel <- nil
}

//Get 데이터 얻기
func (c *Channel) Get() *Packet {
	return <-c.channel
}

//GetTimeOut 데이터 얻기
func (c *Channel) GetTimeOut(timeOut int) *Packet {
	select {
	case <-time.After(time.Duration(timeOut) * time.Millisecond):
		return nil
	case packet := <-c.channel:
		return packet
	}
	return nil
}

//GetDefault 데이터 얻기
func (c *Channel) GetDefault() *Packet {
	select {
	case packet := <-c.channel:
		return packet
	default:
		return nil
	}
	return nil
}

//Put 추가
func (c *Channel) Put(packet *Packet) bool {
	c.channel <- packet
	return true
}
