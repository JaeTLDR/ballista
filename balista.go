package main

import (
	"fmt"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
	"golang.org/x/sys/windows"
	"sync"
	"syscall"
	"unsafe"
)

type (
	HANDLE uintptr
	HWND   HANDLE
)

type balista struct {
	hwnd  HWND
	title string
}

var (
	kernal                  = windows.NewLazyDLL("Kernel32.dll")
	procGetConsoleWindow    = kernal.NewProc("GetConsoleWindow")
	mod                     = windows.NewLazyDLL("user32.dll")
	procGetWindowText       = mod.NewProc("GetWindowTextW")
	procGetWindowTextLength = mod.NewProc("GetWindowTextLengthW")
	procSetActiveWindow     = mod.NewProc("SetActiveWindow")
	procMessageBox          = mod.NewProc("MessageBoxW")
	procSetForegroundWindow = mod.NewProc("SetForegroundWindow")
	procGetWindow           = mod.NewProc("GetWindow")

	keyMap   map[string]*hotkey.Hotkey
	numbers  = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	numCodes = []hotkey.Key{hotkey.Key1, hotkey.Key2, hotkey.Key3, hotkey.Key4, hotkey.Key5, hotkey.Key6, hotkey.Key7, hotkey.Key8, hotkey.Key9}
	balistas [10]balista
)

func init() {
	var err error
	keyMap = make(map[string]*hotkey.Hotkey)
	ctrl := []hotkey.Modifier{hotkey.ModCtrl}
	ctrlShift := []hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}

	rhk := hotkey.New(ctrl, hotkey.Key0)
	if rerr := rhk.Register(); rerr != nil {
		panic(fmt.Sprintf("could not register ghk%d", 0))
	}
	keyMap[fmt.Sprintf("ghk%d", 0)] = rhk
	for i, n := range numbers {
		hk := hotkey.New(ctrl, numCodes[i])
		err = hk.Register()
		if err != nil {
			panic(fmt.Sprintf("could not register ghk%d", n))
		}
		keyMap[fmt.Sprintf("ghk%d", n)] = hk
		hks := hotkey.New(ctrlShift, numCodes[i])
		err = hks.Register()
		if err != nil {
			panic(fmt.Sprintf("could not register ghks%d", n))
		}
		keyMap[fmt.Sprintf("ghk%d-set", n)] = hks
	}
}

func cleanup() {
	for _, v := range keyMap {
		v.Register()
	}
}

func main() { 
	mainthread.Init(fn) 
}

func fn() {
	rootHWND, _, _ := procGetConsoleWindow.Call()
	balistas[0] = balista{hwnd: HWND(rootHWND), title: "Balista List/Help"}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		listenForKeys()
		wg.Done()
	}()
	wg.Wait()
	cleanup()
}

func setActiveWindow(action string, bid int) {
	hwnd := balistas[bid].hwnd
	if bid != 0 {
		fmt.Printf("switching to #%d %s(%d)\n", bid, balistas[bid].title, hwnd)
	}
	ret, _, _ := procSetForegroundWindow.Call(
		uintptr(hwnd),
	)
	if int(ret) == 0 {
		fmt.Println("failed to switch to ", balistas[bid].title)
	}
}

func setBalistaWindow(action string, bid int) {
	fmt.Println("setBalistawinddow", bid)
	hwnd := getWindow("GetForegroundWindow")
	if hwnd != 0 {
		text := GetWindowText(HWND(hwnd))
		fmt.Printf("setting #%d to %s(%d)\n", bid, text, hwnd)
		balistas[bid] = balista{
			hwnd:  HWND(hwnd),
			title: text,
		}
	}
}

func getWindow(funcName string) uintptr {
	proc := mod.NewProc(funcName)
	hwnd, _, _ := proc.Call()
	return hwnd
}

func GetWindowTextLength(hwnd HWND) int {
	ret, _, _ := procGetWindowTextLength.Call(
		uintptr(hwnd))
	return int(ret)
}

func GetWindowText(hwnd HWND) string {
	textLen := GetWindowTextLength(hwnd) + 1
	buf := make([]uint16, textLen)
	procGetWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen))

	return syscall.UTF16ToString(buf)
}

func listenForKeys() {
	for {
		select {
		case <-keyMap["ghk0"].Keydown():
			printWindows()
			setActiveWindow("ghk0", 0)
		case <-keyMap["ghk1"].Keydown():
			setActiveWindow("ghk1", 1)
		case <-keyMap["ghk1-set"].Keydown():
			setBalistaWindow("ghk1-set", 1)
		case <-keyMap["ghk2"].Keydown():
			setActiveWindow("ghk2", 2)
		case <-keyMap["ghk2-set"].Keydown():
			setBalistaWindow("ghk2-set", 2)
		case <-keyMap["ghk3"].Keydown():
			setActiveWindow("ghk3", 3)
		case <-keyMap["ghk3-set"].Keydown():
			setBalistaWindow("ghk3-set", 3)
		case <-keyMap["ghk4"].Keydown():
			setActiveWindow("ghk4", 4)
		case <-keyMap["ghk4-set"].Keydown():
			setBalistaWindow("ghk4-set", 4)
		case <-keyMap["ghk5"].Keydown():
			setActiveWindow("ghk5", 5)
		case <-keyMap["ghk5-set"].Keydown():
			setBalistaWindow("ghk5-set", 5)
		case <-keyMap["ghk6"].Keydown():
			setActiveWindow("ghk6", 6)
		case <-keyMap["ghk6-set"].Keydown():
			setBalistaWindow("ghk6-set", 6)
		case <-keyMap["ghk7"].Keydown():
			setActiveWindow("ghk7", 7)
		case <-keyMap["ghk7-set"].Keydown():
			setBalistaWindow("ghk7-set", 7)
		case <-keyMap["ghk8"].Keydown():
			setActiveWindow("ghk8", 8)
		case <-keyMap["ghk8-set"].Keydown():
			setBalistaWindow("ghk8-set", 8)
		case <-keyMap["ghk9"].Keydown():
			setActiveWindow("ghk9", 9)
		case <-keyMap["ghk9-set"].Keydown():
			setBalistaWindow("ghk9-set", 9)
		}
	}
}

func printWindows() {
	for i, v := range balistas {
		fmt.Printf("%d: %s #{%d}\n", i, v.title, v.hwnd)
	}
	fmt.Println("Help: ctrl+# switch; ctrl+shift+# set window to key")
}
