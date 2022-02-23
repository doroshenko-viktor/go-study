package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Go concurrency examples")
}

// returns data from channel if it received during given period of time, else will return default value
func TryReceiveWithTimeout(c <-chan int, duration time.Duration) (data int, more, ok bool) {
	select {
	case data, more = <-c:
		return data, more, true
	case <-time.After(duration):
		return 0, true, false
	}
}

// sends data from In to firts channel, which will not be blocked
func Funout(In <-chan int, OutA, OutB chan<- int) {
	for data := range In { // this will receive data from In untill it will be closed
		select {
		case OutA <- data:
		case OutB <- data:
		}
	}
}

// receives data from first non blocked input channels and sends
// it to the first not blocked output channel
func Turnout(InA, InB <-chan int, OutA, OutB chan<- int) {
	for {
		select {
		case data, more := <-InA:
			if !more {
				return
			}
			select {
			case OutA <- data:
			case OutB <- data:
			}
		case data, more := <-InB:
			if !more {
				return
			}
			select {
			case OutA <- data:
			case OutB <- data:
			}
		}
	}
}

func TurnoutWithQuit(Quit <-chan int, InA, InB <-chan int, OutA, OutB chan<- int) {
	for {
		select {
		case data := <-InA:
		case data := <-InB:
		case <-Quit:
			close(InA)
			close(InB)
			Fanout(InA, OutA, OutB)
			Fanout(InB, OutA, OutB)
			return
		}
	}
}
