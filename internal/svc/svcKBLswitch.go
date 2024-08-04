package svc

import (
	T "kblswitch/internal/types"
	"kblswitch/internal/winapi"
	"time"
	"unsafe"
)

var _ T.ISvc = (*KBLSwitch)(nil)

type KBLSwitch struct {
	log        T.ILog
	user32     *winapi.User32
	kernel32   *winapi.Kernel32
	swapTable  map[uint16]uint16
	swapBuff   *RingBuff[uint16]
	isBuffLock bool
}

func NewKBLSwitch(log T.ILog) *KBLSwitch {
	user32 := winapi.NewUser32(log)
	kernel32 := winapi.NewKernel32(log)
	return &KBLSwitch{
		log:        log,
		user32:     user32,
		kernel32:   kernel32,
		swapTable:  *makeSwapTable(),
		swapBuff:   NewRingBuff[uint16](1024),
		isBuffLock: false,
	}
}

func (k *KBLSwitch) Start() {
	k.setWinApiHook()
	k.kernel32.GetCurrentThreadId()
}

func (k *KBLSwitch) KeepAlive() {
	var msg T.MSG
	for k.user32.GetMessageA(&msg, 0, 0, 0) {
		k.user32.TranslateMessage(&msg)
		k.user32.DispatchMessage(&msg)
	}
}

func (k *KBLSwitch) Stop() {
	k.unsetWinApiHook()
}

func (k *KBLSwitch) checkKBLayout() uintptr {
	w := k.user32.GetForegroundWindow()
	p := k.user32.GetWindowThreadProcessId(w)
	r := k.user32.GetKeyboardLayout(p)
	return r
}

func (k *KBLSwitch) unsetWinApiHook() {
	k.user32.UnhookWindowsHookEx(k.user32.KBHookId)
}

func (k *KBLSwitch) setWinApiHook() {
	k.user32.SetWindowsHookExA(T.WH_KEYBOARD_LL,
		(T.HOOKPROC)(func(nCode int, wparam uintptr, lparam uintptr) uintptr {
			if nCode == 0 && wparam == T.WM_KEYDOWN {
				kbdstruct := (*T.KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
				wVirtKey := kbdstruct.VkCode
				wScanCode := kbdstruct.ScanCode
				k.user32.GetKeyState(T.VK_SHIFT)
				var keyStateBuff [256]byte
				k.user32.GetKeyboardState(&keyStateBuff)

				switch { // ENTER - dropBuff, SHIFT+ESC - quit, PAUSE - textSwitch, CTRL+PAUSE - textSwitch from OS buffer(Ctrl+C)
				case wVirtKey == T.VK_PAUSE:
					hwnd := k.user32.GetForegroundWindow()
					k.user32.PostMessageA(hwnd, T.WM_INPUTLANGCHANGEREQUEST, 0, 0)

					var b T.Input
					b.InputType = T.INPUT_KEYBOARD
					b.Ki.VkCode = T.VK_BACK
					b.Ki.ScanCode = 0
					b.Ki.Flags = 0

					k.isBuffLock = true
					bufLen := k.swapBuff.DataLen()
					for i := 0; i < bufLen; i++ {
						swapChar, ok := k.swapTable[k.swapBuff.Read(i)]
						if ok {
							k.swapBuff.Change(i, swapChar)
						}
						time.Sleep(10 * time.Millisecond)
						k.user32.SendInput(1, &b, b)
					}
					time.Sleep(100 * time.Millisecond)
					b.Ki.Flags = T.KEYEVENTF_UNICODE
					b.Ki.VkCode = 0
					for i := 0; i < bufLen; i++ {
						b.Ki.ScanCode = k.swapBuff.Read(i)
						time.Sleep(10 * time.Millisecond)
						k.user32.SendInput(1, &b, b)
					}
					k.isBuffLock = false
					// fmt.Printf("%s\n", k.swapBuff.ToString())
				case wVirtKey == T.VK_ENTER:
					k.swapBuff.Clear()
				case wVirtKey == T.VK_ESCAPE:
					if keyStateBuff[T.VK_SHIFT] > 1 {
						k.user32.PostThreadMessageA(k.kernel32.CurrThreadId, T.WM_QUIT, 0, 0)
					}
				case wVirtKey == T.VK_BACK:
					if !k.isBuffLock {
						k.swapBuff.Back()
					}
				default:
					if (wVirtKey != T.VK_LSHIFT) && (wVirtKey != T.VK_RSHIFT) && (!k.isBuffLock) {
						const outSize = 1
						var outBuff [outSize]uint16
						k.user32.ToUnicodeEx(wVirtKey, wScanCode, &keyStateBuff, &outBuff, outSize, 0, k.checkKBLayout())
						if outBuff[0] != 0 {
							k.swapBuff.Set(outBuff[0])
						}
						// fmt.Printf("%s\n", k.swapBuff.ToString())
						/* if keyStateBuff[T.VK_SHIFT] > 1 {
							fmt.Println("SHIFT", keyStateBuff[T.VK_SHIFT])
						}
						if keyStateBuff[T.VK_CONTROL] > 1 {
							fmt.Println("CONTROL", keyStateBuff[T.VK_CONTROL])
						} */
					}
				}
			}
			return k.user32.CallNextHookEx(k.user32.KBHookId, nCode, wparam, lparam)
		}), 0, 0)
}
