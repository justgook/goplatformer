package qoi

/*

QOI - The “Quite OK Image” format for fast, lossless image compression

Original version by Dominic Szablewski - https://phoboslab.org
Go version by Xavier-Frédéric Moulet

*/

import (
	"bufio"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"io"

	"github.com/justgook/goplatformer/pkg/util"
)

const (
	qoiIndex byte = 0b00_000000
	qoiDiff  byte = 0b01_000000
	qoiLuma  byte = 0b10_000000
	qoiRun   byte = 0b11_000000
	qoiRgb   byte = 0b1111_1110
	qoiRgba  byte = 0b1111_1111

	qoiMask2 byte = 0b11_000000
)

const qoiMagic = "qoif"

const qoiPixelsMax = 400_000_000 // 400 million pixels ought to be enough for anybody

func qoi_COLOR_HASH(r, g, b, a byte) byte {
	return byte(r*3 + g*5 + b*7 + a*11)
}

type pixel [4]byte

func Decode(r io.Reader) (image.Image, error) {
	cfg, err := DecodeConfig(r)
	if err != nil {
		return nil, err
	}
	NBPixels := cfg.Width * cfg.Height
	if NBPixels == 0 || NBPixels > qoiPixelsMax {
		return nil, errors.New("bad image dimensions")
	}

	b := bufio.NewReader(r)

	img := image.NewNRGBA(image.Rect(0, 0, cfg.Width, cfg.Height))

	var index [64]pixel

	run := 0

	pixels := img.Pix // pixels yet to write
	px := pixel{0, 0, 0, 255}
	for len(pixels) > 0 {
		if run > 0 {
			run--
		} else {
			b1, err := b.ReadByte()
			if err == io.EOF {
				return img, nil
			}
			if err != nil {
				return nil, err
			}

			switch {
			case b1 == qoiRgb:
				_, err = io.ReadFull(b, px[:3])
				if err != nil {
					return nil, err
				}
			case b1 == qoiRgba:
				_, err = io.ReadFull(b, px[:])
				if err != nil {
					return nil, err
				}
			case b1&qoiMask2 == qoiIndex:
				px = index[b1]
			case b1&qoiMask2 == qoiDiff:
				px[0] += ((b1 >> 4) & 0x03) - 2
				px[1] += ((b1 >> 2) & 0x03) - 2
				px[2] += (b1 & 0x03) - 2
			case b1&qoiMask2 == qoiLuma:
				b2, err := b.ReadByte()
				if err != nil {
					return nil, err
				}
				vg := (b1 & 0b00111111) - 32
				px[0] += vg - 8 + ((b2 >> 4) & 0x0f)
				px[1] += vg
				px[2] += vg - 8 + (b2 & 0x0f)
			case b1&qoiMask2 == qoiRun:
				run = int(b1 & 0b00111111)
			default:
				px = pixel{255, 0, 255, 255} // should not happen
			}

			index[int(qoi_COLOR_HASH(px[0], px[1], px[2], px[3]))%len(index)] = px
		}

		// TODO stride ..
		copy(pixels[:4], px[:])
		pixels = pixels[4:] // advance
	}
	return img, nil
}

