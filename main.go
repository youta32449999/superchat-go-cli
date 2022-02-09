package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"strconv"
	"strings"
)

//go:embed "template/midori.png"
var templateBytes []byte

//go:embed "template/icon192.png"
var iconBytes []byte

//go:embed "font/Koruri/Koruri-Regular.ttf"
var regularFontBytes []byte

//go:embed "font/Koruri/Koruri-Semibold.ttf"
var semiBoldFontBytes []byte

func main() {
	flag.Parse()
	args := flag.Args()
	var name  = args[0]
	var amount, err = strconv.Atoi(args[1])
	if err != nil {
		log.Fatal(err)
	}
	var message = args[2]

	templateImage, _, err := image.Decode(bytes.NewReader(templateBytes))
	if err != nil {
		log.Fatal(err)
	}

	iconImage, _, err := image.Decode(bytes.NewReader(iconBytes))
	if err != nil {
		log.Fatal(err)
	}

	// 元画像をリサイズ
	resized := image.NewRGBA(image.Rect(0, 0, iconImage.Bounds().Size().X*3/5, iconImage.Bounds().Size().Y*3/5))
	draw.CatmullRom.Scale(resized, resized.Bounds(), iconImage, iconImage.Bounds(), draw.Over, nil)

	//オリジナル画像上のどこからlogoイメージを重ねるか
	startPointLogo := image.Point{45, 30}

	logoRectangle := image.Rectangle{startPointLogo, startPointLogo.Add(resized.Bounds().Size())}
	originRectangle := image.Rectangle{image.Point{0, 0}, templateImage.Bounds().Size()}

	// 画像の合成
	rgba := image.NewRGBA(originRectangle)
	draw.Draw(rgba, logoRectangle, resized, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, originRectangle, templateImage, image.Point{0, 0}, draw.Over)

	drawName(rgba, name)
	drawAmount(rgba, amount)
	drawText(rgba, message)
	outputImage(rgba)

	fmt.Println("画像の出力が完了しました")
}

func drawName(img draw.Image, name string) {
	// フォントファイルを読み込み
	ft, err := truetype.Parse(regularFontBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	opt := truetype.Options{
		Size:              34,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}

	text := name

	face := truetype.NewFace(ft, &opt)

	dr := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	dr.Dot.X = 190 * 64
	dr.Dot.Y = fixed.I(65)

	dr.DrawString(text)
}

func drawAmount(img draw.Image, amount int) {
	// フォントファイルを読み込み
	ft, err := truetype.Parse(semiBoldFontBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	opt := truetype.Options{
		Size:              38,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}

	text := convert(amount)

	face := truetype.NewFace(ft, &opt)

	dr := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	dr.Dot.X = 220 * 64
	dr.Dot.Y = fixed.I(115)

	dr.DrawString(text)
}

func drawText(img draw.Image, text string) {
	ft, err := truetype.Parse(semiBoldFontBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	opt := truetype.Options{
		Size:              34,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}

	face := truetype.NewFace(ft, &opt)

	dr := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	dr.Dot.X = 50 * 64
	dr.Dot.Y = fixed.I(205)

	dr.DrawString(text)
}

func outputImage(img draw.Image){
	const outputFile = "./spacha.png"
	// 画像の出力
	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err)
	}
	defer out.Close()

	png.Encode(out, img)
}

func convert(integer int) string {
	arr := strings.Split(fmt.Sprintf("%d", integer), "")
	cnt := len(arr) - 1
	res := ""
	i2 := 0
	for i := cnt; i >= 0; i-- {
		if i2 > 2 && i2%3 == 0 {
			res = fmt.Sprintf(",%s", res)
		}
		res = fmt.Sprintf("%s%s", arr[i], res)
		i2++
	}
	return res
}
