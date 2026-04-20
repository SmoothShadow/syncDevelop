//go:build windows

package clipboard

import (
	"syscall"
	"unsafe"
)

var (
	user32       = syscall.NewLazyDLL("user32.dll")
	procOpenClip = user32.NewProc("OpenClipboard")
	procCloseClip = user32.NewProc("CloseClipboard")
	procGetClipData = user32.NewProc("GetClipboardData")
	procSetClipData = user32.NewProc("SetClipboardData")
	procEmptyClip = user32.NewProc("EmptyClipboard")
	procGlobalAlloc = user32.NewProc("GlobalAlloc")
	procGlobalLock = user32.NewProc("GlobalLock")
	procGlobalUnlock = user32.NewProc("GlobalUnlock")
	procGlobalSize = user32.NewProc("GlobalSize")
	procIsClipFmtAvail = user32.NewProc("IsClipboardFormatAvailable")
)

const (
	CF_TEXT    = 1
	CF_UNICODE = 13
	CF_DIB     = 8
	GMEM_MOVEABLE = 0x0002
)

func readText() (string, error) {
	ret, _, _ := procOpenClip.Call(0)
	if ret == 0 {
		return "", nil
	}
	defer procCloseClip.Call()

	// 优先尝试Unicode文本
	if ret, _, _ := procIsClipFmtAvail.Call(CF_UNICODE); ret != 0 {
		h, _, _ := procGetClipData.Call(CF_UNICODE)
		if h == 0 {
			return "", nil
		}
		size, _, _ := procGlobalSize.Call(h)
		if size == 0 {
			return "", nil
		}
		ptr, _, _ := procGlobalLock.Call(h)
		if ptr == 0 {
			return "", nil
		}
		defer procGlobalUnlock.Call(h)
		// 读取UTF16字符串
		n := size / 2
		wchars := make([]uint16, n)
		for i := uintptr(0); i < uintptr(n); i++ {
			wchars[i] = *(*uint16)(unsafe.Pointer(ptr + i*2))
			if wchars[i] == 0 {
				wchars = wchars[:i]
				break
			}
		}
		return syscall.UTF16ToString(wchars), nil
	}

	// 退回ANSI文本
	h, _, _ := procGetClipData.Call(CF_TEXT)
	if h == 0 {
		return "", nil
	}
	ptr, _, _ := procGlobalLock.Call(h)
	if ptr == 0 {
		return "", nil
	}
	defer procGlobalUnlock.Call(h)
	size, _, _ := procGlobalSize.Call(h)
	if size == 0 {
		return "", nil
	}
	data := make([]byte, size)
	for i := uintptr(0); i < uintptr(size); i++ {
		data[i] = *(*byte)(unsafe.Pointer(ptr + i))
	}
	// 去掉末尾的\0
	for len(data) > 0 && data[len(data)-1] == 0 {
		data = data[:len(data)-1]
	}
	return string(data), nil
}

func writeText(text string) error {
	ret, _, _ := procOpenClip.Call(0)
	if ret == 0 {
		return nil
	}
	defer procCloseClip.Call()

	procEmptyClip.Call()

	// 写入Unicode文本
	wstr, err := syscall.UTF16FromString(text)
	if err != nil {
		return err
	}
	size := uintptr(len(wstr)) * 2
	h, _, _ := procGlobalAlloc.Call(GMEM_MOVEABLE, size)
	if h == 0 {
		return nil
	}
	ptr, _, _ := procGlobalLock.Call(h)
	if ptr == 0 {
		return nil
	}
	for i, c := range wstr {
		*(*uint16)(unsafe.Pointer(ptr + uintptr(i)*2)) = c
	}
	procGlobalUnlock.Call(h)
	procSetClipData.Call(CF_UNICODE, h)
	return nil
}

func readImage() ([]byte, error) {
	// Windows图片读取使用DIB格式
	ret, _, _ := procOpenClip.Call(0)
	if ret == 0 {
		return nil, nil
	}
	defer procCloseClip.Call()

	if ret, _, _ := procIsClipFmtAvail.Call(CF_DIB); ret == 0 {
		return nil, nil
	}

	h, _, _ := procGetClipData.Call(CF_DIB)
	if h == 0 {
		return nil, nil
	}
	ptr, _, _ := procGlobalLock.Call(h)
	if ptr == 0 {
		return nil, nil
	}
	defer procGlobalUnlock.Call(h)
	size, _, _ := procGlobalSize.Call(h)
	if size == 0 {
		return nil, nil
	}
	data := make([]byte, size)
	for i := uintptr(0); i < uintptr(size); i++ {
		data[i] = *(*byte)(unsafe.Pointer(ptr + i))
	}
	return data, nil
}

func writeImage(data []byte) error {
	// DIB格式写入剪贴板
	ret, _, _ := procOpenClip.Call(0)
	if ret == 0 {
		return nil
	}
	defer procCloseClip.Call()

	procEmptyClip.Call()

	size := uintptr(len(data))
	h, _, _ := procGlobalAlloc.Call(GMEM_MOVEABLE, size)
	if h == 0 {
		return nil
	}
	ptr, _, _ := procGlobalLock.Call(h)
	if ptr == 0 {
		return nil
	}
	for i := uintptr(0); i < uintptr(size); i++ {
		*(*byte)(unsafe.Pointer(ptr + uintptr(i))) = data[i]
	}
	procGlobalUnlock.Call(h)
	procSetClipData.Call(CF_DIB, h)
	return nil
}
