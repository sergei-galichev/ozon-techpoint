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

	file, fErr := os.Open("3")
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

	start := time.Now()

	for i := 0; i < setsNumber; i++ {
		var orderCount int
		var carCount int

		log.Printf("----------------begin set #%d----------------", i)

		_, err = fmt.Fscan(in, &orderCount)
		if err != nil {
			log.Fatalf("Count of orders scan error: %s", err.Error())
		}

		startSet := time.Now()

		orders := readOrderArrivals(in, orderCount)

		log.Printf("'readOrderArrivals' time execution: %s", time.Since(startSet))

		_, err = fmt.Fscan(in, &carCount)
		if err != nil {
			log.Fatalf("Count of cars scan error: %s", err.Error())
		}

		startSet = time.Now()

		cars := readCarData(in, carCount)

		log.Printf("'readCarData' time execution: %s", time.Since(startSet))

		startSet = time.Now()

		processedOrders := orderProcessing(orders, cars)

		log.Printf("'orderProcessing' time execution: %s", time.Since(startSet))

		startSet = time.Now()

		printOrders(out, processedOrders)

		log.Printf("'printOrders' time execution: %s", time.Since(startSet))

		log.Printf("----------------end set #%d----------------", i)

	}

	log.Printf("all time execution: %s", time.Since(start))
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
	slices.SortFunc(cars, sortByStartThenByNumber)

	//log.Printf("is there exist car with zero capacity: %v", checkCarZeroCapacity(cars))

	for i := 0; i < len(orders); i++ {
		middle := len(cars) / 2

		carIdxChan1 := make(chan int)
		carIdxChan2 := make(chan int)
		done := make(chan struct{})

		go func(middleIdx int) {
			if len(cars) != 0 {
				if middleIdx == 0 {
					if isInLoadingTime(cars[middleIdx], orders[i].ArrivalTime) {
						carIdxChan1 <- middleIdx
						done <- struct{}{}
						close(carIdxChan1)

						return
					}
				}

				for j := 0; j < middleIdx; j++ {
					if isInLoadingTime(cars[j], orders[i].ArrivalTime) {
						carIdxChan1 <- j
						done <- struct{}{}
						close(carIdxChan1)

						return
					}
				}
			}

			carIdxChan1 <- -1
			close(carIdxChan1)
		}(middle)

		go func(middleIdx int) {
			if len(cars) != 0 {
				j := middleIdx

				for {
					select {
					case <-done:
						carIdxChan2 <- -1
						close(done)
						close(carIdxChan2)

						return

					default:
						if j == 0 || j == len(cars) {
							//	if isInLoadingTime(cars[j], orders[i].ArrivalTime) {
							//		carIdxChan2 <- j
							//		close(carIdxChan2)
							//
							//		return
							//	}
							carIdxChan2 <- -1
							close(carIdxChan2)

							return
						}

						//if len(cars)-j == 0 {
						//	carIdxChan2 <- -1
						//	break
						//}
						if isInLoadingTime(cars[j], orders[i].ArrivalTime) {
							carIdxChan2 <- j
							close(carIdxChan2)

							return
						}

						j++
					}
				}

			}

			carIdxChan2 <- -1
			close(carIdxChan2)
		}(middle)

		carIdx1 := <-carIdxChan1
		carIdx2 := <-carIdxChan2
		//close(done)

		idx := -1

		if carIdx1 == -1 && carIdx2 == -1 {
			orders[i].CarNumber = idx

			continue
		} else if carIdx1 == -1 && carIdx2 >= 0 {
			idx = carIdx2
		} else if carIdx2 == -1 && carIdx1 >= 0 {
			idx = carIdx1
		} else if carIdx1 < carIdx2 || carIdx1 == carIdx2 {
			idx = carIdx1
		} else {
			idx = carIdx2
		}

		orders[i].CarNumber = cars[idx].Number + 1
		cars[idx].Size++

		if isCarFull(cars[idx]) {
			cars = append(cars[:idx], cars[idx+1:]...)
		}

		//for j := 0; j < len(cars); j++ {
		//	if orders[i].ArrivalTime >= cars[j].Start && orders[i].ArrivalTime <= cars[j].End {
		//		orders[i].CarNumber = cars[j].Number + 1
		//		cars[j].Size++
		//
		//		if isCarFull(cars[j]) {
		//			cars = append(cars[:j], cars[j+1:]...)
		//		}
		//		break
		//	}
		//}
	}

	slices.SortFunc(orders, sortByNumber)

	return orders
}

func printOrders(out *bufio.Writer, orders []Order) {

	//result := strings.Builder{}
	//
	//for i := 0; i < len(orders); i++ {
	//	result.WriteString(strconv.Itoa(orders[i].CarNumber))
	//
	//	if i != len(orders)-1 {
	//		result.WriteString(" ")
	//	}
	//}

	//fmt.Fprintln(out, result.String())

	result := make([]byte, 0, len(orders)*10)

	for i := 0; i < len(orders); i++ {
		result = strconv.AppendInt(result, int64(orders[i].CarNumber), 10)

		if i != len(orders)-1 {
			result = append(result, ' ')
		}
	}

	result = append(result, '\n')

	//out.Write(result)

}

func sortByArrival(orderA, orderB Order) int {
	return cmp.Compare(orderA.ArrivalTime, orderB.ArrivalTime)
}

func sortByNumber(orderA, orderB Order) int {
	return cmp.Compare(orderA.Number, orderB.Number)
}

func sortByStartThenByNumber(carA, carB Car) int {
	if carA.Start == carB.Start {
		return cmp.Compare(carA.Number, carB.Number)
	}

	return cmp.Compare(carA.Start, carB.Start)
}

func findMinAvailableCarNumber(cars []Car, prevIdx, minStart int) int {
	if len(cars) == 0 {
		return -1
	}

	idxMiddle := len(cars) / 2

	if idxMiddle == 0 {
		if isInLoadingTime(cars[idxMiddle], minStart) {
			cars[idxMiddle].Size++

			return cars[idxMiddle].Number
		}

		return cars[prevIdx].Number
	}

	tmpIdx := prevIdx

	if isInLoadingTime(cars[idxMiddle], minStart) && prevIdx != -1 && idxMiddle <= prevIdx {
		tmpIdx = idxMiddle
	}

	if idxMiddle-1 > 0 && cars[idxMiddle-1].Start <= minStart {
		tmpIdx = findMinAvailableCarNumber(cars[:idxMiddle], tmpIdx, minStart)
	}

	if idxMiddle+1 <= len(cars)-1 && cars[idxMiddle+1].Start <= minStart {
		tmpIdx = findMinAvailableCarNumber(cars[idxMiddle+1:], tmpIdx, minStart)
	}

	if tmpIdx == -1 {
		return -1
	}

	return cars[tmpIdx].Number
}

func isCarFull(car Car) bool {
	return car.Size >= car.Capacity
}

func isInLoadingTime(car Car, start int) bool {
	return car.Start <= start && car.End >= start
}

func checkCarZeroCapacity(cars []Car) bool {
	for i := 0; i < len(cars); i++ {
		if cars[i].Capacity <= 0 {
			return true
		}
	}

	return false
}
