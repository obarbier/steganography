package raster_image

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/jpeg"
	"io"
	"log"
)

type LSB struct {
	i    *image.Image
	data *LSBData
}

func (j *LSB) Write(writer io.Writer) (int, error) {
	var opt jpeg.Options
	opt.Quality = 98
	err := jpeg.Encode(writer, *j.i, &opt) // https://en.wikipedia.org/wiki/Chroma_subsampling
	if err != nil {
		return 0, err
	}
	//return len(j.i.Pix), nil
	return 120, nil
}

type LSBData struct {
	// payload is a secret seperated in chunks of r bit
	payload     []uint8
	r           uint8
	payloadSize uint32
}

func NewLsbData(secrets []byte, r uint8) *LSBData {
	p := separateSecretsIntoRValues(secrets, r)
	return &LSBData{
		payload:     p,
		r:           r,
		payloadSize: uint32(len(p)),
	}
}

func (l *LSBData) GetBytes() []byte {
	bo := new(bytes.Buffer)
	bo.WriteByte(byte(l.r))

	a := make([]byte, 4)
	binary.LittleEndian.PutUint32(a, l.payloadSize)
	bo.Write(a)

	bo.Write(l.payload)
	return bo.Bytes()
}

func New(reader io.Reader) (*LSB, error) {
	source, format, err := image.Decode(reader)
	log.Printf("reading format %s", format)

	//source, err := jpeg.Decode(reader)

	if err != nil {
		return nil, err
	}
	//bounds := source.Bounds()
	//m := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	//draw.Draw(m, m.Bounds(), source, bounds.Min, draw.Src)
	return &LSB{i: &source}, nil
}
