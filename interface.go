package main

import (
	"fmt"
	"math"
)

type Circle struct {
	x, y, r float64
}

type Rectangle struct {
	x1, y1, x2, y2 float64
}

func distance(x, y float64) float64 {
	return math.Abs(x - y)
}

func (r *Rectangle) Area() float64 {
	l := distance(r.x1, r.x2)
	w := distance(r.y1, r.y2)
	return l * w
}

type Shape interface {
	Area() float64
}

func TotalArea(shapes ...Shape) {
	var ret float64

	for _, s := range shapes {
		ret += s.Area()
	}
}

func (c *Circle) Area() float64 {
	return math.Pi * c.r * c.r
}

func CircleArea(c Circle) float64 {
	return math.Pi * c.r * c.r
}

func main() {
	c := Circle{1, 1, 5}
	fmt.Println(c, CircleArea(c), c.Area())

	r := Rectangle{0, 0, 10, 10}
	fmt.Println(r, r.Area())
}
