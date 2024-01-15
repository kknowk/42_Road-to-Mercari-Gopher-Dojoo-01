// package imgconv_test
package imgconv_test

import (
	"convert/imgconv"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"fmt"
)

// imagePath は通常の画像ファイルのパスを生成します。
func imagePath(t *testing.T, filename string, expectedErr bool, encodeErr bool) string {
	t.Helper()
	if (!expectedErr && !encodeErr) || (expectedErr && encodeErr) {
		return "../testdata/images/" + filename
	} else {
		return "../testdata/error_images/" + filename
	}
}

// TestImageConversion is a test function
func TestImageConversion(t *testing.T) {
	test := []struct {
		name         string
		inputPath    string
		outputFormat string
		quality      int
		expectedErr  bool
		encodeErr    bool
	}{
		{"PNG to JPEG", "jpg_png_top.png", "jpg", 100, false, false},
		{"JPEG to PNG", "image/ナポレオン.jpg", "png", 100, false, false},
		{"change quality PNG to JPEG", "jpg_png_top.png", "jpg", 10, false, false},
		{"Corrupted Image File JPG", "not_Decode/not_Decode_ナポレオン.jpg", "png", 100, true, false},
		{"Corrupted Image File PNG", "not_Decode/not_Decode_jpg_png_top.png", "jpg", 100, true, false},
		{"Cannot Create File JPG", "not_Create/かけないナポレオン.jpg", "png", 100, true, false},
		{"Cannot Create File PNG", "not_Create/かけないjpg_png_top.png", "jpg", 100, true, false},
		{"Encode Failure JPG", "jpg_png_top.png", "jpg", 100, true, true},
		{"Encode Failure PNG", "image/ナポレオン.jpg", "png", 100, true, true},
	}

	for _, tt := range test {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockEncoder := &imgconv.MockEncoder{FailEncode: tt.encodeErr}
			ic := imgconv.ImageConverter{
				JpegEncoder: mockEncoder,
				PngEncoder:  mockEncoder,
				Quality:     tt.quality,
			}

			var err error

			tt.inputPath = imagePath(t, tt.inputPath, tt.expectedErr, tt.encodeErr)
			if tt.outputFormat == "jpg" {
				err = ic.ToJpg(tt.inputPath)
			} else if tt.outputFormat == "png" {
				err = ic.ToPng(tt.inputPath)
			}

			if tt.expectedErr {
				if err == nil {
					t.Errorf("%s: expected error, got none", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("%s: unexpected error: %s", tt.name, err)
				}
				// check if file exists
				outputPath := strings.Replace(tt.inputPath, filepath.Ext(tt.inputPath), "."+tt.outputFormat, -1)
				if _, err := os.Stat(outputPath); err != nil {
					t.Errorf("%s: expected file %s to exist, got none %s", tt.name, outputPath, err)
					fmt.Println(err)
				}
			}
		})
	}

}
