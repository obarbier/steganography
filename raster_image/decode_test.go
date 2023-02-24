package raster_image

import (
	"reflect"
	"testing"
)

func Test_encode_decode(t *testing.T) {
	type args struct {
		s string
		r uint8
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid",
			args: args{
				s: "Hello World",
				r: 2,
			},
			want: "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			separated := separateSecretsIntoRValues([]byte(tt.args.s), tt.args.r)
			got := combineChunksIntoSecrets(separated, tt.args.r)
			t.Logf("strGot=%s, strWant=%s", got, tt.want)
			if !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("encode then decode failed = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func Test_combineChunksIntoSecrets(t *testing.T) {
	type args struct {
		s []byte
		r uint8
	}
	tests := []struct {
		name string
		args args
		want []uint8
	}{
		// TODO: Add test cases.
		{
			name: "valid",
			args: args{
				s: []uint8{3, 3, 3, 3, 2, 2, 3, 1, 0, 0, 3, 3, 0, 0, 0, 2},
				r: 2,
			},
			want: []byte{0xFF, 0xAD, 0x0F, 0x02},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := combineChunksIntoSecrets(tt.args.s, tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("combineChunksIntoSecrets() = %v, want %v", got, tt.want)
			}
		})
	}
}
