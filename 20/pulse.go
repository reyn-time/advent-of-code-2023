package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type ModuleType int

const (
	ModuleTypeFlipFlop    ModuleType = iota
	ModuleTypeConjunction ModuleType = 1
	ModuleTypeBroadcast   ModuleType = 2
)

type Module struct {
	Type        ModuleType
	SourceNames []string
	TargetNames []string
	Memory      map[string]bool
	Toggled     bool
}

type Pulse struct {
	IsHigh                 bool
	SourceName, TargetName string
}

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("pulse.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	nameToModule := make(map[string]*Module)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " -> ")
		moduleName := parts[0][1:]
		moduleType := parts[0][0]
		var mt ModuleType
		switch moduleType {
		case '%':
			mt = ModuleTypeFlipFlop
		case '&':
			mt = ModuleTypeConjunction
		case 'b':
			mt = ModuleTypeBroadcast
		}
		if mt == ModuleTypeBroadcast {
			moduleName = "broadcaster"
		}

		targets := strings.Split(parts[1], ", ")
		nameToModule[moduleName] = &Module{
			Type:        mt,
			SourceNames: []string{},
			TargetNames: targets,
			Memory:      make(map[string]bool),
		}
	}

	// Populate source names for each module
	for name, module := range nameToModule {
		for _, target := range module.TargetNames {
			targetModule := nameToModule[target]
			if targetModule == nil {
				continue
			}
			targetModule.SourceNames = append(targetModule.SourceNames, name)
		}
	}

	// Part 1
	if !*secondFlag {
		high, low := 0, 0
		for i := 0; i < 1000; i++ {
			h, l, _ := pressButton(nameToModule, nil)
			high += h
			low += l
		}
		fmt.Println(high * low)
		return
	}

	// Part 2
	// Find module with 'rx' as target
	var rootModule *Module
	for name, module := range nameToModule {
		for _, target := range module.TargetNames {
			if target == "rx" {
				rootModule = nameToModule[name]
			}
		}
	}
	rootModuleSources := rootModule.SourceNames

	period := map[string]int{}
	i := 0
	for len(period) != len(rootModuleSources) {
		i++
		_, _, highAtRoot := pressButton(nameToModule, rootModule)
		for name := range highAtRoot {
			if _, ok := period[name]; !ok {
				period[name] = i
			}
		}
	}

	product := 1
	for _, p := range period {
		product = lcm(product, p)
	}
	fmt.Println(product)
}

func pressButton(nameToModule map[string]*Module, rootModule *Module) (int, int, map[string]bool) {
	queue := []Pulse{{SourceName: "button", TargetName: "broadcaster", IsHigh: false}}
	high, low := 0, 0
	highAtRoot := make(map[string]bool)

	for len(queue) > 0 {
		pulse := queue[0]
		queue = queue[1:]

		module := nameToModule[pulse.TargetName]

		if module == rootModule && pulse.IsHigh {
			highAtRoot[pulse.SourceName] = true
		}

		if pulse.IsHigh {
			high++
		} else {
			low++
		}

		if module == nil {
			continue
		}

		if module.Type == ModuleTypeFlipFlop {
			if !pulse.IsHigh {
				module.Toggled = !module.Toggled
				for _, target := range module.TargetNames {
					queue = append(queue, Pulse{SourceName: pulse.TargetName, TargetName: target, IsHigh: module.Toggled})
				}
			}
		} else if module.Type == ModuleTypeConjunction {
			module.Memory[pulse.SourceName] = pulse.IsHigh
			allHigh := true
			for _, source := range module.SourceNames {
				if !module.Memory[source] {
					allHigh = false
					break
				}
			}
			for _, target := range module.TargetNames {
				queue = append(queue, Pulse{SourceName: pulse.TargetName, TargetName: target, IsHigh: !allHigh})
			}
		} else if module.Type == ModuleTypeBroadcast {
			for _, target := range module.TargetNames {
				queue = append(queue, Pulse{SourceName: pulse.TargetName, TargetName: target, IsHigh: false})
			}
		}
	}

	return high, low, highAtRoot
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}

	return a
}
