package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Basic use of a channel

// func main() {
// 	c := make(chan string)
// 	// var c chan int
// 	// c = make(chan int)

// 	go boring("boring!", c)
// 	for i := 0; i < 5; i++ {
// 		fmt.Printf("You say: %q\n", <-c)
// 	}
// 	fmt.Println("You're boring; I'm leaving.")
// }

// func boring(msg string, c chan string) {
// 	for i := 0; ; i++ {
// 		c <- fmt.Sprintf("%s %d", msg, i)
// 		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
// 	}
// }

// Generator: function that returns a channel

// func main() {
// 	c := boring("boring!")
// 	for i := 0; i < 5; i++ {
// 		fmt.Printf("You say: %q\n", <-c)
// 	}
// 	fmt.Println("You're boring; I'm leaving.")
// }

// func boring(msg string) <-chan string { // returns receive-only channel of strings
// 	c := make(chan string)
// 	go func() {
// 		for i := 0; ; i++ {
// 			c <- fmt.Sprintf("%s %d", msg, i)
// 			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
// 		}
// 	}()
// 	return c // return the channel to the caller
// }

// Channels as a handle on a service

// func main() {
// 	joe := boring("Joe")
// 	ann := boring("Ann")
// 	for i := 0; i < 5; i++ {
// 		fmt.Println(<-joe)
// 		fmt.Println(<-ann) // ann is blocked until main receives joe's value
// 	}
// 	fmt.Println("You're both boring; I'm leaving.")
// }

// Decouple the execution between joe and ann by doing a fan in function

// func fanIn(input1, input2 <-chan string) <-chan string {
// 	c := make(chan string)
// 	go func() {
// 		for {
// 			c <- <-input1
// 		}
// 	}()
// 	go func() {
// 		for {
// 			c <- <-input2
// 		}
// 	}()
// 	return c
// }

// func main() {
// 	c := fanIn(boring("Joe"), boring("Ann"))
// 	for i := 0; i < 10; i++ {
// 		fmt.Println(<-c)
// 	}
// 	fmt.Println("You're both boring; I'm leaving.")
// }

// Using select

func fanIn(input1, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s
			}
		}
	}()
	return c
}

// Using select to timeout

// func main() {
// 	c := boring("Joe")
// 	for {
// 		select {
// 		case s := <-c:
// 			fmt.Println(s)
// 		case <-time.After(1 * time.Second): // timing out each message
// 			fmt.Println("You talk too much.")
// 			return
// 		}
// 	}
// }

// Timing out the loop

// func main() {
// 	c := boring("Joe")
// 	timeout := time.After(5 * time.Second)

// 	for {
// 		select {
// 		case s := <-c:
// 			fmt.Println(s)
// 		case <-timeout:
// 			fmt.Println("You talk too much.")
// 			return
// 		}
// 	}
// }

// Quit channel

// func boring(msg string, quit chan bool) <-chan string { // returns receive-only channel of strings
// 	c := make(chan string)
// 	go func() {
// 		for i := 0; ; i++ {
// 			select {
// 			case c <- fmt.Sprintf("%s %d", msg, i):
// 				// do nothing
// 			case <-quit:
// 				return
// 			}
// 		}
// 	}()
// 	return c // return the channel to the caller
// }

// func main() {
// 	quit := make(chan bool)
// 	c := boring("Joe", quit)
// 	for i := rand.Intn(10); i >= 0; i-- {
// 		fmt.Println(<-c)
// 	}
// 	quit <- true
// }

// Notify main when quitting

// func cleanup() {
// 	fmt.Printf("Cleaned!\n")
// }

// func boring(msg string, quit chan string) <-chan string {
// 	c := make(chan string)
// 	go func() {
// 		for i := 0; ; i++ {
// 			select {
// 			case c <- fmt.Sprintf("%s %d", msg, i):
// 				// do nothing
// 			case <-quit:
// 				cleanup()
// 				quit <- "See you!"
// 				return
// 			}
// 		}
// 	}()
// 	return c // return the channel to the caller
// }

// func main() {
// 	quit := make(chan string)
// 	c := boring("Joe", quit)
// 	for i := rand.Intn(10); i >= 0; i-- {
// 		fmt.Println(<-c)
// 	}
// 	quit <- "Bye!"
// 	fmt.Printf("Joe says: %q\n", <-quit)
// }

// Naive Google search 1.0 with blocking operations

var (
	Web    = fakeSearch("web")
	Image  = fakeSearch("image")
	Video  = fakeSearch("video")
	Web1   = fakeSearch("web")
	Web2   = fakeSearch("web")
	Image1 = fakeSearch("image")
	Image2 = fakeSearch("image")
	Video1 = fakeSearch("video")
	Video2 = fakeSearch("video")
)

type Search func(query string) string

func fakeSearch(kind string) Search {
	return func(query string) string {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return fmt.Sprintf("%s result for %q\n", kind, query)
	}
}

// func Google(query string) (results []string) {
// 	results = append(results, Web(query))
// 	results = append(results, Image(query))
// 	results = append(results, Video(query))
// 	return
// }

// func main() {
// 	rand.Seed(time.Now().UnixNano())
// 	start := time.Now()
// 	results := Google("golang")
// 	elapsed := time.Since(start)
// 	fmt.Println(results)
// 	fmt.Println(elapsed)
// }

// Doing Google search 2.0 concurrently (only waiting for the slowest goroutine)

// func Google(query string) (results []string) {
// 	c := make(chan string)

// 	go func() {
// 		c <- Web(query)
// 	}()
// 	go func() {
// 		c <- Image(query)
// 	}()
// 	go func() {
// 		c <- Video(query)
// 	}()

// 	for i := 0; i < 3; i++ {
// 		result := <-c
// 		results = append(results, result)
// 	}
// 	return
// }

// Google search 2.1 with timeouts

// func Google(query string) (results []string) {
// 	c := make(chan string)

// 	go func() {
// 		c <- Web(query)
// 	}()
// 	go func() {
// 		c <- Image(query)
// 	}()
// 	go func() {
// 		c <- Video(query)
// 	}()

// 	timeout := time.After(80 * time.Millisecond)
// 	for i := 0; i < 3; i++ {
// 		select {
// 		case result := <-c:
// 			results = append(results, result)
// 		case <-timeout:
// 			fmt.Println("timed out")
// 			return
// 		}
// 	}
// 	return
// }

// Google search 3.0 reduce tail latency using replicated search servers

func First(query string, replicas ...Search) string {
	c := make(chan string)
	searchReplica := func(i int) {
		c <- replicas[i](query)
	}
	for i := range replicas {
		go searchReplica(i)
	}
	return <-c
}

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}

func Google(query string) (results []string) {
	c := make(chan string)

	go func() {
		c <- First(query, Web1, Web2)
	}()
	go func() {
		c <- First(query, Image1, Image2)
	}()
	go func() {
		c <- First(query, Video1, Video2)
	}()

	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	return
}
