package user32

import (
	"syscall"
	"unsafe"

	T "kblswitch/internal/types"
)

var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procGetKeyboardLayout        = user32.NewProc("GetKeyboardLayout")
	procPostMessageA             = user32.NewProc("PostMessageA")
	procSendMessageA             = user32.NewProc("SendMessageA")
	procSendInput                = user32.NewProc("SendInput")
	procSetWindowsHookExA        = user32.NewProc("SetWindowsHookExA")
	procUnhookWindowsHookEx      = user32.NewProc("UnhookWindowsHookEx")
	procCallNextHookEx           = user32.NewProc("CallNextHookEx")
	procGetMessageA              = user32.NewProc("GetMessageA")
	procTranslateMessage         = user32.NewProc("TranslateMessage")
	procDispatchMessage          = user32.NewProc("DispatchMessage")
	procGetKeyboardState         = user32.NewProc("GetKeyboardState")
	procToUnicodeEx              = user32.NewProc("ToUnicodeEx")
	procGetKeyState              = user32.NewProc("GetKeyState")
)

const ErrSuccessfull = "The operation completed successfully."

type User32 struct {
	log            T.ILog
	KeyboardHookId uintptr
}

func NewUser32(log T.ILog) *User32 {
	return &User32{
		log:            log,
		KeyboardHookId: 0,
	}
}

func (u *User32) GetForegroundWindow() uintptr {
	w1, _, e := procGetForegroundWindow.Call(0)
	if e != nil && e.Error() != ErrSuccessfull {
		u.log.LogError(e, "user32.GetForegroundWindow() error")
	}
	return w1
}

func (u *User32) GetWindowThreadProcessId(hwind uintptr) uintptr {
	p1, _, e := procGetWindowThreadProcessId.Call(hwind)
	if e != nil && e.Error() != ErrSuccessfull {
		u.log.LogError(e, "user32.GetWindowThreadProcessId() error")
	}
	return p1
}

func (u *User32) GetKeyboardLayout(proc uintptr) uintptr {
	l1, _, e := procGetKeyboardLayout.Call(proc)
	if e != nil && e.Error() != ErrSuccessfull {
		u.log.LogError(e, "user32.GetKeyboardLayout() error")
	}
	return l1
}

func (u *User32) GetKeyboardState(keyState *[256]byte) bool {
	ks1, _, e := procGetKeyboardState.Call(uintptr(unsafe.Pointer(keyState)))
	if e != nil && e.Error() != ErrSuccessfull {
		u.log.LogError(e, "user32.GetKeyboardState() error")
	}
	return (ks1 != 0)
}

func (u *User32) GetKeyState(keyCode int) uint16 {
	s1, _, e := procGetKeyState.Call(uintptr(keyCode))
	if e != nil && e.Error() != ErrSuccessfull {
		u.log.LogError(e, "user32.GetKeyState() error")
	}
	return uint16(s1)
}

func (u *User32) ToUnicodeEx(wVirtKey uint16) {

}

func (u *User32) PostMessageA(hwnd uintptr, msg int, wparam uintptr, lparam uintptr) bool {
	p1, _, e := procPostMessageA.Call(hwnd, uintptr(msg), wparam, lparam)
	if e != nil && e.Error() != ErrSuccessfull {
		u.log.LogError(e, "user32.PostMessageA() error")
	}
	return (p1 != 0)
}

func (u *User32) SendMessageA(hwnd uintptr, msg int, wparam uintptr, lparam uintptr) uintptr {
	s1, _, e := procSendMessageA.Call(hwnd, uintptr(msg), wparam, lparam)
	if e != nil && e.Error() != ErrSuccessfull {
		u.log.LogError(e, "user32.SendMessageA() error")
	}
	return s1
}

func (u *User32) SendInput(cinputs int, pinput *T.Input, cbSize T.Input) uint {
	i1, _, e := procSendInput.Call(uintptr(cinputs), uintptr(unsafe.Pointer(pinput)), uintptr(unsafe.Sizeof(cbSize)))
	if e != nil && e.Error() != ErrSuccessfull {
		u.log.LogError(e, "user32.SendInput() error")
	}
	return uint(i1)
}

func (u *User32) SetWindowsHookExA(idHook int, lpfn T.HOOKPROC, hMod int, dwThreadId int) uintptr {
	h1, _, e := procSetWindowsHookExA.Call(uintptr(idHook), syscall.NewCallback(lpfn), uintptr(hMod), uintptr(dwThreadId))
	u.KeyboardHookId = h1
	if e != nil && e.Error() != ErrSuccessfull {
		u.log.LogError(e, "user32.SetWindowsHookExA() error")
	}
	return h1
}

func (u *User32) UnhookWindowsHookEx(keyboardHookId uintptr) bool {
	u1, _, e := procUnhookWindowsHookEx.Call(keyboardHookId)
	if e != nil && e.Error() != ErrSuccessfull {
		u.log.LogError(e, "user32.UnhookWindowsHookEx() error")
	}
	return (u1 != 0)
}

func (u *User32) CallNextHookEx(keyboardHookId uintptr, nCode int, wparam uintptr, lparam uintptr) uintptr {
	n1, _, e := procCallNextHookEx.Call(keyboardHookId, uintptr(nCode), wparam, lparam)
	if e != nil && e.Error() != ErrSuccessfull {
		u.log.LogError(e, "user32.CallNextHookEx() error")
	}
	return n1
}

func (u *User32) GetMessageA(msg *T.MSG, hWnd uint32, wMsgFilterMin uintptr, wMsgFilterMax uintptr) bool {
	m1, _, e := procGetMessageA.Call(uintptr(unsafe.Pointer(msg)), uintptr(hWnd), wMsgFilterMin, wMsgFilterMax)
	if e != nil && e.Error() != ErrSuccessfull {
		u.log.LogError(e, "user32.GetMessageA() error")
	}
	return (m1 != 0)
}

func (u *User32) TranslateMessage(msg *T.MSG) bool {
	t1, _, e := procTranslateMessage.Call(uintptr(unsafe.Pointer(msg)))
	if e != nil && e.Error() != ErrSuccessfull {
		u.log.LogError(e, "user32.TranslateMessage() error")
	}
	return (t1 != 0)
}

func (u *User32) DispatchMessage(msg *T.MSG) uintptr {
	d1, _, e := procDispatchMessage.Call(uintptr(unsafe.Pointer(msg)))
	if e != nil && e.Error() != ErrSuccessfull {
		u.log.LogError(e, "user32.DispatchMessage() error")
	}
	return d1
}
