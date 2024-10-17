package filtering

import (
	"sync"
	"time"
)

var (
	BufferSize      int           = 3
	TimeBufferClear time.Duration = 5 * time.Second
)

type CircularBuffer struct {
	mu       sync.Mutex
	buffer   []int
	position int
	size     int
}

func NewCircularBuffer() *CircularBuffer {
	return &CircularBuffer{
		mu:       sync.Mutex{},
		buffer:   make([]int, BufferSize),
		position: -1,
		size:     BufferSize,
	}
}

func (c *CircularBuffer) Push(numb int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.position == c.size-1 {
		for i := 1; i <= c.size-1; i++ {
			c.buffer[i-1] = c.buffer[i]
		}
		c.buffer[c.position] = numb
	} else {
		c.position++
		c.buffer[c.position] = numb
	}
}

func (c *CircularBuffer) Get() []int {
	if c.position < 0 {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	numbOut := c.buffer[:c.position+1]
	c.position = -1
	return numbOut
}

func FilterNegative(numbIn, numbOut chan int) {
	for num := range numbIn {
		if num > 0 {
			numbOut <- num
		}
	}
}

func FilterNumberNotMultipleThree(numbIn, numbOut chan int) {
	for numb := range numbIn {
		if numb != 0 && numb%3 == 0 {
			numbOut <- numb
		}
	}

}

func Buffering(numbIn, numbOut chan int) {

	circBuff := NewCircularBuffer()
	for {
		if circBuff.size != BufferSize {
			circBuff.size = BufferSize
			circBuff.buffer = make([]int, BufferSize)
		}
		select {
		case numb := <-numbIn:
			circBuff.Push(numb)
		case <-time.After(TimeBufferClear):
			cl := circBuff.Get()
			for _, num := range cl {
				numbOut <- num
			}
		}
	}

}
