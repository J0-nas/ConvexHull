package main

import (
	"fmt"
	"math"
	"sort"
)

type Point struct {
	X, Y float64
}

type ByX []Point

func (a ByX) Len() int {
	return len(a)
}

func (a ByX) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByX) Less(i, j int) bool {
	if a[i].X == a[j].X {
		return a[i].Y < a[j].Y
	}
	return a[i].X < a[j].X
}

func main() {
	fmt.Println("Test output!")

	points := []Point{
		{1, 2},
		{8, 2},
		{6, 3},
		{3, 5},
		{7, 7},
		{7, 1},
		{4, 9},
	}

	fmt.Println(points)

	sort.Sort(ByX(points))

	fmt.Println(points)

	length := len(points)

	currentPoint := points[length-1]
	originPoint := points[length-1]

	points = points[:(length - 1)]

	fmt.Println(points)

	var a, c float64
	results := []Point{
		originPoint,
	}
	atOrigin := false

	//TEST
	/*alpha := math.Pi/2 - 0.01
	a = math.Tan(alpha)
	c = currentPoint.Y + a*currentPoint.X

	fmt.Println(a)
	fmt.Println()
	fmt.Println(c)
	fmt.Println()

	for _, p := range points {
		fmt.Println(distance(a, c, p))

		if checkPoint(a, c, p) {
			fmt.Println(p)
		}
	}*/

	for alpha := math.Pi/2 - 0.1; alpha > -math.Pi/2+0.1; alpha -= 0.02 {
		a = -math.Tan(alpha)
		c = currentPoint.Y + a*currentPoint.X
		for i, p := range points {
			if checkPoint(a, c, p) {
				currentPoint = p
				if currentPoint == originPoint {
					atOrigin = true
				}
				points = append(points[:i], points[i+1:]...)
				results = append(results, p)

				fmt.Println(a)
				fmt.Println()
				fmt.Println(c)
				fmt.Println()
				fmt.Println()

				break
			}
		}
		if atOrigin {
			break
		}
	}

	fmt.Println(results)

}

func distance(a, c float64, p Point) float64 {
	return math.Abs((a*p.X + p.Y - c) / math.Sqrt(a*a+1))
}

func checkPoint(a, c float64, p Point) bool {
	return distance(a, c, p) <= 0.7
}
