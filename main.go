package main

import (
	"fmt"
	"sola.love/dmt-io-emulator/card"
	"sola.love/dmt-io-emulator/io"
)

func main() {
	signal := make(chan bool)
	go card.Start(signal)
	<-signal
	go io.Start(signal)
	<-signal
	fmt.Printf("dmt-io-emulator is ready!\n")
	<-signal
}
