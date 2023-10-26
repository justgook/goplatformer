package resources

import (
	"bytes"
	"encoding/gob"
	"image"
	"image/png"

	"github.com/justgook/goplatformer/pkg/util"
)

func Load(data []byte, a any) error {
	var outputBuffer bytes.Buffer
	outputBuffer.Write(data)

	if err := gob.NewDecoder(&outputBuffer).Decode(a); err != nil {
		return err
	}

	return nil
}

func Save(a any) ([]byte, error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)

	if err := encoder.Encode(a); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func BytesToImage(input []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, util.Catch(err)
	}
	return img, nil
}

func ImageToBytes(input image.Image) ([]byte, error) {
	var b bytes.Buffer
	if err := png.Encode(&b, input); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
