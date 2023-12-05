package main

import (
	"fmt"
	"sync"
)

/*
You have a data processing system that reads data from multiple sources concurrently, processes it, and then aggregates the results. Each data source provides a stream of integers. Your task is to design a solution using Go that concurrently processes data from these sources, sums the integers, and aggregates the results into a final total.
*/

type DataSource struct {
	Name    string
	Channel chan int
}

func processDataSource(dataSource DataSource, resultChannel chan int) {
	for input := range dataSource.Channel {
		resultChannel <- input
	}
	close(resultChannel)
}

func aggregateResults(resultChannels []chan int) int {
	var wg sync.WaitGroup
	sum := 0
	sumChannel := make(chan int)
	for _, resultChannel := range resultChannels {
		wg.Add(1)
		go func(resultChannel chan int) {
			defer wg.Done()
			for i := range resultChannel {
				sumChannel <- i
			}
		}(resultChannel)
	}

	go func() {
		wg.Wait()
		close(sumChannel)
	}()

	for i := range sumChannel {
		sum += i
	}

	return sum
}

func main() {
	dataSources := []DataSource{
		{
			Name:    "Source 1",
			Channel: make(chan int),
		},
		{
			Name:    "Source 2",
			Channel: make(chan int),
		},
		{
			Name:    "Source 3",
			Channel: make(chan int),
		},
	}

	resultChannels := make([]chan int, len(dataSources))

	for i, dataSource := range dataSources {
		resultChannels[i] = make(chan int)
		go processDataSource(dataSource, resultChannels[i])
	}

	go func() {
		for _, datasrouce := range dataSources {
			for j := 0; j < 10; j++ {
				datasrouce.Channel <- j
			}
			close(datasrouce.Channel)
		}
	}()
	totalSum := aggregateResults(resultChannels)
	fmt.Printf("Final sum is: %d\n", totalSum)
}
