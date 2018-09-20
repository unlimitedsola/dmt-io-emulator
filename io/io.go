package io

import (
	"fmt"
	"github.com/tarm/serial"
	"os"
	"sola.love/dmt-io-emulator/hotkey"
)

func Start(signal chan<- bool) {
	portConfig := &serial.Config{Name: "COM8", Baud: 9600}
	port, err := serial.OpenPort(portConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to open COM8 port, IO control won't function")
		signal <- false
		return
	} else {
		signal <- true
	}
	hotkey.Register(&hotkey.HotKey{
		Id:       1,
		KeyCode:  'Z',
		Callback: func() { coin(port) },
	})
	hotkey.Register(&hotkey.HotKey{
		Id:       2,
		KeyCode:  'X',
		Callback: func() { service(port) },
	})
	hotkey.Register(&hotkey.HotKey{
		Id:       3,
		KeyCode:  'C',
		Callback: func() { test(port) },
	})
	fmt.Println("Z -> Coin")
	fmt.Println("X -> Service")
	fmt.Println("C -> Test")
	hotkey.Start()
}

func coin(port *serial.Port) {
	port.Write([]byte("\xaa\x01\xa5"))
}

func test(port *serial.Port) {
	port.Write([]byte("\xaa\x03\xa5"))
}

func service(port *serial.Port) {
	port.Write([]byte("\xaa\x18\xa5"))
}
