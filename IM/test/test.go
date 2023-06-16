package main

import "fmt"

func main() {
	t := 0
	ch := make(chan int, 1)
	var ch1 chan<- int
	ch1 = ch
	ch1 <- 22
	t = <-ch
	fmt.Println(t)
	close(ch)
	fmt.Println(<-ch)
}
