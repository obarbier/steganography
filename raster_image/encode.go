package raster_image

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"reflect"
)
import _ "image/jpeg"

const DEFAULT_R = uint8(2)

// TODO: this package should only be Encode
//func Encode(reader io.Reader, secret string) ([]byte, error) {
//	steg, err := newLSBSteg(reader)
//	if err != nil {
//		return nil, err
//	}
//
//	encodeErr := steg.Encode([]byte(secret))
//	if encodeErr != nil {
//		return nil, encodeErr
//	}
//
//}

func (j *LSB) Encode(secret string) error {
	log.Printf("secretBytes=%x\n", []byte(secret))
	dataTobeEncoded := []byte(secret)
	j.data = NewLsbData([]byte(secret), DEFAULT_R)
	rgba64, err := lsbEncode(*j.i, j.data)
	if err != nil {
		return err
	}
	j.i = &rgba64
	log.Printf("payloadSize=%d", j.data.payloadSize)

	decodedByte, decodeErr := j.Decode()
	if decodeErr != nil {
		log.Printf("error while trying to validate encoding: %s", decodeErr)
		return err
	}
	if !reflect.DeepEqual(dataTobeEncoded, decodedByte) {
		log.Fatalf("encoded bytes %x doesn't equal decoded byte %x", dataTobeEncoded, decodedByte)
	}

	return nil
}

func lsbEncode(m image.Image, data *LSBData) (image.Image, error) {
	sStar := data.payload
	log.Printf("dataToEncode=%d\n", sStar)
	data.payloadSize = 0
	switch imageType := m.(type) {
	case *image.YCbCr:
		nM := m.(*image.YCbCr)
		bounds := nM.Bounds()
		o := image.NewYCbCr(bounds, nM.SubsampleRatio)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				c := nM.YCbCrAt(x, y)
				yd := c.Y
				cb := c.Cb
				cr := c.Cr
				shifter := DEFAULT_R
				if len(sStar) != 0 {
					_s := sStar[0]
					iStar := yd & (0xff << shifter)
					yd = iStar ^ _s
					sStar = sStar[1:]
					data.payloadSize++
				}
				//if len(sStar) != 0 {
				//	_s := sStar[0]
				//	iStar := cb & (0xff << shifter)
				//	cb = iStar ^ _s
				//	sStar = sStar[1:]
				//	data.payloadSize++
				//}
				//if len(sStar) != 0 {
				//	_s := sStar[0]
				//	iStar := cr & (0xff << shifter)
				//	cr = iStar ^ _s
				//	sStar = sStar[1:]
				//	data.payloadSize++
				//}

				o.Y[o.YOffset(x, y)] = yd
				o.Cb[o.COffset(x, y)] = cb
				o.Cr[o.COffset(x, y)] = cr

			}
		}

		return o, nil
	case *image.RGBA:
		nM := m.(*image.RGBA)
		bounds := nM.Bounds()
		o := image.NewRGBA(bounds)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := nM.At(x, y).RGBA()
				nR := uint8(r >> 8)
				nG := uint8(g >> 8)
				nB := uint8(b >> 8)
				shifter := DEFAULT_R
				if len(sStar) != 0 {
					_s := sStar[0]
					iStar := nR & (0xff << shifter)
					nR = iStar ^ _s
					sStar = sStar[1:]
					data.payloadSize++
				}
				if len(sStar) != 0 {
					_s := sStar[0]
					iStar := nG & (0xff << shifter)
					nG = iStar ^ _s
					sStar = sStar[1:]
					data.payloadSize++
				}
				if len(sStar) != 0 {
					_s := sStar[0]
					iStar := nB & (0xff << shifter)
					nB = iStar ^ _s
					sStar = sStar[1:]
					data.payloadSize++
				}
				o.Set(x, y, color.RGBA{
					R: nR,
					G: nG,
					B: nB,
				})
			}
		}

		return o, nil

	default:
		return nil, errors.New(fmt.Sprintf("%s unsupported image type", imageType))
	}

}

func separateSecretsIntoRValues(s []byte, r uint8) []uint8 {
	p := make([]uint8, 0)
	size := uint8(8)
	for _, _s := range s {
		available := size
		for available != 0 {
			p = append(p, _s>>(size-r))
			available -= r
			if available < r && available != 0 {
				p = append(p, _s>>(size-available))
				break
			}
			_s = _s << r
		}
	}
	return p
}
