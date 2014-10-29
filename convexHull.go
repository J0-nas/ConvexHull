package main

import (
	"fmt"
	"sort"
)

type Point struct {
	X, Y int
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

	p1 := points[len(points)-1]

	p1origin := p1
}
