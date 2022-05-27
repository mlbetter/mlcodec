package ycbyr

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"io"
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
	switch y.format {
	case image.YCbCrSubsampleRatio420:
		return y.readYUV420(r)
	default:
		log.Println("not surpport now")
		return fmt.Errorf("not surrport now")
	}
}

func (y *YCbCr) GetOneFrame() *image.YCbCr {
	f := y.frames[y.frameIndex]
	y.frameIndex++
	if y.frameIndex >= y.frameCnt {
		y.frameIndex = 0
	}
	return f
}

func (y *YCbCr) GetFrameCnt() int {
	return y.frameCnt
}

func (y *YCbCr) readYUV420(r *bufio.Reader) error {
	bufIter := make([]byte, y.weight*y.height*3/2*8)

	for {
		n, err := r.Read(bufIter)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			log.Println("read file done")
			break
		}
		if n != y.weight*y.height*2 {
			log.Println("read length error, len:", n)
			return fmt.Errorf("read file error")
		}

		imageIter := image.NewYCbCr(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: y.weight, Y: y.height},
		}, image.YCbCrSubsampleRatio420)

		i := 0
		// y
		for i < y.weight*y.height*8 {
			imageIter.Y[i] = BytesToUINT8(bufIter[i : i+8])
			i = i + 8
		}

		// u
		for i < y.weight*y.height/4*8 {
			imageIter.Cb[i] = BytesToUINT8(bufIter[i : i+8])
			i = i + 8
		}

		// v
		for i < y.weight*y.height/4*8 {
			imageIter.Cr[i] = BytesToUINT8(bufIter[i : i+8])
			i = i + 8
		}

		y.frames = append(y.frames, imageIter)
	}
	return nil
}

func BytesToUINT8(bys []byte) uint8 {
	bytebuff := bytes.NewBuffer(bys)
	var data uint8
	binary.Read(bytebuff, binary.BigEndian, &data)
	return data
}
