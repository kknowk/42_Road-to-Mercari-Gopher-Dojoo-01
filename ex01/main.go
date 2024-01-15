// godoc -http=:6060
package main

import (
	"convert/imgconv"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// this is main function
func main() {
	if len(os.Args) < 2 {
		fmt.Println("error: invaild argument")
		return
	}

	// default format
	inputFormat := flag.String("i", "jpg", "input format")
	outputFormat := flag.String("o", "png", "input format")
	qualityFormat := flag.Int("q", 100, "jpeg quality")

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("error: invaild argument")
		fmt.Println("usage: ./convert -i=inputFormat -o=outputFormat <directory>")
		return
	}

	dir := flag.Arg(0)
	ic := imgconv.ImageConverter{
		Quality:     *qualityFormat,			   // JPEGのクオリティ
		JpegEncoder: &imgconv.MockEncoder{}, // JPEGエンコーダーの実装
		PngEncoder:  &imgconv.MockEncoder{},   // PNGエンコーダーの実装
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ディレクトリの場合はスキップ
		if info.IsDir() {
			return nil
		}

		fbyte, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}
		ftype := http.DetectContentType(fbyte)

		switch {
		case ftype == "image/png" && *inputFormat == "png" && *outputFormat == "jpg":
			ic.ToJpg(path)
		case ftype == "image/jpeg" && *inputFormat == "jpg" && *outputFormat == "png":
			ic.ToPng(path)
		case ftype != "image/png" && ftype != "image/jpeg":
			fmt.Printf("error: %s is not a valid file\n", path)
		}

		return nil
	})
	if err != nil {
		// fmt.Printf("error: %s: no such file or directory", dir)
		fmt.Printf("error: %s\n", err)
	}
}
