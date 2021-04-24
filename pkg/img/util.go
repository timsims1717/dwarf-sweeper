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
