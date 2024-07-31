package types

const (
	INPUT_KEYBOARD            = 1
	WM_INPUTLANGCHANGEREQUEST = 80
	WH_KEYBOARD_LL            = 13
	WM_KEYDOWN                = 256
	WM_KEYUP                  = 257
	WM_CHAR                   = 258
	VK_PAUSE                  = 19
	VK_BACK                   = 8
	VK_ENTER                  = 13
	VK_LSHIFT                 = 160
	VK_RSHIFT                 = 161
	KBLayoutRus               = 68748313
	KBLayoutEng               = 67699721
	KEYEVENTF_KEYUP           = 2
)

type Input struct {
	InputType uint32
	Ki        KBDLLHOOKSTRUCT
	Padding   uint64
}

type KBDLLHOOKSTRUCT struct {
	VkCode      uint16
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

type MSG struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

type POINT struct {
	X, Y int32
}

type HOOKPROC func(int, uintptr, uintptr) uintptr

type (
	DWORD     uint32
	WPARAM    uintptr
	LPARAM    uintptr
	LRESULT   uintptr
	HANDLE    uintptr
	HINSTANCE HANDLE
	HHOOK     HANDLE
	HWND      HANDLE
)
