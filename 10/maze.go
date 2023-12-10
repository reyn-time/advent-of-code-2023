package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type cell struct {
	row, col int
}

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("maze.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	maze := make([]string, 0)
	startRow, startCol := 0, 0
	i := 0
	for scanner.Scan() {
		t := scanner.Text()
		maze = append(maze, t)

		sPos := strings.IndexByte(t, 'S')
		if sPos != -1 {
			startRow, startCol = i, sPos
		}
		i++
	}

	// Create map of all cells within the main loop
	inLoop := make(map[cell]bool)
	inLoop[cell{startRow, startCol}] = true

	// Find the path connecting S
	currRow, currCol := startRow, startCol
	direction := " "
	if currRow > 0 && strings.IndexByte("7|F", maze[currRow-1][currCol]) != -1 {
		direction = "N"
		currRow -= 1
	} else if currCol < len(maze[currRow])-1 && strings.IndexByte("J-7", maze[currRow][currCol+1]) != -1 {
		direction = "E"
		currCol += 1
	} else if currRow < len(maze)-1 && strings.IndexByte("|LJ", maze[currRow+1][currCol]) != -1 {
		direction = "S"
		currRow += 1
	} else if currCol > 0 && strings.IndexByte("-LF", maze[currRow][currCol-1]) != -1 {
		direction = "W"
		currCol -= 1
	}
	initDirection := direction

	// Follow the path
	step := 1
	for maze[currRow][currCol] != 'S' {
		inLoop[cell{currRow, currCol}] = true

		switch string(maze[currRow][currCol]) + direction {
		case "|N", "LW", "JE":
			direction = "N"
			currRow -= 1
		case "-E", "LS", "FN":
			direction = "E"
			currCol += 1
		case "|S", "7E", "FW":
			direction = "S"
			currRow += 1
		case "-W", "JS", "7N":
			direction = "W"
			currCol -= 1
		default:
			fmt.Println("Invalid path")
			return
		}
		step++
	}

	// Update what the start tile should be
	startTile := " "
	switch direction + initDirection {
	case "NN", "SS":
		startTile = "|"
	case "EE", "WW":
		startTile = "-"
	case "NE", "WS":
		startTile = "F"
	case "SE", "WN":
		startTile = "L"
	case "SW", "EN":
		startTile = "J"
	case "NW", "ES":
		startTile = "7"
	default:
		fmt.Println("Invalid path")
		return
	}
	maze[startRow] = maze[startRow][:startCol] + startTile + maze[startRow][startCol+1:]

	if !*secondFlag {
		fmt.Println(step / 2)
		return
	}

	// Count all cells within the loop
	// For each cell, shoot a ray to the (lower) left and count the number of walls it intersects.
	// If the number is odd, the cell is inside the loop.
	count := 0
	for i := 0; i < len(maze); i++ {
		walls := 0
		for j := 0; j < len(maze[i]); j++ {
			if inLoop[cell{i, j}] {
				if maze[i][j] == 'F' || maze[i][j] == '7' || maze[i][j] == '|' {
					walls += 1
				}
			} else if walls%2 == 1 {
				// Not in loop and intersected odd number of walls. It is inside the loop.
				count += 1
			}
		}
	}
	fmt.Println(count)
}
