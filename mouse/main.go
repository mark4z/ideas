package main

import (
	"github.com/go-vgo/robotgo"
	"time"
)

func main() {
	//4788230845369078
	robotgo.MouseSleep = 100
	for i := 0; i < 9999999; i++ {
		robotgo.KeyTap("down")
		time.Sleep(time.Millisecond * 1)
	}
}
