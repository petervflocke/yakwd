package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

const waitforkey = "/usr/bin/waitforkey"

// Kbd defines pressed or release key and it's state 0 release, 1 pressed
type Kbd struct {
	key   int
	state int
}

// KeyboardWorker waits for key and add it to channel
func KeyboardWorker() {
	var k Kbd
	var err0 error
	var err1 error
	for {
		cmd := exec.Command(waitforkey)
		cmdOut, err := cmd.Output()
		if err != nil {
			log.Fatalln(waitforkey, "ended with ", err)
		}
		//fmt.Println("Keyboard :", string(cmdOut))
		// waitfor key returns string og two values key code (104, 105, ..) and state code (1 or 0) plus "\n"
		// convert it to Kbd structure and send to keyboard channel
		tmps := strings.Split(strings.TrimSuffix(string(cmdOut), "\n"), " ")
		k.key, err0 = strconv.Atoi(tmps[0])
		k.state, err1 = strconv.Atoi(tmps[1])
		// fmt.Println("errors", err, err0, err1)
		if err == nil && err0 == nil && err1 == nil {
			//fmt.Println("Added: Key:", k.key, " Pressed: ", k.state)
			keyboard <- k
		}
	}
}

// MenuWorker consumes key from channel and updates program states
// at this moment only exits
func MenuWorker() {
	// keyboard scan codes:
	// 158 back
	// 29 keyboard
	// 105 left
	// 106 right
	// 103 up
	// 108 down
	// 194 enter / middle
	// 139 menu
	// 109 right page back
	// 191 right page next
	// 193 left page back
	// 104 left page next

	for {
		// non blocking check if channel not empty
		select {
		case k, ok := <-keyboard:
			if ok { // not empty, then process the key
				// fmt.Println("Got Key:", k.key, " Pressed: ", k.state)
				switch k.key {
				case 194: // enter / middle pressed or released
					// fmt.Println("Pressed", k.key, " exiting")
					wg.Done()
				default: // .... next function to be built in, next city, detailed forecast or whatsoever
					fmt.Println("Got Key:", k.key, " Pressed: ", k.state)
				}
			} else {
				// that should not happened, keyboard worker never ends
				log.Fatalln("Channel closed!")
			}
		default:
		}
	}
}
