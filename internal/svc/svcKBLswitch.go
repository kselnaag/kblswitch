package svc

import (
	"fmt"
	T "kblswitch/internal/types"
	U32 "kblswitch/internal/user32"
	"time"
	"unsafe"
)

var _ T.ISvc = (*KBLSwitch)(nil)

type KBLSwitch struct {
	log        T.ILog
	user32     *U32.User32
	swapTable  map[uint16]uint16
	swapBuff   *RingBuff[uint16]
	isBuffLock bool
}

func NewKBLSwitch(log T.ILog) *KBLSwitch {
	user32 := U32.NewUser32(log)
	return &KBLSwitch{
		log:        log,
		user32:     user32,
		swapTable:  *makeSwapTable(),
		swapBuff:   NewRingBuff[uint16](10),
		isBuffLock: false,
	}
}

func (k *KBLSwitch) Start() {
	k.setWinApiHook()
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
	k.user32.UnhookWindowsHookEx(k.user32.KeyboardHookId)
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

				switch {
				case wVirtKey == T.VK_PAUSE:
					// hwnd := k.user32.GetForegroundWindow()
					// k.user32.PostMessageA(hwnd, T.WM_INPUTLANGCHANGEREQUEST, 0, 0)

					var b T.Input
					b.InputType = T.INPUT_KEYBOARD
					b.Ki.VkCode = T.VK_BACK
					bufLen := k.swapBuff.DataLen()

					k.isBuffLock = true
					for i := 0; i < bufLen; i++ {
						swapChar, ok := k.swapTable[k.swapBuff.Read(i)]
						if ok {
							k.swapBuff.Change(i, swapChar)
						}
						k.user32.SendInput(1, &b, b)
						time.Sleep(time.Millisecond)
					}
					for i := 0; i < bufLen; i++ {
						b.Ki.VkCode = 'R' // k.swapBuff.Read(i)
						k.user32.SendInput(1, &b, b)
						time.Sleep(time.Millisecond)
					}

					// utf16char := UTF16.Encode([]rune{k.swapBuff.Read(i)})
					// fmt.Printf("%c", utf16char[0])

					k.isBuffLock = false

					fmt.Printf("%s\n", k.swapBuff.ToString())
				case wVirtKey == T.VK_ENTER:
					k.swapBuff.Clear()
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
						fmt.Printf("%s\n", k.swapBuff.ToString())
						/* if keyStateBuff[T.VK_SHIFT] > 1 {
							fmt.Println("SHIFT", keyStateBuff[T.VK_SHIFT])
						}
						if keyStateBuff[T.VK_CONTROL] > 1 {
							fmt.Println("CONTROL", keyStateBuff[T.VK_CONTROL])
						} */
					}
				}
			}
			return k.user32.CallNextHookEx(k.user32.KeyboardHookId, nCode, wparam, lparam)
		}), 0, 0)
}
