package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
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

	validateArray, err := scanArray(in)
	if err != nil {
		return false
	}

	if len(validateArray) != numberCount {
		return false
	}

	slices.Sort(srcArray)

	if slices.Compare(srcArray, validateArray) != 0 {
		return false
	}

	return true
}

func scanNumberCount(in *bufio.Reader) (int, error) {
	var numberCount int

	_, err := fmt.Fscan(in, &numberCount)
	if err != nil {

		return -1, err
	}

	return numberCount, nil
}

func scanSourceArrayOfIntegers(in *bufio.Reader, count int) ([]int, error) {
	srcArray := make([]int, count)

	for i := 0; i < count; i++ {
		_, err := fmt.Fscan(in, &srcArray[i])
		if err != nil {
			log.Printf("scan source error: %s", err.Error())
			return nil, err
		}
	}

	return srcArray, nil
}

func scanArray(in *bufio.Reader) ([]int, error) {
	var value string

	cycle := 0

	for {
		str, e := in.ReadString('\n')
		if e != nil {
			fmt.Println(e)
			return nil, e
		}

		if cycle == 0 {
			cycle++

			continue
		}

		value = str
		break
	}

	if strings.HasPrefix(value, " ") || strings.HasSuffix(value, " ") {
		return nil, errors.New("has spaces before or after")
	}

	strArray := strings.Split(value, " ")

	validateArray, err := convertToIntSlice(strArray)
	if err != nil {
		return nil, err
	}

	return validateArray, nil
}

func convertToIntSlice(input []string) ([]int, error) {
	result := make([]int, len(input))

	for i := 0; i < len(input); i++ {
		var value string

		if i == len(input)-1 {
			value = strings.TrimSpace(input[i])
		} else {
			value = input[i]
		}

		if strings.HasPrefix(value, "0") {
			return nil, errors.New("value has leading zero")
		}

		if strings.HasPrefix(value, "-0") {
			return nil, errors.New("value has prefix '-0'")
		}

		number, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}

		result[i] = number
	}

	return result, nil
}
