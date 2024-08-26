package server

import (
	"bytes"
	"image"
)

type WallPaper struct {
	id              int64  `json:"id"`
	image           []byte `json:"image"`
	Width, Height   int
	Format          string
	Name            string
	Decription      string
	Tag             string
	Rating          float32
	NumberOFRatings int64
}

func NewWallPaper(img []byte) (*WallPaper, error) {
	imageInfo, format, err := image.Decode(bytes.NewBuffer(img))

	if err != nil {
		return nil, err
	}

	var wallPaper *WallPaper = new(WallPaper)

	wallPaper.Width = imageInfo.Bounds().Dx()
	wallPaper.Height = imageInfo.Bounds().Dy()
	wallPaper.Format = format
	wallPaper.image = img

	return wallPaper, nil
}
