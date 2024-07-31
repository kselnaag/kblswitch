package svc

import "fmt"

type RingBuff[T any] struct {
	size int
	buff []T
	head int
	tail int
}

func NewRingBuff[T any](size int) *RingBuff[T] {
	return &RingBuff[T]{
		size: size,
		buff: make([]T, size),
		head: 0,
		tail: 0,
	}
}

func (rb *RingBuff[T]) Clear() {
	var val T
	rb.head = 0
	rb.tail = 0
	rb.buff[0] = val
}

func (rb *RingBuff[T]) Size() int {
	return rb.size
}

func (rb *RingBuff[T]) DataLen() int {
	switch {
	case rb.head < rb.tail:
		return rb.tail - rb.head
	case rb.head > rb.tail:
		return rb.size - rb.head + rb.tail
	default:
		return 0
	}
}

func (rb *RingBuff[T]) Set(val T) {
	rb.buff[rb.tail] = val
	rb.tail++
	if rb.tail == rb.size {
		rb.tail = 0
	}
	if rb.tail == rb.head {
		rb.head++
		if rb.head == rb.size {
			rb.head = 0
		}
	}
}

func (rb *RingBuff[T]) Get() T {
	var res T
	if rb.head == rb.tail {
		return res
	}
	res = rb.buff[rb.head]
	rb.head++
	if rb.head == rb.size {
		rb.head = 0
	}
	return res
}

func (rb *RingBuff[T]) Change(pos int, val T) {
	rb.buff[(rb.head+pos)%rb.size] = val
}

func (rb *RingBuff[T]) Read(pos int) T {
	return rb.buff[(rb.head+pos)%rb.size]
}

func (rb *RingBuff[T]) Back() {
	if rb.tail == rb.head {
		return
	}
	rb.tail--
	if rb.tail < 0 {
		rb.tail = rb.size - 1
	}
}

func (rb *RingBuff[T]) Iterate() []T {
	switch {
	case rb.head < rb.tail:
		return rb.buff[rb.head:rb.tail]
	case rb.head > rb.tail:
		return append(rb.buff[rb.head:], rb.buff[:rb.tail]...)
	default:
		return []T{}
	}
}

func (rb *RingBuff[T]) ToString() string {
	return fmt.Sprintf("%c", rb.buff)
}
