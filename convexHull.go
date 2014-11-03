package main

import (
	"fmt"
	"github.com/J0-nas/drawGraph"
	"log"
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
		{10, 20},
		{80, 20},
		{60, 30},
		{30, 50},
		{70, 70},
		{70, 10},
		{40, 90},
	}

	convexHull := buildConvexHull(points)
	fmt.Println(convexHull)

	d := drawGraph.Instance()
	d.NewImage(100, 100)
	for j := 0; j < len(convexHull)-1; j++ {
		if j == len(convexHull)-1 {
			d.AddLine(convexHull[0].X, convexHull[0].Y, convexHull[j].X, convexHull[j].Y)
		} else {
			d.AddLine(convexHull[j].X, convexHull[j].Y, convexHull[j+1].X, convexHull[j+1].Y)
		}
	}
	for j := 0; j < len(points)-1; j++ {
		err := d.AddPoint(convexHull[j].X, convexHull[j].Y)
		if err != nil {
			log.Println(err)
		}
	}
	//d.AddPoint(50.0, 100.0)
	//d.AddLine(50.0, 100.0, 100.0, 100.0)
	d.SaveImage("test")

}

func buildConvexHull(points []Point) []Point {
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
	return results
}

func distance(a, c float64, p Point) float64 {
	return math.Abs((a*p.X + p.Y - c) / math.Sqrt(a*a+1))
}

func checkPoint(a, c float64, p Point) bool {
	return distance(a, c, p) <= 0.7
}

//d.PrintSomething()
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
