package img

import (
	"github.com/faiface/pixel"
	"github.com/pkg/errors"
	"image"
	_ "image/png"
	"os"
)

func LoadImage(path string) (pixel.Picture, error) {
	errMsg := "load image"
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	return pixel.PictureDataFromImage(img), nil
}

func Reverse(o []*pixel.Sprite) []*pixel.Sprite {
	var s []*pixel.Sprite
	for i := len(o) - 1; i >= 0; i-- {
		s = append(s, o[i])
	}
	return s
}
