package main

import (
	"github.com/go-vgo/robotgo"
	"time"
)

func main() {
	robotgo.MouseSleep = 100
	for i := 0; i < 9999999; i++ {
		robotgo.KeyTap("up")
		time.Sleep(time.Millisecond * 1)
	}
}
