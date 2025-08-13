package lib

import (
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/png"
	"os"
)

var (
	MapFilePath   = GetEnv("MAP_FILE_PATH", "layers_osrs/mapsquares/-1/2")
	MapTilePixels = GetEnvInt("MAP_TILE_PIXELS", 4)
)

func CreateThumbnail(x, y, width, height int) (*image.RGBA, error) {
	chunkX := (x >> 3) / 8
	chunkY := (y >> 3) / 8

	chunksImage := createCombinedChunksImage(chunkX, chunkY)
	croppedImage := createCroppedImage(x, y, chunkX, chunkY, width, height, chunksImage)
	drawCrossHair(croppedImage)

	scaled := image.NewRGBA(image.Rect(0, 0, croppedImage.Bounds().Max.X*2, croppedImage.Bounds().Max.Y*2))
	draw.NearestNeighbor.Scale(scaled, scaled.Rect, croppedImage, croppedImage.Bounds(), draw.Over, nil)

	return scaled, nil
}

func drawCrossHair(img *image.RGBA) {
	imageWidth := img.Bounds().Dx()
	imageHeight := img.Bounds().Dy()
	crossHairSize := 24
	crossHairColor := color.RGBA{R: 255, A: 255}

	x := (imageWidth / 2) - (crossHairSize / 2)
	y := (imageHeight / 2) - (crossHairSize / 2)
	for pos := 0; pos <= crossHairSize; pos++ {
		img.SetRGBA(x, y, crossHairColor)
		img.SetRGBA(x, y+1, crossHairColor)
		img.SetRGBA(x, y-1, crossHairColor)
		x++
		y++
	}

	x = (imageWidth / 2) - (crossHairSize / 2)
	y = (imageHeight / 2) + (crossHairSize / 2)
	for pos := 0; pos <= crossHairSize; pos++ {
		img.SetRGBA(x, y, crossHairColor)
		img.SetRGBA(x, y+1, crossHairColor)
		img.SetRGBA(x, y-1, crossHairColor)
		x++
		y--
	}
}

func createCroppedImage(x, y, chunkX, chunkY, width, height int, img *image.RGBA) *image.RGBA {
	minX := ((chunkX - 1) << 3) * 8
	minY := ((chunkY - 1) << 3) * 8

	mapHeight := img.Bounds().Dy()

	mapX := ((x - minX) * MapTilePixels) - MapTilePixels
	mapY := (mapHeight - (y-minY)*MapTilePixels) - MapTilePixels

	imageCoords := image.Rect(mapX-(width/2), mapY+(height/2), mapX+(width/2), mapY-(height/2))
	subImage := img.SubImage(imageCoords)

	baseImage := image.NewRGBA(image.Rect(0, 0, subImage.Bounds().Dx(), subImage.Bounds().Dy()))
	draw.Draw(baseImage, baseImage.Bounds(), subImage, subImage.Bounds().Min, draw.Src)

	return baseImage
}

func createCombinedChunksImage(chunkX, chunkY int) *image.RGBA {
	chunkImageSize := 64 * MapTilePixels
	surroundingChunks := 1
	chunks := 1 + (surroundingChunks * 2)

	baseImage := image.NewRGBA(image.Rect(0, 0, chunks*chunkImageSize, chunks*chunkImageSize))

	for imageChunkX := -surroundingChunks; imageChunkX <= surroundingChunks; imageChunkX++ {
		for imageChunkY := surroundingChunks; imageChunkY >= -surroundingChunks; imageChunkY-- {
			chunkImage, err := loadChunkFile(chunkX+imageChunkX, chunkY-imageChunkY)
			if err != nil {
				continue
			}
			x := (imageChunkX + 1) * chunkImageSize
			y := (imageChunkY + 1) * chunkImageSize
			bounds := image.Rectangle{
				Min: image.Point{
					X: x,
					Y: y,
				},
				Max: image.Point{
					X: x + chunkImageSize,
					Y: y + chunkImageSize,
				},
			}
			img := *chunkImage
			draw.Draw(baseImage, bounds, *chunkImage, img.Bounds().Min, draw.Src)
		}
	}
	return baseImage
}

func loadChunkFile(x, y int) (*image.Image, error) {
	imageFile, err := os.Open(fmt.Sprintf("%s/0_%d_%d.png", MapFilePath, x, y))
	if err != nil {
		return nil, err
	}

	defer func(imageFile *os.File) {
		err := imageFile.Close()
		if err != nil {
			panic(err)
		}
	}(imageFile)

	img, err := png.Decode(imageFile)
	if err != nil {
		return nil, err
	}
	return &img, nil
}
