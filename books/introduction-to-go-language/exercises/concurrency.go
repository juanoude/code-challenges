package exercises

import (
	"fmt"
	"math/rand"
	"time"
)

func FirstConcurrencyExample() {
	for i := 0; i < 10; i++ {
		go printNumbers(i)
	}

	var input string
	fmt.Scanln(&input)
}

func printNumbers(identifier int) {
	for i := 0; i < 20; i++ {
		fmt.Println(identifier, ":", i)
		randomDuration := time.Duration(rand.Intn(250))
		time.Sleep(time.Millisecond * randomDuration)
	}
}

func FirstChannelExample() {
	// Notice that the emitter is blocked until the receiver is ready again
	var c chan string = make(chan string)
	go pinger(c)
	go printer(c)

	var input string
	fmt.Scanln(&input)
}

func pinger(c chan string) {
	for i := 0; ; i++ {
		c <- "ping"
	}
}

func printer(c chan string) {
	for {
		msg := <-c
		fmt.Println(msg)
		time.Sleep(time.Second * 1) // This examplifies the blocking nature of channels
	}
}

func SecondChannelExample() {
	// Now it will alternate between ping and pong
	// It makes sense since they share the same channel
	var c chan string = make(chan string)
	go pinger(c)
	go ponger(c)
	go printer(c)

	var input string
	fmt.Scanln(&input)
}

func ponger(c chan string) {
	for i := 0; ; i++ {
		c <- "pong"
	}
}
