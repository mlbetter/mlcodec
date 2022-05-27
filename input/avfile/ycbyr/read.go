package ycbyr

import (
	"bufio"
	"image"
	"log"
	"os"
)

type YCbCr struct {
	filePath   string
	format     image.YCbCrSubsampleRatio
	weight     int
	height     int
	frameCnt   int
	frameIndex int
	frames     []*image.YCbCr
}

func NewReader(filePath string, w, h int, format image.YCbCrSubsampleRatio) *YCbCr {
	return &YCbCr{
		filePath:   filePath,
		weight:     w,
		height:     h,
		format:     format,
		frameCnt:   0,
		frameIndex: 0,
	}
}

func (y *YCbCr) Read() error {
	f, err := os.Open(y.filePath)
	if err != nil {
		log.Panicln(err)
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	log.Println(r.Size())

	// todo: read buffer by video fmt
	return nil
}

func (y *YCbCr) GetOneFrame() *image.YCbCr {
	return nil
}

func (y *YCbCr) GetFrameCnt() int {
	return 0
}
