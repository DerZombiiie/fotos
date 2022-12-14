package textimage

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"io/ioutil"
	"sync"
)

var f *truetype.Font
var fMu sync.RWMutex

// most systems have this one:
const fontPath = "/usr/share/fonts/liberation/LiberationSerif-Regular.ttf"

// UpdateFont sets the font used
func UpdateFont(b []byte) error {
	fMu.Lock()
	defer fMu.Unlock()

	var err error

	f, err = truetype.Parse(b)
	if err != nil {
		return err
	}

	return nil
}

// Generate Image creats an image thats centerd (no newline support)
func GenerateImage(width, height int, text string) *image.RGBA {
	if f == nil {
		b, err := ioutil.ReadFile(fontPath)
		if err != nil {
			return nil
		}

		UpdateFont(b)
	}

	fMu.RLock()
	defer fMu.RUnlock()

	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(125)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.NewUniform(color.NRGBA{0, 255, 0, 128}))
	c.SetHinting(font.HintingNone)

	// Truetype stuff
	opts := truetype.Options{
		Size: 125.0,
	}
	face := truetype.NewFace(f, &opts)

	advance, _ := StrAdvance(face, text)

	pt := freetype.Pt((width/2)-advance.Round()/2, height/2+int(125)/4)
	fmt.Printf("center: %d %d; w/h: %d/%d \n", pt.X.Round(), pt.Y.Round(), width, height)
	rgba.Set(pt.X.Round()-10, pt.Y.Round(), color.NRGBA{0, 255, 0, 128})

	_, err := c.DrawString(text, pt)
	if err != nil {
		fmt.Errorf("Error drawing string: %s \n", err)
	}

	return rgba
}

// like GlyphAdvance but for entire strings
// if ok != true one glyph isn't in font
func StrAdvance(face font.Face, str string) (i fixed.Int26_6, ok bool) {
	for _, x := range str {
		awidth, k := face.GlyphAdvance(x)
		if !k {
			ok = k
		}

		i += awidth - fixed.I(9)
	}

	return i, true
}
