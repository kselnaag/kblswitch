package types

const (
	INPUT_KEYBOARD            = 1
	WM_INPUTLANGCHANGEREQUEST = 80
	WM_QUIT                   = 18
	WH_KEYBOARD_LL            = 13
	WM_KEYDOWN                = 256
	WM_KEYUP                  = 257
	WM_CHAR                   = 258
	VK_PAUSE                  = 19
	VK_BACK                   = 8
	VK_ENTER                  = 13
	VK_SHIFT                  = 16
	VK_CONTROL                = 17
	VK_ESCAPE                 = 27
	VK_LSHIFT                 = 160
	VK_RSHIFT                 = 161
	KBLayoutRus               = 68748313
	KBLayoutEng               = 67699721
	KEYEVENTF_KEYUP           = 2
	KEYEVENTF_UNICODE         = 4
)

const ErrSuccessfull = "The operation completed successfully."

type Input struct {
	InputType uint32
	Ki        KBDLLHOOKSTRUCT
	Padding   uint64
}

type KBDLLHOOKSTRUCT struct {
	VkCode      uint16
	ScanCode    uint16
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
	SHORT     uint16
	WORD      uint16
	DWORD     uint32
	WPARAM    uintptr
	LPARAM    uintptr
	LRESULT   uintptr
	HANDLE    uintptr
	HINSTANCE HANDLE
	HHOOK     HANDLE
	HWND      HANDLE
)
