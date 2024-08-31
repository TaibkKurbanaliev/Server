package server

import (
	"bytes"
	"image/png"
	"log"
)

type WallPaper struct {
	Id              int64   `json:"id"`
	UserId          int64   `json:"userID"`
	ImagePath       string  `json:"image"`
	Width           int     `json:"width"`
	Height          int     `json:"height"`
	Format          string  `json:"format"`
	Name            string  `json:"name"`
	Decription      string  `json:"description"`
	Tag             string  `json:"tag"`
	Rating          float32 `json:"rating"`
	NumberOfRatings int64   `json:"numberOfRatings"`
}

type JsonImage struct {
	FileName string `json:"name"`
	Image    []byte `json:"image"`
}

func NewWallPaper(jsonImage JsonImage, storagePath string) (*WallPaper, error) {

	var wallPaper *WallPaper = new(WallPaper)
	imageInfo, err := png.Decode(bytes.NewBuffer(jsonImage.Image))

	if err != nil {
		log.Panic(err)
		return nil, err
	}

	wallPaper.Id = 1
	wallPaper.UserId = 1
	wallPaper.Width = imageInfo.Bounds().Dx()
	wallPaper.Height = imageInfo.Bounds().Dy()
	wallPaper.Format = "png"
	wallPaper.ImagePath = storagePath + jsonImage.FileName
	wallPaper.Name = jsonImage.FileName
	wallPaper.Decription = ""
	wallPaper.Tag = ""
	wallPaper.Rating = 0.0
	wallPaper.NumberOfRatings = 0

	return wallPaper, nil
}
