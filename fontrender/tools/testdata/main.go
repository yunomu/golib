package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

var (
	fontFile = flag.String("font", "./testdata/NotoSans-Black.ttf", "ttf path")
	outFile  = flag.String("o", "data_test.go", "")
)

const srcTemplate = `package fontrender

var fontData = []byte{
	{{.}}
}
`

func init() {
	flag.Parse()
}

func main() {
	bs, err := ioutil.ReadFile(*fontFile)
	if err != nil {
		log.Fatalln(err)
	}

	var ss []string
	for _, b := range bs {
		ss = append(ss, fmt.Sprintf("0x%02x,", b))
	}

	t, err := template.New("fontdata").Parse(srcTemplate)
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.Create(*outFile)
	if err != nil {
		log.Fatalln(err)
	}

	if err := t.Execute(f, strings.Join(ss, " ")); err != nil {
		log.Fatalln(err)
	}
}
