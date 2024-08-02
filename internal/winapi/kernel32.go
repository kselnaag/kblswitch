package winapi

import (
	T "kblswitch/internal/types"
	"syscall"
)

var (
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	procGetCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")
	procExitThread         = kernel32.NewProc("ExitThread")
)

type Kernel32 struct {
	log          T.ILog
	CurrThreadId uintptr
}

func NewKernel32(log T.ILog) *Kernel32 {
	return &Kernel32{
		log:          log,
		CurrThreadId: 0,
	}
}

func (k *Kernel32) ExitThread(code uintptr) {
	_, _, e := procExitThread.Call()
	if e != nil && e.Error() != T.ErrSuccessfull {
		k.log.LogError(e, "kernel32.GetCurrentThreadId() error")
	}
}

func (k *Kernel32) GetCurrentThreadId() uintptr {
	t1, _, e := procGetCurrentThreadId.Call()
	k.CurrThreadId = t1
	if e != nil && e.Error() != T.ErrSuccessfull {
		k.log.LogError(e, "kernel32.GetCurrentThreadId() error")
	}
	return t1
}