func Encode(w io.Writer, m image.Image) error {

	var out = bufio.NewWriter(w)

	minX := m.Bounds().Min.X
	maxX := m.Bounds().Max.X
	minY := m.Bounds().Min.Y
	maxY := m.Bounds().Max.Y

	NBPixels := (maxX - minX) * (maxY - minY)
	if NBPixels == 0 || NBPixels >= qoiPixelsMax {
		return errors.New("Bad image Size")
	}

	// write header to output
	if err := binary.Write(out, binary.BigEndian, []byte(qoiMagic)); err != nil {
		return err
	}
	// width
	if err := binary.Write(out, binary.BigEndian, uint32(maxX-minX)); err != nil {
		return err
	}
	// height
	if err := binary.Write(out, binary.BigEndian, uint32(maxY-minY)); err != nil {
		return err
	}
	// channels
	if err := binary.Write(out, binary.BigEndian, uint8(4)); err != nil {
		return err
	}
	// 0b0000rgba colorspace
	if err := binary.Write(out, binary.BigEndian, uint8(0)); err != nil {
		return err
	}

	var index [64]pixel
	px_prev := pixel{0, 0, 0, 255}
	run := 0

	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			// extract pixel and convert to non-premultiplied
			c := color.NRGBAModel.Convert(m.At(x, y))
			c_r, c_g, c_b, c_a := c.RGBA()
			px := pixel{byte(c_r >> 8), byte(c_g >> 8), byte(c_b >> 8), byte(c_a >> 8)}

			if px == px_prev {
				run++
				last_pixel := x == (maxX-1) && y == (maxY-1)
				if run == 62 || last_pixel {
					err := out.WriteByte(qoiRun | byte(run-1))
					if err != nil {
						return err
					}
					run = 0
				}
			} else {
				if run > 0 {
					err := out.WriteByte(qoiRun | byte(run-1))
					if err != nil {
						return err
					}
					run = 0
				}
				var index_pos byte = qoi_COLOR_HASH(px[0], px[1], px[2], px[3]) % 64
				if index[index_pos] == px {
					err := out.WriteByte(qoiIndex | index_pos)
					if err != nil {
						return err
					}
				} else {
					index[index_pos] = px

					if px[3] == px_prev[3] {
						vr := int8(int(px[0]) - int(px_prev[0]))
						vg := int8(int(px[1]) - int(px_prev[1]))
						vb := int8(int(px[2]) - int(px_prev[2]))

						vg_r := vr - vg
						vg_b := vb - vg

						if vr > -3 && vr < 2 && vg > -3 && vg < 2 && vb > -3 && vb < 2 {

							if err := out.WriteByte(qoiDiff | byte((vr+2)<<4|(vg+2)<<2|(vb+2))); err != nil {
								return util.Catch(err)
							}
						} else if vg_r > -9 && vg_r < 8 && vg > -33 && vg < 32 && vg_b > -9 && vg_b < 8 {
							if err := out.WriteByte(qoiLuma | byte(vg+32)); err != nil {
								return util.Catch(err)
							}

							if err := out.WriteByte(byte((vg_r+8)<<4) | byte(vg_b+8)); err != nil {
								return util.Catch(err)
							}
						} else {
							if err := out.WriteByte(qoiRgb); err != nil {
								return util.Catch(err)
							}

							if err := out.WriteByte(px[0]); err != nil {
								return util.Catch(err)
							}

							if err := out.WriteByte(px[1]); err != nil {
								return util.Catch(err)
							}

							if err := out.WriteByte(px[2]); err != nil {
								return util.Catch(err)
							}
						}

					} else {
						if err := out.WriteByte(qoiRgba); err != nil {
							return util.Catch(err)
						}
						for i := 0; i < 4; i++ {
							if err := out.WriteByte(px[i]); err != nil {
								return util.Catch(err)
							}
						}
					}

				}
			}

			px_prev = px
		}
	}

	if err := binary.Write(out, binary.BigEndian, uint32(0)); err != nil {
		return util.Catch(err)
	} // padding

	if err := binary.Write(out, binary.BigEndian, uint32(1)); err != nil {
		return util.Catch(err)
	} // padding

	return out.Flush()
}

func DecodeConfig(r io.Reader) (cfg image.Config, err error) {
	var header [4 + 4 + 4 + 1 + 1]byte
	if _, err = io.ReadAtLeast(r, header[:], len(header)); err != nil {
		return
	}

	if string(header[:4]) != qoiMagic {
		return cfg, errors.New("Invalid magic")
	}
	// only decodes as NRGBA images
	return image.Config{
		Width:      int(binary.BigEndian.Uint32(header[4:])),
		Height:     int(binary.BigEndian.Uint32(header[8:])),
		ColorModel: color.NRGBAModel,
	}, err
}

func init() {
	image.RegisterFormat("qoi", qoiMagic, Decode, DecodeConfig)
}
