// Стадия фильтрации отрицательных чисел (не пропускать отрицательные числа).
// Стадия фильтрации чисел, не кратных 3 (не пропускать такие числа), исключая также и 0.
// В этой стадии предусмотреть опустошение буфера (и соответственно, передачу этих данных, если они есть, дальше) с определённым интервалом во времени.
// Значения размера буфера и этого интервала времени сделать настраиваемыми.

// Написать источник данных для конвейера. Непосредственным источником данных должна быть консоль.

// Также написать код потребителя данных конвейера. Данные от конвейера можно направить снова в консоль построчно, сопроводив
// их каким-нибудь поясняющим текстом, например: «Получены данные …».

// При написании источника данных подумайте о фильтрации нечисловых данных, которые можно ввести через консоль. Как и где их фильтровать, решайте сами.

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"src/HW2021/filtering"
	"strconv"
	"time"
)

func main() {

	numbers := make(chan int)
	go dataSource(numbers)

	filterNegative := make(chan int)
	go filtering.FilterNegative(numbers, filterNegative)

	filterNumberNotMultipleThree := make(chan int)
	go filtering.FilterNumberNotMultipleThree(filterNegative, filterNumberNotMultipleThree)

	buff := make(chan int)
	go filtering.Buffering(filterNumberNotMultipleThree, buff)

	go consumer(buff)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	select {
	case sig := <-c:
		fmt.Printf("Got %s signal. Aborting ... \n", sig)
		os.Exit(0)
	}

}

func dataSource(numb chan int) {
	menu := fmt.Sprintf("Menu:\n \"Menu\" - open the menu\n \"buffer\" - enter the buffer value\n \"timer\" - enter the timer value\n \"exit\" - exiting the program\n")
	fmt.Print(menu)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		switch scanner.Text() {
		case "menu":
			fmt.Print(menu)
		case "exit":
			fmt.Println("exit from the program")
			os.Exit(0)
		case "buffer":
			fmt.Println("enter the buffer value")
			scanner.Scan()
			buff, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("only need to enter an integer B")
				continue
			}
			filtering.BufferSize = buff
		case "timer":
			fmt.Println("enter a number after how many seconds to clear the buffer")
			scanner.Scan()
			sec, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("only need to enter an integer T")
				continue
			}
			filtering.TimeBufferClear = time.Duration(sec) * time.Second
		default:
			num, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("only need to enter integers")
				continue
			}
			numb <- num
		}

	}

}

func consumer(numbIn chan int) {
	for numb := range numbIn {
		fmt.Printf("Получены данные: %d\n", numb)
	}
}
