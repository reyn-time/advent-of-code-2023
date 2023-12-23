package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

type Point struct {
	x, y int
}

type QueueInfo struct {
	point     Point
	direction int
}

type Edge struct {
	dest Point
	cost int
}

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("hike.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	maze := make([][]rune, 0)
	for scanner.Scan() {
		maze = append(maze, []rune(scanner.Text()))
	}

	offset := [][2]int{{0, 1}, {-1, 0}, {1, 0}, {0, -1}} // right, up, down, left
	visited := make(map[Point]bool)
	queue := []QueueInfo{{Point{0, 1}, 2}}
	graph := make(map[Point][]Edge)
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		start := Point{curr.point.x, curr.point.y}
		o := map[rune][2]int{'>': {0, -1}, '^': {1, 0}, 'v': {-1, 0}, '<': {0, 1}}[maze[start.x][start.y]]
		start.x += o[0]
		start.y += o[1]

		steps := 0
		for (maze[curr.point.x][curr.point.y] == '.' || steps == 0) && !(curr.point.x == len(maze)-1 && curr.point.y == len(maze[0])-2) {
			steps++
			for i, o := range offset {
				x, y := curr.point.x+o[0], curr.point.y+o[1]
				if x < 0 || x >= len(maze) || y < 0 || y >= len(maze[0]) || maze[x][y] == '#' || i+curr.direction == 3 {
					// Don't go out of bounds, don't go into walls, don't go back the way we came
					continue
				}
				curr.point.x = x
				curr.point.y = y
				curr.direction = i
				break
			}
		}

		o, found := map[rune][2]int{'>': {0, 1}, '^': {-1, 0}, 'v': {1, 0}, '<': {0, -1}}[maze[curr.point.x][curr.point.y]]
		curr.point.x += o[0]
		curr.point.y += o[1]

		if found {
			graph[start] = append(graph[start], Edge{curr.point, steps + 2})
			if *secondFlag {
				graph[curr.point] = append(graph[curr.point], Edge{start, steps + 2})
			}
		} else {
			graph[start] = append(graph[start], Edge{curr.point, steps})
			if *secondFlag {
				graph[curr.point] = append(graph[curr.point], Edge{start, steps})
			}
		}

		if visited[curr.point] {
			continue
		}
		visited[curr.point] = true
		if found {
			for i, c := range []rune{'>', '^', 'v', '<'} {
				nextX, nextY := curr.point.x+offset[i][0], curr.point.y+offset[i][1]
				if nextX >= 0 && nextX < len(maze) && nextY >= 0 && nextY < len(maze[0]) && c == maze[nextX][nextY] {
					queue = append(queue, QueueInfo{Point{nextX, nextY}, i})
				}
			}
		}
	}

	maxCost := 0
	visited = make(map[Point]bool)
	dfs(Point{0, 1}, Point{len(maze) - 1, len(maze[0]) - 2}, graph, visited, 0, &maxCost)
	fmt.Println(maxCost)
}

func dfs(point Point, dest Point, graph map[Point][]Edge, visited map[Point]bool, currCost int, maxCost *int) {
	if point == dest {
		if currCost > *maxCost {
			*maxCost = currCost
		}
	}

	for _, edge := range graph[point] {
		if !visited[edge.dest] {
			visited[edge.dest] = true
			dfs(edge.dest, dest, graph, visited, currCost+edge.cost, maxCost)
			visited[edge.dest] = false
		}
	}
}
