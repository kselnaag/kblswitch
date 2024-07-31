package svc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	S "kblswitch/internal/svc"
)

func TestSvc(t *testing.T) {
	asrt := assert.New(t)
	defer func() {
		err := recover()
		asrt.Nil(err)
	}()

	t.Run("RingBuff", func(t *testing.T) {
		rb := S.NewRingBuff[byte](10)
		asrt.Equal(rb.Size(), 10)
		asrt.Equal(rb.DataLen(), 0)

		rb.Set(1)
		rb.Set(2)
		rb.Set(3)
		rb.Set(4)
		rb.Set(5)
		rb.Set(6)
		asrt.Equal(rb.DataLen(), 6)
		asrt.Equal(rb.ToString(), "[1 2 3 4 5 6 0 0 0 0]")
		b := rb.Get()
		b = rb.Get()
		b = rb.Get()
		b = rb.Get()
		b = rb.Get()
		b = rb.Get()
		asrt.Equal(b, byte(6))
		asrt.Equal(rb.DataLen(), 0)

		rb.Set(1)
		rb.Set(2)
		rb.Set(3)
		rb.Set(4)
		rb.Set(5)
		rb.Set(6)
		asrt.Equal(rb.DataLen(), 6)
		asrt.Equal(rb.ToString(), "[5 6 3 4 5 6 1 2 3 4]")
		b = rb.Get()
		b = rb.Get()
		b = rb.Get()
		b = rb.Get()
		b = rb.Get()
		b = rb.Get()
		asrt.Equal(b, byte(6))
		asrt.Equal(rb.DataLen(), 0)

	})

}
