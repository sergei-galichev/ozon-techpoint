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

	for i := 0; i < setsNumber; i++ {
		var count string
		//var srcArray string
		//var sortedArray string

		_, err = fmt.Fscan(in, &count)
		if err != nil {
			log.Printf("Data scan error: %s", err.Error())

			fmt.Fprintln(out, "no")

			continue
		}

		str, err := in.ReadString('\n')
		if err != nil {
			fmt.Fprintln(out, "no")
			continue
		}

		log.Println(str)
		//log.Println(sortedArray)

		//sortedArray, err := scanString(in)
		//if err != nil {
		//	fmt.Fprintln(out, "no")
		//	continue
		//}
		//
		//log.Printf("Sorted array: %s", sortedArray)

		//if !checkLengths(sourceArray, sortedArray) {
		//	fmt.Fprintln(out, "no")
		//	continue
		//}

	}
}

func scanString(in *bufio.Reader) ([]any, error) {
	//var str []string

	str := make([]any, 0)

	_, err := fmt.Fscan(in, str...)
	if err != nil {
		log.Printf("String scan error: %s\n", err.Error())

		return nil, err
	}

	return str, nil
}

func checkLengths(sourceArray, sortedArray string) bool {
	return len(sourceArray) == len(sortedArray)
}
