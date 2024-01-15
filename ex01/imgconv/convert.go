// godoc -http=:6060
// this package converts png to jpg and jpg to png
package imgconv

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"
	"errors"
)

// ImageEncoder インターフェイスは、JPEG画像のエンコード処理を抽象化します。
type ImageEncoder interface {
	Encode(w io.Writer, m image.Image, o *jpeg.Options) error
}

// PngEncoder インターフェイスは、PNG画像のエンコード処理を抽象化します。
type PngEncoder interface {
	EncodePNG(w io.Writer, m image.Image) error
}

type MockEncoder struct {
	// FailEncode が true の場合にのみエンコードを失敗させる
	FailEncode bool
}

func (m *MockEncoder) Encode(w io.Writer, img image.Image, o *jpeg.Options) error {
	if m.FailEncode {
		return errors.New("mock encode error")
	}
	return jpeg.Encode(w, img, o)
}

func (m *MockEncoder) EncodePNG(w io.Writer, img image.Image) error {
	if m.FailEncode {
		return errors.New("mock encode error")
	}
	return png.Encode(w, img)
}

// ImageConverter is a struct
type ImageConverter struct {
	Quality     int // jpeg quality
	JpegEncoder ImageEncoder
	PngEncoder  PngEncoder
}

// convert image to jpg
func (ic *ImageConverter) ToJpg(path string) error {
	// open file
	file, _ := os.Open(path)
	// main.goでos.ReadFile(path)で読み込んでいるので、ここには来ない
	// if err != nil {
	// 	return fmt.Errorf("error opening file: %w", err)
	// }
	defer file.Close()

	// decode image
	img, err := png.Decode(file)
	if err != nil {
		return fmt.Errorf("error decoding image: %w", err)
	}

	// create new file
	newFile, err := os.Create(strings.Replace(path, ".png", ".jpg", -1))
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer newFile.Close()

	// encode image
	err = ic.JpegEncoder.Encode(newFile, img, &jpeg.Options{Quality: ic.Quality})
	if err != nil {
		return fmt.Errorf("error encoding image: %w", err)
	}

	return nil
}

// convert image to png
func (ic *ImageConverter) ToPng(path string) error {
	// open file
	file, _ := os.Open(path)
	// main.goでos.ReadFile(path)で読み込んでいるので、ここには来ない
	// if err != nil {
	// 	return fmt.Errorf("error opening file: %w", err)
	// }
	defer file.Close()

	// decode image
	img, err := jpeg.Decode(file)
	if err != nil {
		return fmt.Errorf("error decoding image: %w", err)
	}

	// create new file
	newFile, err := os.Create(strings.Replace(path, ".jpg", ".png", -1))
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer newFile.Close()


	// encode image
	err = ic.PngEncoder.EncodePNG(newFile, img)
	if err != nil {
		return fmt.Errorf("error encoding image: %w", err)
	}

	return nil
}
