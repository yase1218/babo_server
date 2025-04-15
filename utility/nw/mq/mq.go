package mq

import "container/list"

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// block

type DefaultMsgChan struct {
	in   chan<- []byte
	out  <-chan []byte
	done chan struct{}
	size int
}

func NewDefaultQueue(size int, done chan struct{}) *DefaultMsgChan {
	if size < 10 {
		size = 10
	}
	msgChan := make(chan []byte, size)
	return &DefaultMsgChan{
		in:   msgChan,
		out:  msgChan,
		done: done,
		size: size,
	}
}

func (this *DefaultMsgChan) GetInChan() chan<- []byte {
	return this.in
}

func (this *DefaultMsgChan) GetOutChan() <-chan []byte {
	return this.out
}

func (this *DefaultMsgChan) Len() int {
	return len(this.in)
}

func (this *DefaultMsgChan) Size() int {
	return this.size
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// unblock

type NonBlockingMsgChan struct {
	*DefaultMsgChan
	elements *list.List
	count    int
}

func NewNonBlockingMsgChan(size int, done chan struct{}) *NonBlockingMsgChan {
	if size < 10 {
		size = 10
	}
	in := make(chan []byte, size)
	out := make(chan []byte, size)
	mq := &NonBlockingMsgChan{
		DefaultMsgChan: &DefaultMsgChan{in: in, out: out, done: done, size: size},
		elements:       list.New(),
		count:          0,
	}

	go mq.run(in, out)
	return mq
}

func (this *NonBlockingMsgChan) run(in <-chan []byte, out chan<- []byte) {
	for {
		if this.in == nil && this.elements.Len() == 0 {
			close(out)
			break
		}
		var (
			outChan chan<- []byte
			outVal  []byte
		)
		if this.elements.Len() > 0 {
			outChan = out
			outVal = this.elements.Front().Value.([]byte)
		}
		select {
		case i, ok := <-in:
			if ok {
				this.elements.PushBack(i)
				this.count++
			} else {
				in = nil
			}
		case outChan <- outVal:
			this.elements.Remove(this.elements.Front())
			this.count--
		case <-this.done:
			return
		}
	}
}

func (this *NonBlockingMsgChan) Len() int {
	return this.count
}

func (this *NonBlockingMsgChan) Size() int {
	return this.size
}
