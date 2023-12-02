package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("trebuchet.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	total := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		total += getCalibrationTwo(line)
	}

	fmt.Println(total)
}

// func getCalibrationOne(s string) int {
// 	first := int(s[strings.IndexAny(s, "0123456789")] - '0')
// 	last := int(s[strings.LastIndexAny(s, "0123456789")] - '0')
// 	return first*10 + last
// }

func getCalibrationTwo(s string) int {
	numbers := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	firstIndex := strings.IndexAny(s, "0123456789")
	first := int(s[firstIndex] - '0')
	lastIndex := strings.LastIndexAny(s, "0123456789")
	last := int(s[lastIndex] - '0')

	for i, num := range numbers {
		f := strings.Index(s, num)
		if f != -1 && f < firstIndex {
			firstIndex = f
			first = i
		}
		l := strings.LastIndex(s, num)
		if l != -1 && l > lastIndex {
			lastIndex = l
			last = i
		}
	}

	return first*10 + last
}
