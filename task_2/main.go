package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

	// TODO: remove it after complete
	fmt.Println("Enter count of sets:")

	_, err := fmt.Fscan(in, &setsNumber)
	if err != nil {
		log.Fatalf("Count of sets scan error: %s", err.Error())
	}

	for i := 0; i < setsNumber; i++ {
		result := scanSetAndValidate(in)
		if !result {
			fmt.Fprintln(out, "no")
		} else {
			fmt.Fprintln(out, "yes")
		}
	}
}

func scanSetAndValidate(in *bufio.Reader) bool {
	numberCount, err := scanNumberCount(in)
	if err != nil {
		return false
	}

	srcArray, err := scanSourceArrayOfIntegers(in, numberCount)
	if err != nil {
		return false
	}

	validateArray, err := scanArrayToValidate(in, numberCount)
	if err != nil {
		return false
	}

	if !checkLengths(srcArray, validateArray) {
		return false
	}

	return true
}

func scanNumberCount(in *bufio.Reader) (int, error) {
	var numberCount int
	// TODO: remove it after complete
	fmt.Println("Enter count of numbers:")

	_, err := fmt.Fscan(in, &numberCount)
	if err != nil {
		// TODO: remove it after complete
		log.Printf("Count of numbers scan error: %s", err.Error())

		return -1, err
	}

	return numberCount, nil
}

func scanSourceArrayOfIntegers(in *bufio.Reader, count int) ([]int, error) {
	srcArray := make([]int, count)

	// TODO: remove it after complete
	fmt.Println("Enter source array with space delimiter:")

	for i := 0; i < count; i++ {
		_, err := fmt.Fscan(in, &srcArray[i])
		if err != nil {
			// TODO: remove it after complete
			log.Printf("Source array scan error: %s", err.Error())

			return nil, err
		}
	}

	return srcArray, nil
}

func scanArrayToValidate(in *bufio.Reader, count int) ([]int, error) {
	// TODO: remove it after complete
	fmt.Println("Enter array to validate with space delimiter:")

	validateArray := make([]int, count)

	for {
		str, e := in.ReadString('\n')
		if e != nil {
			if e == io.EOF {
				break
			} else {
				fmt.Println(e)
				return nil, e
			}
		}

		value := strings.TrimSpace(str)

		if value == "" {
			log.Printf("entered value is empty")
			continue
		}

		log.Printf("validate array: %s", value)
		//log.Printf("validate array: %d", len(strings.TrimSpace(str)))
	}

	//for i := 0; i < count; i++ {
	//	_, err := fmt.Fscan(in, &validateArray[i])
	//	if err != nil {
	//TODO: remove it after complete
	//log.Printf("Validate array scan error: %s", err.Error())
	//
	//return nil, err
	//}
	//}

	return validateArray, nil
}

func checkLengths(srcArray, validateArray []int) bool {
	return len(srcArray) == len(validateArray)
}
