package raster_image

import (
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
	err := lsbEncode(j)
	if err != nil {
		return err
	}
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

func lsbEncode(j *LSB) error {
	sStar := j.data.payload
	log.Printf("dataToEncode=%d\n", sStar)
	bounds := j.i.Bounds()
	j.data.payloadSize = 0
	o := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := j.i.At(x, y).RGBA()
			nR := uint8(r)
			nG := uint8(g)
			nB := uint8(b)
			shifter := DEFAULT_R
			if len(sStar) != 0 {
				_s := sStar[0]
				iStar := nR & (0xff << shifter)
				nR = iStar ^ _s
				sStar = sStar[1:]
				j.data.payloadSize++
			}
			if len(sStar) != 0 {
				_s := sStar[0]
				iStar := nG & (0xff << shifter)
				nG = iStar ^ _s
				sStar = sStar[1:]
				j.data.payloadSize++
			}
			if len(sStar) != 0 {
				_s := sStar[0]
				iStar := nB & (0xff << shifter)
				nB = iStar ^ _s
				sStar = sStar[1:]
				j.data.payloadSize++
			}
			o.Set(x, y, color.NRGBA{
				R: nR,
				G: nG,
				B: nB,
				A: uint8(a),
			})
		}
	}
	j.i = o
	return nil
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
