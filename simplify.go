package simplifyGo

import "fmt"

type ChartPoint struct {
	X float64
	Y float64
}

type Point interface {
	GetX() float64
	GetY() float64
}

func (p ChartPoint) GetX() float64 {
	return p.X
}

func (p ChartPoint) GetY() float64 {
	return p.Y
}

func getSqDist(p1 Point, p2 Point) float64 {

	dx := p1.GetX() - p2.GetX()
	dy := p1.GetY() - p2.GetY()

	return dx*dx + dy*dy
}

func getSqSegDist(p Point, p1 Point, p2 Point) float64 {

	x := p1.GetX()
	y := p1.GetY()
	dx := p2.GetX() - x
	dy := p2.GetY() - y

	if dx != 0 || dy != 0 {
		t := ((p.GetX()-x)*dx + (p.GetY()-y)*dy) / (dx*dx + dy*dy)

		if t > 1 {
			x = p2.GetX()
			y = p2.GetY()
		} else if t > 0 {
			x += dx * t
			y += dy * t
		}
	}

	dx = p.GetX() - x
	dy = p.GetY() - y

	return dx*dx + dy*dy
}

func simplifyRadialDist(points []Point, sqTolerance float64) []Point {

	prevPoint := points[0]
	newPoints := []Point{prevPoint}
	var point Point

	for i := 1; i < len(points); i++ {
		point = points[i]

		if getSqDist(point, prevPoint) > sqTolerance {
			newPoints = append(newPoints, point)
			prevPoint = point
		}
	}

	if &prevPoint != &point {
		newPoints = append(newPoints, point)
	}

	return newPoints
}

func simplifyDPStep(points []Point, first int, last int, sqTolerance float64, simplified []Point) []Point {
	maxSqDist := sqTolerance
	var index int

	for i := first + 1; i < last; i++ {
		sqDist := getSqSegDist(points[i], points[first], points[last])

		if sqDist > maxSqDist {
			index = i
			maxSqDist = sqDist
		}
	}

	if maxSqDist > sqTolerance {
		if index-first > 1 {
			simplified = simplifyDPStep(points, first, index, sqTolerance, simplified)
		}
		simplified = append(simplified, points[index])
		if last-index > 1 {
			simplified = simplifyDPStep(points, index, last, sqTolerance, simplified)
		}
	}

	return simplified
}

func simplifyDouglasPeucker(points []Point, sqTolerance float64) []Point {
	last := len(points) - 1

	simplified := []Point{points[0]}
	simplified = simplifyDPStep(points, 0, last, sqTolerance, simplified)
	simplified = append(simplified, points[last])

	return simplified
}

func Simplify(points *[]Point, tolerance float64, highestQuality bool) []Point {
	arr := *points

	if len(arr) <= 2 {
		return arr
	}

	var sqTolerance float64
	if tolerance == 0 {
		sqTolerance = 1
	} else {
		sqTolerance = tolerance * tolerance
	}

	if !highestQuality {
		arr = simplifyRadialDist(arr, sqTolerance)
	}

	arr = simplifyDouglasPeucker(arr, sqTolerance)
	
	fmt.Println("Simplify:", len(points), "=>", len(arr))

	return arr
}
