package genid

import "sync/atomic"

type IDcounter struct {
	id int64
}

func NewCounter() *IDcounter {
	return &IDcounter{}
}

func (i *IDcounter) GenerateID() int {
	return int(atomic.AddInt64(&i.id, 1))
}
