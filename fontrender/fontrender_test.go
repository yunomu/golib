package fontrender

import (
	"testing"

	"image"

	"github.com/golang/freetype"
)

func TestRenderer(t *testing.T) {
	ttf, err := freetype.ParseFont(fontData)
	if err != nil {
		t.Fatalf(err.Error())
	}

	r := NewRenderer(ttf, image.Rect(0, 0, 300, 300), *content)
	var _ = r
}
