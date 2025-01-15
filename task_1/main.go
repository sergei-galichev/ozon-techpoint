package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

/*
 Рекомендуется использовать быстрый (буферизованный) ввод и вывод
var in *bufio.Reader
var out *bufio.Writer
in = bufio.NewReader(os.Stdin)
out = bufio.NewWriter(os.Stdout)
defer out.Flush()

var a, b int
fmt.Fscan(in, &a, &b)
fmt.Fprint(out, a + b)
*/

func main() {
	var in *bufio.Reader
	var out *bufio.Writer
	var setsNumber int

	in = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
	defer func() {
		err := out.Flush()
		if err != nil {
			panic(err)
		}
	}()

	_, err := fmt.Fscan(in, &setsNumber)
	if err != nil {
		log.Fatalf("Number of sets scan error: %s", err.Error())
	}

	for i := 0; i < setsNumber; i++ {
		var number string
		_, err = fmt.Fscan(in, &number)
		if err != nil {
			log.Fatalf("Number scan error: %s", err.Error())
		}

		maxSalary := getMaxSalaryWithoutOneDigit(number)

		fmt.Fprintln(out, maxSalary)
	}
}

func getMaxSalaryWithoutOneDigit(salary string) string {
	if len(salary) == 1 {
		return "0"
	}

	maxSalary := ""
	isRemoved := false

	for i := 0; i < len(salary); i++ {
		if i == len(salary)-1 && !isRemoved {
			maxSalary = salary[:i] + salary[i+1:]

			break
		}

		if salary[i] < salary[i+1] {
			maxSalary = salary[:i] + salary[i+1:]
			isRemoved = true

			break
		}
	}

	return maxSalary
}
