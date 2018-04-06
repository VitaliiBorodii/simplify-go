package simplifyGo

type ChartPoint struct {
	X float64
	Y float64
}

type Point interface {
	getX() float64
	getY() float64
}

func (p ChartPoint) getX() float64 {
	return p.X
}

func (p ChartPoint) getY() float64 {
	return p.Y
}

func getSqDist(p1 ChartPoint, p2 ChartPoint) float64 {

	dx := p1.getX() - p2.getX()
	dy := p1.getY() - p2.getY()

	return dx*dx + dy*dy
}

func getSqSegDist(p ChartPoint, p1 ChartPoint, p2 ChartPoint) float64 {

	x := p1.getX()
	y := p1.getY()
	dx := p2.getX() - x
	dy := p2.getY() - y

	if dx != 0 || dy != 0 {
		t := ((p.getX()-x)*dx + (p.getY()-y)*dy) / (dx*dx + dy*dy)

		if t > 1 {
			x = p2.getX()
			y = p2.getY()
		} else if t > 0 {
			x += dx * t
			y += dy * t
		}
	}

	dx = p.getX() - x
	dy = p.getY() - y

	return dx*dx + dy*dy
}

func simplifyRadialDist(points []ChartPoint, sqTolerance float64) []ChartPoint {

	prevPoint := points[0]
	newPoints := []ChartPoint{prevPoint}
	var point ChartPoint

	for i := 1; i < len(points); i++ {
		point = points[i]

		if getSqDist(point, prevPoint) > sqTolerance {
			newPoints = append(newPoints, point)
			prevPoint = point
		}
	}

	if prevPoint != point {
		newPoints = append(newPoints, point)
	}

	return newPoints
}

func simplifyDPStep(points []ChartPoint, first int, last int, sqTolerance float64, simplified []ChartPoint) []ChartPoint {
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

func simplifyDouglasPeucker(points []ChartPoint, sqTolerance float64) []ChartPoint {
	last := len(points) - 1

	simplified := []ChartPoint{points[0]}
	simplified = simplifyDPStep(points, 0, last, sqTolerance, simplified)
	simplified = append(simplified, points[last])

	return simplified
}

func Simplify(points []ChartPoint, tolerance float64, highestQuality bool) []ChartPoint {

	if len(points) <= 2 {
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
