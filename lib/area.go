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

func CreateEventArea(x, y int) *Area {
	return &Area{
		delegate: image.Rect(x-EventAreaRadius, y-EventAreaRadius, x+EventAreaRadius, y+EventAreaRadius),
	}
}
