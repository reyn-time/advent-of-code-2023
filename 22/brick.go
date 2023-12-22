package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
)

type Point struct {
	x, y, z int
}

type Brick struct {
	start, end  Point
	bricksAbove []int
	bricksBelow []int
}

type Buffer struct {
	height, brickIndex int
}

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("brick.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	bricks := make([]*Brick, 0)
	for scanner.Scan() {
		var x1, y1, z1, x2, y2, z2 int
		fmt.Sscanf(scanner.Text(), "%d,%d,%d~%d,%d,%d", &x1, &y1, &z1, &x2, &y2, &z2)
		brick := &Brick{Point{x1, y1, z1}, Point{x2, y2, z2}, make([]int, 0), make([]int, 0)}
		bricks = append(bricks, brick)
	}

	// Sort bricks by z values
	slices.SortFunc(bricks, func(i, j *Brick) int {
		return i.start.z - j.start.z
	})

	N := 10
	zBuffer := make([][]Buffer, N)
	for i := 0; i < N; i++ {
		zBuffer[i] = make([]Buffer, N)
	}
	for i, brick := range bricks {
		maxHeight := 1
		bricksBelow := make(map[int]bool)

		// Find bricks immediately below this brick
		for x := brick.start.x; x <= brick.end.x; x++ {
			for y := brick.start.y; y <= brick.end.y; y++ {
				if zBuffer[x][y].height == 0 {
					continue
				}
				if zBuffer[x][y].height > maxHeight {
					maxHeight = zBuffer[x][y].height
					bricksBelow = map[int]bool{zBuffer[x][y].brickIndex: true}
				} else if zBuffer[x][y].height == maxHeight {
					bricksBelow[zBuffer[x][y].brickIndex] = true
				}
			}
		}

		// Update bricksBelow for this brick
		bricksBelowArr := make([]int, 0)
		for brickIndex := range bricksBelow {
			bricksBelowArr = append(bricksBelowArr, brickIndex)
		}
		bricks[i].bricksBelow = bricksBelowArr

		// Update bricksAbove for bricks immediately below this brick
		for brickIndex := range bricksBelow {
			bricks[brickIndex].bricksAbove = append(bricks[brickIndex].bricksAbove, i)
		}

		// Update zBuffer
		zOffset := brick.start.z - maxHeight - 1
		brick.start.z -= zOffset
		brick.end.z -= zOffset
		for x := brick.start.x; x <= brick.end.x; x++ {
			for y := brick.start.y; y <= brick.end.y; y++ {
				zBuffer[x][y].height = brick.end.z
				zBuffer[x][y].brickIndex = i
			}
		}
	}

	if !*secondFlag {
		total := 0
		for _, brick := range bricks {
			canRemove := true
			for _, brickIndex := range brick.bricksAbove {
				canRemove = canRemove && len(bricks[brickIndex].bricksBelow) > 1
			}
			if canRemove {
				total++
			}
		}
		fmt.Println(total)
		return
	}

	total := 0
	for i := range bricks {
		bricksBelowGone := make([]int, len(bricks))
		count := -1
		bricksBelowGone[i] += len(bricks[i].bricksBelow)
		for j := i; j < len(bricks); j++ {
			if len(bricks[j].bricksBelow) == bricksBelowGone[j] && (len(bricks[j].bricksBelow) > 0 || j == i) {
				for _, brickIndex := range bricks[j].bricksAbove {
					bricksBelowGone[brickIndex]++
				}
				count++
			}
		}
		total += count
	}
	fmt.Println(total)
}
