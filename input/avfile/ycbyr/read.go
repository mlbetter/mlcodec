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

// readYUV420 convert byte to image.ycbyr
// adout yuv420 format https://blog.csdn.net/mzpmzk/article/details/81239532
func (y *YCbCr) readYUV420(r *bufio.Reader) error {

	var (
		frameSize = y.weight * y.height * 3 / 2
		yBegin    = 0
		ySize     = y.weight * y.height
		cbBegin   = yBegin + ySize
		cbSize    = y.weight * y.height / 4
		crBegin   = cbBegin + cbSize
		bufIter   = make([]byte, frameSize+41)
	)

	index := 0
	for {
		n, err := r.Read(bufIter)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			log.Println("read file done")
			break
		}
		if index >= 100 {
			return nil
		}
		if index == 0 {
			bufIter = bufIter[41:]
		} else {
			bufIter = bufIter[6:]
		}

		imageIter := image.NewYCbCr(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: y.weight, Y: y.height},
		}, image.YCbCrSubsampleRatio420)

		// y
		for i := yBegin; i < cbBegin; i++ {
			imageIter.Y[i] = BytesToUINT8(bufIter[i : i+1])
		}
		// uv
		j := 0
		for i := cbBegin; i < crBegin; i++ {
			imageIter.Cb[j] = BytesToUINT8(bufIter[i : i+1])
			j++
		}
		j = 0
		// v
		for i := crBegin; i < frameSize; i++ {
			imageIter.Cr[j] = BytesToUINT8(bufIter[i : i+1])
			j++
		}
		y.frames = append(y.frames, imageIter)
		log.Println(index, " done...")
		index++
		bufIter = make([]byte, frameSize+6)
		y.frameCnt++
	}
	return nil
}

func BytesToUINT8(bys []byte) uint8 {
	bytebuff := bytes.NewBuffer(bys)
	var data uint8
	binary.Read(bytebuff, binary.BigEndian, &data)
	return data
}
