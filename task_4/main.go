package main

import (
	"bufio"
	"cmp"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
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
	IsFull   bool
}

func main() {
	var in *bufio.Reader
	var out *bufio.Writer
	var setsNumber int

	//file, fErr := os.Open("1")
	//if fErr != nil {
	//	log.Fatalf("file open error: %s", fErr.Error())
	//}

	in = bufio.NewReader(os.Stdin)
	//in = bufio.NewReader(file)

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

		processedOrders := orderProcessing(orders, cars)

		printOrders(out, processedOrders)
	}
}

func readOrderArrivals(in *bufio.Reader, orderCount int) []Order {
	orders := make([]Order, orderCount)

	for i := 0; i < orderCount; i++ {
		var arrivalTime int
		order := Order{}

		_, err := fmt.Fscan(in, &arrivalTime)
		if err != nil {
			log.Fatalf("Order arrival scan error: %s", err.Error())
		}

		order.Number = i
		order.ArrivalTime = arrivalTime
		order.CarNumber = -1

		orders[i] = order
	}

	return orders
}

func readCarData(in *bufio.Reader, carCount int) []Car {
	cars := make([]Car, carCount)
	tmp := make([]int, 3)

	cycle := 0

	for i := 0; i < carCount; {
		if cycle == 0 {
			fmt.Fscanln(in)
			cycle++

			continue
		}

		car := Car{}

		_, err := fmt.Fscanln(in, &tmp[0], &tmp[1], &tmp[2])
		if err != nil {
			log.Fatalf("Car data scan error: %s", err.Error())
		}

		car.Number = i
		car.Start = tmp[0]
		car.End = tmp[1]
		car.Capacity = tmp[2]

		cars[i] = car

		i++
	}

	return cars
}

func orderProcessing(orders []Order, cars []Car) []Order {
	slices.SortFunc(orders, sortByArrival)
	slices.SortFunc(cars, sortByStartThenByIndex)

	for i := 0; i < len(orders); i++ {
		for j := 0; j < len(cars); j++ {
			if cars[j].IsFull {
				continue
			}

			if !(orders[i].ArrivalTime >= cars[j].Start && orders[i].ArrivalTime <= cars[j].End) {
				continue
			}

			orders[i].CarNumber = cars[j].Number + 1

			cars[j].Size++

			if cars[j].Size == cars[j].Capacity {
				cars[j].IsFull = true
			}

			break
		}
	}

	slices.SortFunc(orders, sortByNumber)

	return orders
}

func printOrders(out *bufio.Writer, orders []Order) {
	var result string
	for i := 0; i < len(orders); i++ {
		result += strconv.Itoa(orders[i].CarNumber)

		if i != len(orders)-1 {
			result += " "
		}
	}

	fmt.Fprint(out, result)

}

func sortByArrival(orderA, orderB Order) int {
	return cmp.Compare(orderA.ArrivalTime, orderB.ArrivalTime)
}

func sortByNumber(orderA, orderB Order) int {
	return cmp.Compare(orderA.Number, orderB.Number)
}

func sortByStartThenByIndex(carA, carB Car) int {
	byStart := cmp.Compare(carA.Start, carB.Start)

	if byStart != 0 {
		return byStart
	}

	return cmp.Compare(carA.Number, carB.Number)
}
