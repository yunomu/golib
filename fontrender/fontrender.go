package fontrender

import (
	"errors"
	"image"
	"image/color"
	"io/ioutil"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

var (
	ErrOutOfBounds = errors.New("out of bounds")
)

func LoadFont(path string) (*truetype.Font, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return freetype.ParseFont(bs)
}

type options struct {
	color           color.Color
	fontSize        float64
	dpi             float64
	allowOutOfBound bool
	hinting         font.Hinting
	pos             image.Point
}

type Option func(*options)

func SetFontSize(size float64) Option {
	return func(r *options) {
		r.fontSize = size
	}
}

func SetDPI(dpi float64) Option {
	return func(r *options) {
		r.dpi = dpi
	}
}

func AllowOutOfBound(f bool) Option {
	return func(r *options) {
		r.allowOutOfBound = f
	}
}

func SetHinting(hinting font.Hinting) Option {
	return func(r *options) {
		r.hinting = hinting
	}
}

func Render(
	ttf *truetype.Font,
	clip image.Rectangle,
	content string,
	opts ...Option,
) (*image.NRGBA, error) {
	options := &options{
		color:    color.Black,
		fontSize: 12,
		dpi:      72,
		pos:      image.ZP,
	}
	for _, f := range opts {
		f(options)
	}

	dst := image.NewNRGBA(clip)

	fctx := freetype.NewContext()
	fctx.SetFont(ttf)
	fctx.SetFontSize(options.fontSize)
	fctx.SetDPI(options.dpi)
	fctx.SetHinting(options.hinting)
	fctx.SetClip(clip)
	fctx.SetSrc(image.NewUniform(options.color))
	fctx.SetDst(dst)

	extent, err := fctx.DrawString(content,
		fixed.Point26_6{
			X: fctx.PointToFixed(float64(options.pos.X)),
			Y: fctx.PointToFixed(float64(options.pos.Y) + options.fontSize),
		},
	)
	if err != nil {
		return nil, err
	}
	if !options.allowOutOfBound && (clip.Max.X < options.pos.X+extent.X.Round() || clip.Max.Y < options.pos.Y+extent.Y.Round()) {
		return nil, ErrOutOfBounds
	}

	return dst, nil
}
