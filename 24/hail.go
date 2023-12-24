package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"

	"gonum.org/v1/gonum/mat"
)

type Vector3D struct {
	x, y, z int
}

type Hailstone struct {
	position, velocity Vector3D
}

func main() {
	secondFlag := flag.Bool("second", false, "Run the solution for the second half of the puzzle")
	flag.Parse()

	file, err := os.Open("hail.in")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	hailstones := make([]Hailstone, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var x, y, z, tx, ty, tz int
		fmt.Sscanf(scanner.Text(), "%d, %d, %d @ %d, %d, %d", &x, &y, &z, &tx, &ty, &tz)
		hailstone := Hailstone{Vector3D{x, y, z}, Vector3D{tx, ty, tz}}
		hailstones = append(hailstones, hailstone)
	}

	if !*secondFlag {
		minBound, maxBound := float64(200000000000000), float64(400000000000000)
		counter := 0
		for i := 0; i < len(hailstones); i++ {
			for j := i + 1; j < len(hailstones); j++ {
				// Check if the two hailstones are on a collision course
				a := hailstones[i]
				b := hailstones[j]
				x := approximateCollisionTimes(a, b)
				t1 := x.At(0, 0)
				t2 := x.At(1, 0)
				if t1 > 0 && t2 > 0 {
					ax := float64(a.position.x) + float64(a.velocity.x)*t1
					ay := float64(a.position.y) + float64(a.velocity.y)*t1
					if ax >= minBound && ax <= maxBound && ay >= minBound && ay <= maxBound {
						counter++
					}
				}
			}
		}
		fmt.Println(counter)
		return
	}
	for testX := -500; testX <= 500; testX++ {
		for testY := -500; testY <= 500; testY++ {
			a := hailstones[0]
			b := hailstones[1]
			a.velocity.x -= testX
			a.velocity.y -= testY
			b.velocity.x -= testX
			b.velocity.y -= testY
			x := approximateCollisionTimes(a, b)
			t1 := x.At(0, 0)
			t2 := x.At(1, 0)
			if !(t1 > 0 && t2 > 0 && isInteger(t1) && isInteger(t2)) {
				continue
			}
			T := toInteger(t1)
			collision := Vector3D{a.position.x + (a.velocity.x+testX)*T, a.position.y + (a.velocity.y+testY)*T, a.position.z + a.velocity.z*T}
			for testZ := -500; testZ <= 500; testZ++ {
				startPos := Vector3D{collision.x - testX*T, collision.y - testY*T, collision.z - testZ*T}
				testVelocity := Vector3D{testX, testY, testZ}
				hailstone := Hailstone{startPos, testVelocity}
				canCollide := true
				for i := 0; i < 100; i++ {
					if !canActuallyCollide(hailstone, hailstones[i]) {
						canCollide = false
						break
					}
				}
				if canCollide {
					fmt.Println(startPos.x + startPos.y + startPos.z)
					return
				}
			}
		}
	}
}

func approximateCollisionTimes(a, b Hailstone) mat.Dense {
	A := mat.NewDense(2, 2, []float64{float64(a.velocity.x), float64(-b.velocity.x), float64(a.velocity.y), float64(-b.velocity.y)})
	B := mat.NewDense(2, 1, []float64{float64(b.position.x - a.position.x), float64(b.position.y - a.position.y)})
	x := mat.NewDense(2, 1, nil)
	if mat.Det(A) == 0 {
		return *mat.NewDense(2, 1, []float64{-1, -1})
	}
	x.Solve(A, B)
	return *x
}

func isInteger(x float64) bool {
	return math.Abs(x-math.Round(x)) < 0.1
}

func toInteger(x float64) int {
	return int(math.Round(x))
}

func canActuallyCollide(a, b Hailstone) bool {
	pos := []int{a.position.x - b.position.x, a.position.y - b.position.y, a.position.z - b.position.z}
	vel := []int{a.velocity.x - b.velocity.x, a.velocity.y - b.velocity.y, a.velocity.z - b.velocity.z}
	prevTime := -1
	for i := 0; i < 3; i++ {
		if pos[i] == 0 {
			if vel[i] == 0 {
				continue
			}
			return false
		}
		if vel[i] == 0 {
			return false
		}
		if pos[i]%vel[i] != 0 {
			return false
		}
		time := pos[i] / vel[i]
		if prevTime != -1 && time != prevTime {
			return false
		}
		prevTime = time
	}
	return true
}
