package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"zzz_helper/internal/utils/file2"
)

// binarizeImage performs binary thresholding on the image with the specified threshold
func binarizeImage(img image.Image, threshold uint8) *image.Gray {
	bounds := img.Bounds()
	grayImage := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Convert RGB to grayscale intensity
			gray := uint8((r*299 + g*587 + b*114) / 1000 >> 8)

			// Apply threshold
			if gray > threshold {
				grayImage.Set(x, y, color.Gray{255}) // White
			} else {
				grayImage.Set(x, y, color.Gray{0}) // Black
			}
		}
	}
	return grayImage
}

func binarizeImageWithBytes(src []byte, threshold uint8) ([]byte, error) {
	reader := bytes.NewReader(src)
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	// Perform binarization
	binaryImage := binarizeImage(img, threshold)

	out := bytes.Buffer{}
	err = png.Encode(&out, binaryImage)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil

}

func main() {
	inFile := "test\\drivers\\Snipaste_2025-01-18_04-21-52.png"
	outFile := "tmpout/output.png"
	src, err := file2.ReadFileBytes(inFile)
	if err != nil {
		panic(err)
	}
	dst, err := binarizeImageWithBytes(src, 128)
	if err != nil {
		panic(err)
	}
	dst, err = binarizeImageWithBytes(src, 128)

	file2.WriteFile(outFile, dst)

}
