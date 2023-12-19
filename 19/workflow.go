package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Workflow struct {
	Rules             []Rule
	WorkflowIfNoMatch string
}

type Rule struct {
	Char        byte
	Lt          bool
	TargetValue int
	ToWorkflow  string
}

type Range struct {
	L, R int
}

type SearchState struct {
	WorkflowName string
	Ranges       map[byte]Range
}

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("workflow.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	workflows := make(map[string]Workflow)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		workflowStr := scanner.Text()
		if workflowStr == "" {
			break
		}

		fields := strings.Split(workflowStr, "{")
		workflowName := fields[0]
		fields = strings.Split(fields[1], ",")
		rules := make([]Rule, len(fields)-1)
		for i := 0; i < len(fields)-1; i++ {
			ruleStr := fields[i]
			ruleFields := strings.Split(ruleStr, ":")
			char := ruleFields[0][0]
			lt := ruleFields[0][1] == '<'
			targetValue, _ := strconv.Atoi(ruleFields[0][2:])
			toWorkflow := ruleFields[1]
			rule := Rule{
				Char:        char,
				Lt:          lt,
				TargetValue: targetValue,
				ToWorkflow:  toWorkflow,
			}
			rules[i] = rule
		}
		lastIndex := len(fields) - 1
		workflowIfNoMatch := fields[lastIndex][:len(fields[lastIndex])-1]

		workflow := Workflow{
			Rules:             rules,
			WorkflowIfNoMatch: workflowIfNoMatch,
		}
		workflows[workflowName] = workflow
	}

	if !*secondFlag {
		// Part 1
		sum := 0
		for scanner.Scan() {
			dictString := scanner.Text()
			fields := strings.Split(dictString[1:len(dictString)-1], ",")
			part := make(map[byte]int)
			for _, field := range fields {
				keyValue := strings.Split(field, "=")
				key := keyValue[0][0]
				value, _ := strconv.Atoi(keyValue[1])
				part[key] = value
			}

			currWorkflowName := "in"
			for currWorkflowName != "A" && currWorkflowName != "R" {
				matched := false
				currWorkflow := workflows[currWorkflowName]
				for _, rule := range currWorkflow.Rules {
					if (rule.Lt && part[rule.Char] < rule.TargetValue) || (!rule.Lt && part[rule.Char] > rule.TargetValue) {
						currWorkflowName = rule.ToWorkflow
						matched = true
						break
					}
				}
				if !matched {
					currWorkflowName = currWorkflow.WorkflowIfNoMatch
				}
			}

			if currWorkflowName == "A" {
				for _, value := range part {
					sum += value
				}
			}
		}
		fmt.Println(sum)
		return
	}

	// Part 2
	sum := 0
	var queue []SearchState
	initialState := SearchState{
		WorkflowName: "in",
		Ranges:       make(map[byte]Range),
	}
	for _, c := range "xmas" {
		initialState.Ranges[byte(c)] = Range{1, 4000}
	}
	queue = append(queue, initialState)
	for len(queue) > 0 {
		// Remove from queue
		state := queue[0]
		queue = queue[1:]

		// Process accepted/rejected state
		if state.WorkflowName == "A" {
			product := 1
			for _, r := range state.Ranges {
				product *= r.R - r.L + 1
			}
			sum += product
			continue
		}
		if state.WorkflowName == "R" {
			continue
		}

		workflow := workflows[state.WorkflowName]
		noMatchRangeValid := true
		for _, rule := range workflow.Rules {
			var matchRange, noMatchRange Range
			if rule.Lt {
				matchRange = Range{state.Ranges[rule.Char].L, rule.TargetValue - 1}
				noMatchRange = Range{rule.TargetValue, state.Ranges[rule.Char].R}
			} else {
				matchRange = Range{rule.TargetValue + 1, state.Ranges[rule.Char].R}
				noMatchRange = Range{state.Ranges[rule.Char].L, rule.TargetValue}
			}

			// Check if match range is valid
			if matchRange.L <= matchRange.R {
				newRanges := make(map[byte]Range)
				for k, v := range state.Ranges {
					newRanges[k] = v
				}
				newRanges[rule.Char] = matchRange
				queue = append(queue, SearchState{
					WorkflowName: rule.ToWorkflow,
					Ranges:       newRanges,
				})
			}

			// Check if no match range is valid
			if noMatchRange.L > noMatchRange.R {
				noMatchRangeValid = false
				break
			} else {
				state.Ranges[rule.Char] = noMatchRange
			}
		}
		if noMatchRangeValid {
			queue = append(queue, SearchState{
				WorkflowName: workflow.WorkflowIfNoMatch,
				Ranges:       state.Ranges,
			})
		}
	}
	fmt.Println(sum)
}
