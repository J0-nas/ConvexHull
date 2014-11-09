package main

import (
	"errors"
	"fmt"
	"github.com/J0-nas/drawGraph"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
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
	var wG sync.WaitGroup
	workChannel := make(chan int)

	m := make(map[int]int)
	for worker := 0; worker < 100; worker++ {
		// Increment the WaitGroup counter.
		wG.Add(1)
		go func(wChannel <-chan int, id int, waitGrp *sync.WaitGroup) {
			defer waitGrp.Done()
			completedWork := 0
			for i := range wChannel {
				//fmt.Println("Worker ", id, " working on ", i)
				convex("convexHull_" + strconv.Itoa(i))
				completedWork++
			}
			m[id] = completedWork
			//fmt.Println("Done.")
		}(workChannel, worker, &wG)
	}

	wG.Add(1)
	go func(wChannel chan int, wGrp *sync.WaitGroup) {
		defer wGrp.Done()
		for work := 0; work < 10000; work++ {
			wChannel <- work
		}
		//fmt.Println("Create work done.")
		close(workChannel)
	}(workChannel, &wG)

	wG.Wait()

	var res []Point
	for key, value := range m {
		res = append(res, Point{float64(key), float64(value)})
	}
	sort.Sort(ByX(res))
	for _, r := range res {
		fmt.Println("Key: ", int(r.X), " computed: ", int(r.Y), " convex hulls")
	}

}

func convex(name string) {
	height := 1000.0
	width := 1000.0

	var points []Point

	seed := time.Now().UTC().UnixNano()
	rand.Seed(seed)
	//fmt.Println("seed: ", seed)

	for start := 0; start < int(rand.Float64()*1000)+3; start++ {
		points = append(points, Point{rand.Float64() * width, rand.Float64() * height})
	}

	d := drawGraph.NewImage(height, width)
	for i := range points {
		err := d.AddPoint(points[i].X, points[i].Y)
		if err != nil {
			log.Println(err)
		}
	}

	convexHull, err := buildConvexHull_2(points)
	//convexHull := buildConvexHull_2(points)
	if err != nil {
		log.Println(err)
	}

	if len(convexHull) <= 2 {
		log.Println("Too few points in convex hull")
		os.Exit(1)
	}
	for j := 0; j <= len(convexHull)-2; j++ {
		err := d.AddLine(convexHull[j].X, convexHull[j].Y, convexHull[j+1].X, convexHull[j+1].Y)
		if err != nil {
			fmt.Println("\n\nError!\n\n")
			log.Println(err)
		}
	}
	d.SaveImage(name)
}

func buildConvexHull_2(points []Point) (result []Point, err error) {
	sort.Sort(ByX(points))
	//fmt.Println("sorted ", points)
	if len(points) < 3 {
		log.Println(points)
		return points, errors.New("Too few points.")
	}
	origin := points[0]
	nextToOrigin := points[1]
	currentPoint := origin
	nextPoint := nextToOrigin
	points = points[2:]

	result = []Point{
		currentPoint,
		nextPoint,
	}
	left := 0.0
	for _, p := range points {
		left = ccw(currentPoint, nextPoint, p)
		if left < 0.0 {
			l := len(result) - 1
			if l == 1 { //only 2 elements in result, special treatment
				result = append(result[:1], p)
				nextPoint = p
			} else {
				for j := l; j >= 1; j-- {
					if ccw(result[j-1], result[j], p) >= 0 {
						result = append(result[:j+1], p)
						//fmt.Println("new hull: ", result)
						currentPoint = result[j]
						nextPoint = p
						break
					}
					if j == 1 {
						result = append(result[:1], p)
						nextPoint = p
						currentPoint = result[0]
					}
				}
			}
		} else {
			result = append(result, p)
			currentPoint = nextPoint
			nextPoint = p
		}
	}

	currentPoint = origin
	nextPoint = nextToOrigin
	result_2 := []Point{
		currentPoint,
		nextPoint,
	}

	right := 0.0
	for _, p := range points {
		right = ccw(currentPoint, nextPoint, p)
		if right > 0.0 {
			l := len(result_2) - 1
			if l == 1 { //only 2 elements in result, special treatment
				result_2 = append(result_2[:1], p)
				nextPoint = p
			} else {
				for j := l; j >= 1; j-- {
					if ccw(result_2[j-1], result_2[j], p) <= 0 {
						result_2 = append(result_2[:j+1], p)
						//fmt.Println("new hull: ", result_2)
						currentPoint = result_2[j]
						nextPoint = p
						break
					}
					if j == 1 {
						result_2 = append(result_2[:1], p)
						nextPoint = p
						currentPoint = result_2[0]
					}
				}
			}
		} else {
			result_2 = append(result_2, p)
			currentPoint = nextPoint
			nextPoint = p
		}
	}

	for h := len(result_2) - 1; h > 0; h-- {
		result = append(result, result_2[h])
	}
	result = append(result, origin)

	return result, nil
}

/*(p1, p2, p3) is counterclockwise iff result > 0 => p3 is on the righthand side*/
func ccw(p1, p2, p3 Point) float64 {
	i := (p1.X-p2.X)*(p3.Y-p2.Y) - (p3.X-p2.X)*(p1.Y-p2.Y)
	//fmt.Println(i)
	return i
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

/*points = []Point{
	{100, 200},
	{800, 200},
	{600, 300},
	{300, 500},
	{700, 700},
	{700, 100},
	{400, 900},
}

result := []Point{
	{100, 200},
	{600, 300},
	{300, 500},
	{700, 700},
	{400, 900},
}

sort.Sort(ByX(points))
sort.Sort(ByX(result))

fmt.Println(points)
for i, p := range points {
	for _, r := range result {
		if p.X > r.X {
			break
		} else if p.Equals(r) {
			if i == len(points)-1 {
				points = points[:i]
			} else if i == 0 {
				points = points[1:]
			} else {
				points = append(points[:i], points[i+1:]...)
			}
		}
	}
}
fmt.Println(result)
fmt.Println(points)*/
