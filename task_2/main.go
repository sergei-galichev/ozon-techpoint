package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

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

	log.Printf("Number of sets: %d", setsNumber)

	for i := 0; i < setsNumber; i++ {
		var count int

		_, err = fmt.Fscan(in, &count)
		if err != nil {
			log.Printf("Count scan error: %s\n", err.Error())

			fmt.Fprintln(out, "no")

			continue
		}

		log.Printf("Count of numbers: %d", count)

		var arrayStr string

		_, err = fmt.Fscan(in, &arrayStr)
		if err != nil {
			log.Printf("Array scan error: %s\n", err.Error())

			fmt.Fprintln(out, "no")

			continue
		}

		var sortedArrayStr string

		_, err = fmt.Fscan(in, &sortedArrayStr)
		if err != nil {
			log.Printf("Sorted array scan error: %s\n", err.Error())

			fmt.Fprintln(out, "no")

			continue
		}

	}
}
