package raster_image

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"reflect"
	"strings"
	"testing"
)

func NewMockedImage(format string) *MockedImage {
	newImage := image.NewRGBA(image.Rect(0, 0, 125, 125))
	newImage.Pix[0] = 255 // 1st pixel red
	newImage.Pix[1] = 0   // 1st pixel green
	newImage.Pix[2] = 0   // 1st pixel blue
	newImage.Pix[3] = 255 // 1st pixel alpha

	return &MockedImage{
		Format: format,
		i:      newImage,
	}
}

type MockedImage struct {
	Format string
	i      *image.RGBA
}

func (m MockedImage) Read(p []byte) (int, error) {
	buf := new(bytes.Buffer)
	switch strings.ToUpper(m.Format) {
	case "JPEG":
		if err := jpeg.Encode(buf, m.i, nil); err != nil {
			return 0, err
		}
	default:
		return 0, errors.New(fmt.Sprintf("unknown format %s", m.Format))
	}
	copy(p, buf.Bytes())
	return len(p), nil
}

var _ io.Reader = MockedImage{}

func TestNew(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *LSB
		wantErr bool
	}{
		{
			name:    "testing jpeg",
			args:    args{reader: NewMockedImage("jpeg")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			secrets := "Hello World"
			err = got.Encode(secrets)
			if err != nil {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			buf := new(bytes.Buffer)
			_, err = got.Write(buf)
			if err != nil {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got2, err := New(buf)
			if err != nil {
				t.Errorf("New 2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			decode, err := got2.Decode()
			if err != nil {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(string(decode), secrets) {
				t.Errorf(" %v, want %v", string(decode), secrets)
			}

		})
	}
}
