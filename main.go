package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tarm/serial"
	"io/ioutil"
	"log"
	"strings"
)

type Config struct {
	CardUID string `json:"card_uid"`
	CardID  string `json:"card_id"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readConfigFromFile() Config {
	config := Config{}
	data, err := ioutil.ReadFile("card.json")
	check(err)
	err = json.Unmarshal(data, &config)
	check(err)
	return config
}

func hexToString(hexStr string) string {
	bytes, err := hex.DecodeString(hexStr)
	check(err)
	return string(bytes)
}

func makeResponseTable(config Config) map[string]string {

	cardUID := hexToString(config.CardUID)
	cardID := strings.Replace(config.CardID, " ", "", -1)
	cardIDPart1 := cardID[:16]
	cardIDPart2 := cardID[16:20] + strings.Repeat("\x00", 12)

	return map[string]string{
		// is there any card? -> NOPE!
		"\x02\x00\x02\x31\x30\x03\x02": "\x02\x00\x03\x31\x30\x4e\x03",
		// eject card -> OK!
		"\x02\x00\x02\x32\x30\x03\x01": "\x02\x00\x03\x32\x30\x59\x03",
		// find card -> We got one!
		"\x02\x00\x02\x35\x30\x03\x06": "\x02\x00\x03\x35\x30\x59\x03",
		// get UID -> return our UID!
		"\x02\x00\x02\x35\x31\x03\x07": "\x02\x00\x07\x35\x31\x59" + cardUID + "\x03",
		// auth key A at sector 0 (key = 37 21 53 6a 72 40) -> You got it!
		"\x02\x00\x09\x35\x32\x00\x37\x21\x53\x6a\x72\x40\x03\x12": "\x02\x00\x04\x35\x32\x00\x59\x03",
		// read sector 0 block 1 -> Our card ID
		"\x02\x00\x04\x35\x33\x00\x01\x03\x02": "\x02\x00\x15\x35\x33\x00\x01\x59" + cardIDPart1 + "\x03",
		"\x02\x00\x04\x35\x33\x00\x02\x03\x01": "\x02\x00\x15\x35\x33\x00\x02\x59" + cardIDPart2 + "\x03",
	}

}

func calcBCC(data string) string {
	bcc := 0
	for e := range data {
		bcc ^= e
	}
	return string(bcc)
}

func main() {
	config := readConfigFromFile()
	table := makeResponseTable(config)
	fmt.Printf("dmt-io-emulator is ready!\n")
	fmt.Printf("Card UID: %s\n", config.CardUID)
	fmt.Printf("Card ID: %s\n", config.CardID)
	portConfig := &serial.Config{Name: "COM9", Baud: 9600}
	port, err := serial.OpenPort(portConfig)
	if err != nil {
		log.Fatal(err)
	}
	cmd := ""
	for {
		code := simpleRead(port)
		cmd += code
		if code == "\x03" {
			bcc := simpleRead(port)
			cmd += bcc
			fmt.Printf("command => % x\n", cmd)
			response, prs := table[cmd]
			if prs {
				_, err := port.Write([]byte("\x06"))
				if err != nil {
					log.Fatal(err)
				}
				confirm := simpleRead(port)
				if confirm == "\x05" {
					port.Write([]byte(response + calcBCC(response)))
					fmt.Printf("response => % x\n", response)
				}
			} else {
				_, err := port.Write([]byte("\x15"))
				if err != nil {
					log.Fatal(err)
				}
			}
			cmd = ""
		}
	}
}

func simpleRead(port *serial.Port) string {
	buf := make([]byte, 1)
	n, err := port.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	return string(buf[:n])
}
