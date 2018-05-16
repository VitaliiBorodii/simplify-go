package simplifyGo

type ChartPoint struct {
	X float64
	Y float64
}

func (p ChartPoint) GetX() float64 {
	return p.X
}

func (p ChartPoint) GetY() float64 {
	return p.Y
}

type Point interface {
	GetX() float64
	GetY() float64
}

func getSqDist(p1 Point, p2 Point) float64 {

	dx := p1.GetX() - p2.GetX()
	dy := p1.GetY() - p2.GetY()

	return dx*dx + dy*dy
}

func getSqSegDist(point *Point, point1 *Point, point2 *Point) float64 {
	p := *point
	p1 := *point1
	p2 := *point2

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

func simplifyRadialDist(points *[]Point, sqTolerance float64) *[]Point {
	arr := *points

	prevPoint := arr[0]
	newPoints := []Point{prevPoint}
	var point Point

	for i := 1; i < len(arr); i++ {
		point = arr[i]

		if getSqDist(point, prevPoint) > sqTolerance {
			newPoints = append(newPoints, point)
			prevPoint = point
		}
	}

	if &prevPoint != &point {
		newPoints = append(newPoints, point)
	}

	return &newPoints
}

func simplifyDPStep(points *[]Point, first int, last int, sqTolerance float64, simplified *[]Point) {
	arr := *points
	maxSqDist := sqTolerance
	var index int
	simpifiedVal := *simplified

	for i := first + 1; i < last; i++ {
		sqDist := getSqSegDist(&arr[i], &arr[first], &arr[last])

		if sqDist > maxSqDist {
			index = i
			maxSqDist = sqDist
		}
	}

	if maxSqDist > sqTolerance {
		if index-first > 1 {
			simplifyDPStep(points, first, index, sqTolerance, simplified)
		}
		simpifiedVal = append(simpifiedVal, arr[index])
		if last-index > 1 {
			simplifyDPStep(points, index, last, sqTolerance, simplified)
		}
	}
}

func simplifyDouglasPeucker(points *[]Point, sqTolerance float64) *[]Point {
	arr := *points
	last := len(arr) - 1

	simplifiedVal := []Point{arr[0]}
	simplified := &simplifiedVal
	simplifyDPStep(points, 0, last, sqTolerance, simplified)
	simplifiedVal = append(simplifiedVal, arr[last])

	return simplified
}

func Simplify(points *[]Point, tolerance float64, highestQuality bool) *[]Point {
	if len(*points) <= 2 {
		return points
	}

	var sqTolerance float64
	if tolerance == 0 {
		sqTolerance = 1
	} else {
		sqTolerance = tolerance * tolerance
	}

	if !highestQuality {
		points = simplifyRadialDist(points, sqTolerance)
	}

	points = simplifyDouglasPeucker(points, sqTolerance)

	return points
}
