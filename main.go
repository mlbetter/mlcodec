package main

import (
	"flag"
	"log"
)

var (
	codec      = flag.String("codec", "enc/dec", "encode or decode")
	sourcePath = flag.String("srcPath", "/source/file/path", "soure file path")
	targetPath = flag.String("targetPath", "/target/file/path", "target file path")
)

func main() {
	flag.Parse()

	log.SetFlags(log.Llongfile | log.Lmicroseconds)

	log.Printf("mlcode start! codec:%s source file:%s target path:%s", *codec, *sourcePath, *targetPath)

	if *codec != "enc" && *codec != "dec" {
		log.Println("codec param error!")
		return
	}

	if *sourcePath == "" || *targetPath == "" {
		log.Println("source file path error or target path error")
		return
	}

	return
}
