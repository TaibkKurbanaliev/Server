package server

import (
	"bytes"
	"image/png"
)

type WallPaper struct {
	Id              int64 `json:"id"`
	UserId          int64
	Image           []byte `json:"image"`
	Width, Height   int
	Format          string
	Name            string
	Decription      string
	Tag             string
	Rating          float32
	NumberOFRatings int64
}

type JsonImage struct {
	Image []byte `json:"image"`
}

func (wallPaper *WallPaper) Init(jsonImage JsonImage) error {

	imageInfo, err := png.Decode(bytes.NewBuffer(jsonImage.Image))

	if err != nil {
		return err
	}

	wallPaper.Id = 1
	wallPaper.UserId = 1
	wallPaper.Width = imageInfo.Bounds().Dx()
	wallPaper.Height = imageInfo.Bounds().Dy()
	wallPaper.Format = "png"
	wallPaper.Image = jsonImage.Image
	wallPaper.Name = ""
	wallPaper.Decription = ""
	wallPaper.Tag = ""
	wallPaper.Rating = 0.0
	wallPaper.NumberOFRatings = 0

	return nil
}
