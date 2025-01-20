package main

import (
	"bufio"
	"cmp"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"time"
)

type Order struct {
	Number      int
	ArrivalTime int
	CarNumber   int
}

type Car struct {
	Number   int
	Start    int
	End      int
	Size     int
	Capacity int
}

func main() {
	var in *bufio.Reader
	var out *bufio.Writer
	var setsNumber int

	file, fErr := os.Open("15")
	if fErr != nil {
		log.Fatalf("file open error: %s", fErr.Error())
	}

	//in = bufio.NewReader(os.Stdin)
	in = bufio.NewReader(file)

	out = bufio.NewWriter(os.Stdout)
	defer func() {
		err := out.Flush()
		if err != nil {
			panic(err)
		}
	}()

	_, err := fmt.Fscan(in, &setsNumber)
	if err != nil {
		log.Fatalf("Count of sets scan error: %s", err.Error())
	}

	for i := 0; i < setsNumber; i++ {
		var orderCount int
		var carCount int

		_, err := fmt.Fscan(in, &orderCount)
		if err != nil {
			log.Fatalf("Count of orders scan error: %s", err.Error())
		}

		orders := readOrderArrivals(in, orderCount)

		_, err = fmt.Fscan(in, &carCount)
		if err != nil {
			log.Fatalf("Count of cars scan error: %s", err.Error())
		}

		cars := readCarData(in, carCount)

		start := time.Now()

		processedOrders := orderProcessing(orders, cars)

		printOrders(out, processedOrders)

		log.Println(time.Since(start))
	}
}

func readOrderArrivals(in *bufio.Reader, orderCount int) []Order {
	orders := make([]Order, orderCount)

	for i := 0; i < orderCount; i++ {
		_, err := fmt.Fscan(in, &orders[i].ArrivalTime)
		if err != nil {
			log.Fatalf("Order arrival scan error: %s", err.Error())
		}

		orders[i].Number = i
		orders[i].CarNumber = -1
	}

	return orders
}

func readCarData(in *bufio.Reader, carCount int) []Car {
	cars := make([]Car, carCount)
	cycle := 0

	for i := 0; i < carCount; {
		if cycle == 0 {
			fmt.Fscanln(in)
			cycle++

			continue
		}

		_, err := fmt.Fscanln(in, &cars[i].Start, &cars[i].End, &cars[i].Capacity)
		if err != nil {
			log.Fatalf("Car data scan error: %s", err.Error())
		}

		cars[i].Number = i

		i++
	}

	return cars
}

func orderProcessing(orders []Order, cars []Car) []Order {
	slices.SortFunc(orders, sortByArrival)
	slices.SortFunc(cars, sortByStartThenByIndex)

	for i := 0; i < len(orders); i++ {
		for j := 0; j < len(cars); j++ {
			if cars[j].Size >= cars[j].Capacity {
				continue
			}

			if orders[i].ArrivalTime >= cars[j].Start && orders[i].ArrivalTime <= cars[j].End {
				orders[i].CarNumber = cars[j].Number + 1
				cars[j].Size++
				break
			}
		}
	}

	slices.SortFunc(orders, sortByNumber)

	return orders
}

func printOrders(out *bufio.Writer, orders []Order) {

	result := make([]byte, 0, len(orders)*10)

	for i := 0; i < len(orders); i++ {
		result = strconv.AppendInt(result, int64(orders[i].CarNumber), 10)

		if i != len(orders)-1 {
			result = append(result, ' ')
		}
	}

	result = append(result, '\n')

	out.Write(result)

}

func sortByArrival(orderA, orderB Order) int {
	return cmp.Compare(orderA.ArrivalTime, orderB.ArrivalTime)
}

func sortByNumber(orderA, orderB Order) int {
	return cmp.Compare(orderA.Number, orderB.Number)
}

func sortByStartThenByIndex(carA, carB Car) int {
	if carA.Start == carB.Start {
		return cmp.Compare(carA.Number, carB.Number)
	}

	return cmp.Compare(carA.Start, carB.Start)
}
