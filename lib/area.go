package lib

import (
	"image"
)

type Area struct {
	delegate image.Rectangle
}

func (a *Area) IntersectsArea(other *Area) bool {
	return a.delegate.Overlaps(other.delegate)
}

func CreateEventArea(x, y, radius int) *Area {
	return &Area{
		delegate: image.Rect(x-radius, y-radius, x+radius, y+radius),
	}
}
