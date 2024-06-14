package tool

import (
	"image"
	"image/png"
	"io"
	"log"
	"os"

	gzip "github.com/klauspost/pgzip"
)

func Gzip() {
	fr, err := os.Open("/Users/tian/projects/hc-flutterflow/editor/assets/module/新建空压机 (1).json")
	if err != nil {
		log.Fatalf("failed to open file to read: %v", err)
	}
	defer fr.Close()

	fw, err := os.Create("/Users/tian/projects/hc-flutterflow/editor/assets/module/新建空压机.gz")
	if err != nil {
		log.Fatalf("failed to open file to write: %v", err)
	}
	defer fw.Close()

	w := gzip.NewWriter(fw)
	defer w.Close()
	// 1MB block with 4 concurrency
	w.SetConcurrency(1<<20, 4)
	_, err = io.Copy(w, fr)
	if err != nil {
		log.Fatalf("faield to gzip: %v", err)
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
