package main

import (
	"fmt"
	"github.com/J0-nas/drawGraph"
	"log"
	"math"
	"math/rand"
	"os"
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
	pTemp := a[i]
	a[i] = a[j]
	a[j] = pTemp
	//	, a[j] = a[j], a[i]
}

func (a ByX) Less(i, j int) bool {
	if a[i].X == a[j].X {
		return a[i].Y < a[j].Y
	}
	return a[i].X < a[j].X
}

func main() {
	fmt.Println("Test output!")

	var points []Point
	height := 1000.0
	width := 1000.0

	for start := 0; start < 20; start++ {
		points = append(points, Point{rand.Float64() * width, rand.Float64() * height})
	}

	/*points := []Point{
		{10, 20},
		{80, 20},
		{60, 30},
		{30, 50},
		{70, 70},
		{70, 10},
		{40, 90},
	}*/
	d := drawGraph.Instance()
	d.NewImage(height, width)
	for i := range points {
		err := d.AddPoint(points[i].X, points[i].Y)
		if err != nil {
			log.Println(err)
		}
	}

	convexHull := buildConvexHull_2(points)

	if len(convexHull) <= 2 {
		log.Println("Too few points in convex hull")
		os.Exit(1)
	}
	for j := 0; j <= len(convexHull)-1; j++ {
		if j == len(convexHull)-1 {
			err := d.AddLine(convexHull[j].X, convexHull[j].Y, convexHull[0].X, convexHull[0].Y)
			if err != nil {
				log.Println(err)
			}
		} else {
			err := d.AddLine(convexHull[j].X, convexHull[j].Y, convexHull[j+1].X, convexHull[j+1].Y)
			if err != nil {
				log.Println(err)
			}
		}
	}
	//d.AddPoint(50.0, 100.0)
	//d.AddLine(50.0, 100.0, 100.0, 100.0)
	d.SaveImage("test")

}

func buildConvexHull_2(points []Point) []Point {
	sort.Sort(ByX(points))
	currentPoint := points[0]
	originPoint := points[0]
	points = points[1:]

	results := []Point{
		originPoint,
	}
	maxAngle := 0.0
	var maxPos int
	for {
		for i, p := range points {
			a := math.Atan((p.Y - currentPoint.Y) / (p.X - currentPoint.X))
			if a > maxAngle {
				maxAngle = a
				maxPos = i
			}
		}
		results = append(results, points[maxPos])
		currentPoint = points[maxPos]
		points = append(points[:maxPos], points[maxPos+1:]...)
		if currentPoint.Equals(originPoint) {
			break
		}
	}
	return results
}

func buildConvexHull_1(points []Point) []Point {
	//fmt.Println(points)

	sort.Sort(ByX(points))

	//fmt.Println(points)

	length := len(points)

	currentPoint := points[length-1]
	originPoint := points[length-1]

	points = points[:(length - 1)]

	//fmt.Println(points)

	var a, c float64
	results := []Point{
		originPoint,
	}
	atOrigin := false

	for alpha := math.Pi/2 - 0.001; alpha > -(math.Pi + math.Pi/2 + 0.001); alpha -= 0.001 {
		a = -math.Tan(alpha)
		c = currentPoint.Y + a*currentPoint.X
		for i, p := range points {
			if checkPoint(a, c, p) {
				currentPoint = p
				if currentPoint == originPoint {
					atOrigin = true
				}
				points = append(points[:i], points[i+1:]...)
				fmt.Println("Added point... ", p)
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
	fmt.Print("Convec Hull: ")
	fmt.Println(results)
	return results
}

func distance(a, c float64, p Point) float64 {
	return math.Abs((a*p.X + p.Y - c) / math.Sqrt(a*a+1))
}

func checkPoint(a, c float64, p Point) bool {
	return distance(a, c, p) <= 0.2
}

func (p *Point) Equals(q Point) bool {
	return p.X == q.X && p.Y == q.Y
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
