package tool

import (
	"encoding/base64"
	"fmt"
	gzip "github.com/klauspost/pgzip"
	"image"
	"image/png"
	"io"
	"log"
	"os"
	"strings"
)

func DeGzip() {

	open, err := os.Open("/Users/tian/projects/hc-flutterflow/editor/assets/module/安装时间.gz")
	if err != nil {
		return
	}
	reader, err := gzip.NewReader(open)
	if err != nil {
		return
	}
	create, err := os.Create("/Users/tian/Downloads/安装时间2.json")
	if err != nil {
		return
	}

	_, err = io.Copy(create, reader)
	if err != nil {
		return
	}
	fmt.Println("Copy Ok")

}

func Base64() {

	// 打开输入文件
	inputFile := "/Users/tian/projects/hc-flutterflow/editor/assets/module/test.gz"
	inFile, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Failed to open input file: %v", err)
	}
	defer inFile.Close()

	// 创建输出文件
	outputFile := "/Users/tian/projects/hc-flutterflow/editor/assets/module/test.txt"
	outFile, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	// 创建Base64编码器
	encoder := base64.NewEncoder(base64.StdEncoding, outFile)
	defer encoder.Close()

	// 将输入文件内容复制到编码器（最终写入到输出文件）
	_, err = io.Copy(encoder, inFile)
	if err != nil {
		log.Fatalf("Failed to encode input file: %v", err)
	}

	log.Println("Base64 encoded data has been written to", outputFile)

}

func Gzip() {

	dirs, err := os.ReadDir("/Users/tian/projects/hc-flutterflow/editor/assets/module/")
	if err != nil {
		return
	}

	for _, dir := range dirs {

		//info, err := dir.Info()
		//if err != nil {
		//	return
		//}
		//
		//return
		//fmt.Println(info)

		if dir.Name() == "config.json" {
			continue
		}
		if !strings.HasSuffix(dir.Name(), ".json") {
			continue
		}

		fr, err := os.Open("/Users/tian/projects/hc-flutterflow/editor/assets/module/" + dir.Name())
		if err != nil {
			log.Fatalf("failed to open file to read: %v", err)
		}
		defer fr.Close()

		baseName := strings.TrimSuffix(dir.Name(), ".json")

		fw, err := os.Create("/Users/tian/projects/hc-flutterflow/editor/assets/module/" + baseName + ".gz")
		if err != nil {
			log.Fatalf("failed to open file to write: %v", err)
		}
		defer fw.Close()

		w := gzip.NewWriter(fw)
		defer func(w *gzip.Writer) {
			err := w.Close()
			if err != nil {
				if err != nil {
					//手机登录微信了？
				}
			}
		}(w)
		// 1MB block with 4 concurrency
		w.SetConcurrency(1<<20, 4)
		_, err = io.Copy(w, fr)
		if err != nil {
			log.Fatalf("faield to gzip: %v", err)
		}

	}

}

func CompressJPEG(inputFile, outputFile string, quality int) error {
	// 打开输入文件
	inFile, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer inFile.Close()

	// 解码图像
	img, _, err := image.Decode(inFile)
	if err != nil {
		return err
	}

	// 创建输出文件
	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 设置 PNG 编码器的压缩级别
	encoder := png.Encoder{
		CompressionLevel: 80,
	}

	// 将图像编码并写入输出文件
	return encoder.Encode(outFile, img)
}
