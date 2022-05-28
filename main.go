package main

import (
	"github/mlbetter/mlcodec/input/avfile/ycbyr"
	"image"
	"image/jpeg"
	"log"
	"net/http"
)

var (
	// codec      = flag.String("codec", "enc/dec", "encode or decode")
	// sourcePath = flag.String("srcPath", "/source/file/path", "soure file path")
	// targetPath = flag.String("targetPath", "/target/file/path", "target file path")
	reader *ycbyr.YCbCr
)

func main() {
	// flag.Parse()

	log.SetFlags(log.Llongfile | log.Lmicroseconds)

	// log.Printf("mlcode start! codec:%s source file:%s target path:%s", *codec, *sourcePath, *targetPath)

	// if *codec != "enc" && *codec != "dec" {
	// 	log.Println("codec param error!")
	// 	return
	// }

	// if *sourcePath == "" || *targetPath == "" {
	// 	log.Println("source file path error or target path error")
	// 	return
	// }

	reader = ycbyr.NewReader("video/park_joy_420_720p50.yuv", 1280, 720, image.YCbCrSubsampleRatio(image.YCbCrSubsampleRatio420))
	if err := reader.Read(); err != nil {
		log.Println(err)
		return
	}

	http.HandleFunc("/show/pic", showPic)
	err := http.ListenAndServe(":9090", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func showPic(w http.ResponseWriter, r *http.Request) {
	i := reader.GetOneFrame()
	jpeg.Encode(w, i, &jpeg.Options{100})
}
