package resources

import (
	"bytes"
	"encoding/gob"
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

