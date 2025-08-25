package main_test

import (
	"github.com/cubeee/ent-notifier/lib"
	"image/png"
	"os"
	"testing"
)

func TestThumbnailGeneration(t *testing.T) {
	thumbnail, err := lib.CreateThumbnail(2715, 3510, 400, 250, 4, "layers_osrs/mapsquares/-1/2")
	if err != nil {
		panic(err)
	}
	output, err := os.Create("test-thumbnail.png")
	if err != nil {
		panic(err)
	}
	if err = png.Encode(output, thumbnail); err != nil {
		panic(err)
	}
}
