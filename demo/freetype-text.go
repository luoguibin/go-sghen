package demo

import (
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/golang/freetype"
)

func main() {
	// 装载字体
	fontSourceBytes, err := ioutil.ReadFile("file/font/DejaVuSans.ttf")
	if err != nil {
		panic(err)
	}
	trueTypeFont, err := freetype.ParseFont(fontSourceBytes)
	if err != nil {
		panic(err)
	}

	// 新建需要保存图片的文件
	file, err := os.Create("file/capture/temp.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 新建一个RGBA位图
	img := image.NewNRGBA(image.Rect(0, 0, 256, 256))

	fc := freetype.NewContext()
	fc.SetDPI(72)
	fc.SetFont(trueTypeFont)
	fc.SetFontSize(28)
	fc.SetClip(img.Bounds())
	fc.SetDst(img)

	fc.SetSrc(image.NewUniform(color.RGBA{33, 200, 233, 255}))
	pt := freetype.Pt(0, 50)
	_, err = fc.DrawString("hello word", pt)
	if err != nil {
		panic(err)
	}

	fc.SetSrc(image.NewUniform(color.RGBA{200, 33, 233, 255}))
	pt = freetype.Pt(0, 100)
	_, err = fc.DrawString("hello word", pt)
	if err != nil {
		panic(err)
	}

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}
