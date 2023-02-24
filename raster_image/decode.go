package raster_image

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"log"
)

func (j *LSB) Decode() ([]byte, error) {
	return lsbDecode(*j.i)
}

func lsbDecode(m image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	size := 44 // FIXME: this should be embedded in image bytes
	switch imageType := m.(type) {
	case *image.YCbCr:
		nM := m.(*image.YCbCr)
		bounds := nM.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				c := nM.YCbCrAt(x, y)
				yd := c.Y
				//cb := c.Cb
				//cr := c.Cr

				shifter := DEFAULT_R
				if size > 0 {
					iStar := yd & (0xff >> (8 - shifter))
					buf.WriteByte(iStar)
					size--
				}
				//if size > 0 {
				//	iStar := cb & (0xff >> (8 - shifter))
				//	buf.WriteByte(iStar)
				//	size--
				//}
				//if size > 0 {
				//	iStar := cr & (0xff >> (8 - shifter))
				//	buf.WriteByte(iStar)
				//	size--
				//}
				if size == 0 {
					goto DONE
				}
			}
		}
	case *image.RGBA:
		nM := m.(*image.RGBA)
		bounds := nM.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := m.At(x, y).RGBA()
				nR := uint8(r)
				nG := uint8(g)
				nB := uint8(b)
				shifter := uint8(DEFAULT_R)
				if size > 0 {
					iStar := uint8(nR & (0xff >> (8 - shifter)))
					buf.WriteByte(iStar)
					size--
				}
				if size > 0 {
					iStar := uint8(nG & (0xff >> (8 - shifter)))
					buf.WriteByte(iStar)
					size--
				}
				if size > 0 {
					iStar := uint8(nB & (0xff >> (8 - shifter)))
					buf.WriteByte(iStar)
					size--
				}
				if size == 0 {
					goto DONE
				}
			}
		}

	default:
		return nil, errors.New(fmt.Sprintf("%s unsupported image type", imageType))
	}
DONE:
	log.Printf("dataToDecode=%d\n", buf.Bytes())
	secretBytes := combineChunksIntoSecrets(buf.Bytes(), DEFAULT_R)
	log.Printf("secretBytes=%x\n", []byte(secretBytes))
	return secretBytes, nil
}
func combineChunksIntoSecrets(s []byte, r uint8) []byte {
	p := make([]uint8, 0)
	size := 8 / r
	for len(s) != 0 {
		var val []uint8
		if len(s) >= int(size) {
			val = s[:size]
			p = append(p, (val[0]<<6)|(val[1]<<4)|(val[2]<<2)|val[3])
		} else {
			var lastVal uint8
			for _, d := range s {
				lastVal = (lastVal << r) | d
			}
			p = append(p, lastVal)
			break
		}
		s = s[size:]
	}
	return p
}
