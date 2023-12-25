package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("cut.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	graph := make(map[string]map[string]int)
	removed := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ": ")
		nodeName := fields[0]
		friends := strings.Split(fields[1], " ")
		for _, friend := range friends {
			if graph[nodeName] == nil {
				graph[nodeName] = make(map[string]int)
			}
			graph[nodeName][friend] = 1
			if graph[friend] == nil {
				graph[friend] = make(map[string]int)
			}
			graph[friend][nodeName] = 1
		}
	}

	// Run simple min-cut algorithm
	groupSizes := make(map[string]int)
	for node := range graph {
		groupSizes[node] = 1
	}
	for {
		t, s, l := suggestMerge(graph, removed)
		if l == 3 || l == 0 {
			fmt.Println(groupSizes[t] * (len(groupSizes) - groupSizes[t]))
			break
		}
		groupSizes[s] = groupSizes[s] + groupSizes[t]
		removed[t] = true
		for neighbour, dist := range graph[t] {
			if removed[neighbour] {
				continue
			}
			graph[neighbour][s] = graph[neighbour][s] + dist
			graph[s][neighbour] = graph[s][neighbour] + dist
		}
	}

}

func suggestMerge(graph map[string]map[string]int, removed map[string]bool) (string, string, int) {
	startNode := ""
	count := 0
	for node := range graph {
		if !removed[node] {
			startNode = node
			count++
		}
	}
	partition := map[string]bool{}
	l := []string{startNode}
	neighbourDists := map[string]int{}
	largestNeighbour := startNode
	largestDist := 0
	for i := 0; i < count-1; i++ {
		partition[largestNeighbour] = true
		for neighbour, dist := range graph[largestNeighbour] {
			if partition[neighbour] || removed[neighbour] {
				continue
			}
			neighbourDists[neighbour] = neighbourDists[neighbour] + dist
		}

		largestDist = 0
		for n, dist := range neighbourDists {
			if dist > largestDist {
				largestNeighbour = n
				largestDist = dist
			}
		}
		l = append(l, largestNeighbour)
		delete(neighbourDists, largestNeighbour)
	}
	return l[len(l)-1], l[len(l)-2], largestDist
}
