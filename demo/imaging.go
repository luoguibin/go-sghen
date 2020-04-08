package demo

import (
	"fmt"
	"path"
	"strings"

	"github.com/disintegration/imaging"
)

func main() {
	fullFilename := "img/P91101-223114.jpg"
	filenameWithSuffix := path.Base(fullFilename)
	fileSuffix := path.Ext(filenameWithSuffix)
	if fileSuffix != ".png" && fileSuffix != ".jpeg" && fileSuffix != ".jpg" {
		fmt.Println("only support .png .jpeg .jpg")
		return
	}

	filename := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	fmt.Println(filename, fileSuffix)
	// test	.png

	// Open a test image.
	srcImage, err := imaging.Open(fullFilename)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Lanczos 				- A high-quality resampling filter for photographic images yielding sharp results.
	// CatmullRom			- A sharp cubic filter that is faster than Lanczos filter while providing similar results.
	// MitchellNetravali 	- A cubic filter that produces smoother results with less ringing artifacts than CatmullRom.
	// Linear 				- Bilinear resampling filter, produces smooth output. Faster than cubic filters.
	// Box 					- Simple and fast averaging filter appropriate for downscaling. When upscaling it's similar to NearestNeighbor.
	// NearestNeighbor 		- Fastest resampling filter, no antialiasing.
	dstWith := 128
	dstImage128 := imaging.Resize(srcImage, dstWith, 0, imaging.Lanczos)

	err = imaging.Save(dstImage128, fmt.Sprintf("img/%s_%d%s", filename, dstWith, fileSuffix))
	if err != nil {
		fmt.Println(err)
	}
}
