package exercises

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
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

func ThirdChannelExample() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		for {
			c1 <- "from 1"
			SleepChannelImplementation(time.Second * 2)
			// time.Sleep(time.Second * 2)
		}
	}()

	go func() {
		for {
			c2 <- "from 2"
			SleepChannelImplementation(time.Second * 3)
			// time.Sleep(time.Second * 3)
		}
	}()

	go func() {
		for {
			select {
			case msg1 := <-c1:
				fmt.Println(msg1)
			case msg2 := <-c2:
				fmt.Println(msg2)
			case <-time.After(time.Second): // Trigger after each second
				fmt.Println("timeout")
				// default:
				// 	fmt.Println("nothing ready")
			}
		}
	}()

	var input string
	fmt.Scanln(&input)
}

func FourthChannelExample() {
	type HomePageSize struct {
		URL  string
		Size int
	}

	urls := []string{
		"http://www.apple.com",
		"http://www.amazon.com",
		"http://www.google.com",
		"http://www.microsoft.com",
	}

	results := make(chan HomePageSize)
	for _, url := range urls {
		go func(url string) {
			res, err := http.Get(url)
			if err != nil {
				panic(err)
			}

			defer res.Body.Close()

			bs, err := io.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}

			results <- HomePageSize{
				URL:  url,
				Size: len(bs),
			}
		}(url)
	}

	var biggest HomePageSize

	for range urls {
		result := <-results
		if result.Size > biggest.Size {
			biggest = result
		}
	}

	fmt.Println("The biggest home page: ", biggest.URL)
}

func SleepChannelImplementation(deadline time.Duration) {
	for {
		_ = <-time.After(deadline)
		return
	}
}

func BufferedChannelTest() {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	// ch <- 3 - Overflows the buffer and blocks the execution
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
