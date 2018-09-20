package hotkey

import (
	"bytes"
	"fmt"
	"syscall"
	"unsafe"
)

const (
	ModAlt = 1 << iota
	ModCtrl
	ModShift
	ModWin
)

type HotKey struct {
	Id        int // Unique id
	Modifiers int // Mask of modifiers
	KeyCode   int // Key code, e.g. 'A'
	Callback  func()
}

type MSG struct {
	HWnd    uintptr
	Message uintptr
	WParam  int16
	LParam  int64
	Time    int32
	Point   struct{ X, Y int64 }
}

// String returns a human-friendly display name of the hot key
// such as "HotKey[Id: 1, Alt+Ctrl+O]"
func (h *HotKey) String() string {
	mod := &bytes.Buffer{}
	if h.Modifiers&ModAlt != 0 {
		mod.WriteString("Alt+")
	}
	if h.Modifiers&ModCtrl != 0 {
		mod.WriteString("Ctrl+")
	}
	if h.Modifiers&ModShift != 0 {
		mod.WriteString("Shift+")
	}
	if h.Modifiers&ModWin != 0 {
		mod.WriteString("Win+")
	}
	return fmt.Sprintf("HotKey[Id: %d, %s%c]", h.Id, mod, h.KeyCode)
}

var (
	user32         = syscall.MustLoadDLL("user32")
	registerHotKey = user32.MustFindProc("RegisterHotKey")
	getMessage     = user32.MustFindProc("GetMessageW")
)

var keys = map[int16]*HotKey{}

func Start() {
	for {
		var msg = &MSG{}
		getMessage.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

		// Registered id is in the WParam field:
		if hotKey, ok := keys[msg.WParam]; ok {
			hotKey.Callback()
		}
	}
}

func Register(key *HotKey) {
	keys[int16(key.Id)] = key
	registerHotKey.Call(0, uintptr(key.Id), uintptr(key.Modifiers), uintptr(key.KeyCode))
}
