package main

import (
	"flag"
	"image"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/disintegration/imaging"

	"github.com/yunomu/golib/fontrender"
)

var (
	fontFile = flag.String("font", "./testdata/NotoSans-Black.ttf", "ttf path")
	content  = flag.String("content", "Hello", "")
	srcFile  = flag.String("src", "", "Source image file")

	outFile = flag.String("o", "out.png", "")
)

func init() {
	flag.Parse()
	log.SetOutput(os.Stderr)
}

func main() {
	ttf, err := fontrender.LoadFont(*fontFile)
	if err != nil {
		log.Fatalln(err)
	}

	var in io.Reader
	switch *srcFile {
	case "":
		log.Fatalln("-src is required")
	case "-":
		in = os.Stdin
	default:
		f, err := os.Open(*srcFile)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		in = f
	}
	srcImg, err := png.Decode(in)
	if err != nil {
		log.Fatalln(err)
	}

	img, err := fontrender.Render(ttf, srcImg.Bounds(), *content)
	if err != nil {
		log.Fatalln(err)
	}

	dstImg := imaging.Overlay(srcImg, img, image.ZP, 1.0)

	var out io.Writer = os.Stdout
	if *outFile != "-" {
		f, err := os.Create(*outFile)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		out = f
	}

	if err := png.Encode(out, dstImg); err != nil {
		log.Fatalln(err)
	}
}
