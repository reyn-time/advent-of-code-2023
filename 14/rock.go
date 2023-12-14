package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("rock.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var rockMap [][]rune

	for scanner.Scan() {
		rockMap = append(rockMap, []rune(scanner.Text()))
	}

	if !*secondFlag {
		northShift(rockMap)
		fmt.Println(calculateCost(rockMap))
	} else {
		rockMapMemory := make(map[string]int)
		var history []string
		var time, prevTime int
		N := 1000000000
		for ; time < N; time++ {
			cycle(rockMap)
			s := rockMapToString(rockMap)

			if t, found := rockMapMemory[s]; found {
				prevTime = t
				break
			} else {
				rockMapMemory[s] = time
				history = append(history, s)
			}
		}

		lastRockMap := stringToRockMap(history[prevTime+(N-1-prevTime)%(time-prevTime)])
		fmt.Println(calculateCost(lastRockMap))
	}
}

func cycle(rockMap [][]rune) {
	northShift(rockMap)
	westShift(rockMap)
	southShift(rockMap)
	eastShift(rockMap)
}

func northShift(rockMap [][]rune) {
	colCount := len(rockMap[0])
	rowCount := len(rockMap)

	for j := 0; j < colCount; j++ {
		var nextRockPos int
		for i := 0; i < rowCount; i++ {
			if rockMap[i][j] == 'O' {
				rockMap[i][j], rockMap[nextRockPos][j] = rockMap[nextRockPos][j], rockMap[i][j]
				nextRockPos++
			} else if rockMap[i][j] == '#' {
				nextRockPos = i + 1
			}
		}
	}
}

func westShift(rockMap [][]rune) {
	colCount := len(rockMap[0])
	rowCount := len(rockMap)

	for i := 0; i < rowCount; i++ {
		var nextRockPos int
		for j := 0; j < colCount; j++ {
			if rockMap[i][j] == 'O' {
				rockMap[i][j], rockMap[i][nextRockPos] = rockMap[i][nextRockPos], rockMap[i][j]
				nextRockPos++
			} else if rockMap[i][j] == '#' {
				nextRockPos = j + 1
			}
		}
	}
}

func southShift(rockMap [][]rune) {
	colCount := len(rockMap[0])
	rowCount := len(rockMap)

	for j := 0; j < colCount; j++ {
		nextRockPos := rowCount - 1
		for i := rowCount - 1; i >= 0; i-- {
			if rockMap[i][j] == 'O' {
				rockMap[i][j], rockMap[nextRockPos][j] = rockMap[nextRockPos][j], rockMap[i][j]
				nextRockPos--
			} else if rockMap[i][j] == '#' {
				nextRockPos = i - 1
			}
		}
	}
}

func eastShift(rockMap [][]rune) {
	colCount := len(rockMap[0])
	rowCount := len(rockMap)

	for i := 0; i < rowCount; i++ {
		nextRockPos := colCount - 1
		for j := colCount - 1; j >= 0; j-- {
			if rockMap[i][j] == 'O' {
				rockMap[i][j], rockMap[i][nextRockPos] = rockMap[i][nextRockPos], rockMap[i][j]
				nextRockPos--
			} else if rockMap[i][j] == '#' {
				nextRockPos = j - 1
			}
		}
	}
}

func rockMapToString(rockMap [][]rune) string {
	var sb strings.Builder
	for _, row := range rockMap {
		for _, col := range row {
			sb.WriteRune(col)
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func stringToRockMap(s string) [][]rune {
	var rockMap [][]rune
	for _, line := range strings.Split(s, "\n") {
		rockMap = append(rockMap, []rune(line))
	}
	return rockMap
}

func calculateCost(rockMap [][]rune) int {
	colCount := len(rockMap[0])
	var cost int
	for i, row := range rockMap {
		for _, col := range row {
			if col == 'O' {
				cost += colCount - i
			}
		}
	}
	return cost
}
