package raster_image

import (
	"reflect"
	"testing"
)

func Test_separateSecretsIntoRValues(t *testing.T) {
	type args struct {
		s []byte
		r uint8
	}
	tests := []struct {
		name string
		args args
		want []uint8
	}{
		{
			name: "basic test",
			args: args{
				s: []byte{0xFF, 0xAD, 0x0F, 0x02},
				r: 2,
			},
			want: []uint8{3, 3, 3, 3, 2, 2, 3, 1, 0, 0, 3, 3, 0, 0, 0, 2},
		},
		{
			name: "Hello world test",
			args: args{
				s: []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64},
				r: 2,
			},
			want: []uint8{1, 2, 2, 0, 1, 2, 1, 1, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 3, 0, 2, 0, 0, 1, 3, 1, 3, 1, 2, 3, 3, 1, 3, 0, 2, 1, 2, 3, 0, 1, 2, 1, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := separateSecretsIntoRValues(tt.args.s, tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("separateSecretsIntoRValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
