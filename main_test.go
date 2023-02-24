package main

import (
	"reflect"
	"testing"
)

func Test_lsbEncode(t *testing.T) {
	type args struct {
		s []byte
		i []byte
		r int
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "z equal to len(i)",
			args: args{
				i: []byte{0x6c, 0x6c, 0xC4, 0xD4},
				s: []byte{0xDD},
				r: 2,
			},
			want:    []byte{0x6F, 0x6D, 0xC7, 0xD5},
			wantErr: false,
		},
		{
			name: "len(i) much larger than len(s) ",
			args: args{
				i: []byte{0x6c, 0x6c, 0xC4, 0xD4, 0x6c, 0x6c, 0xC4, 0xD4, 0x6c, 0x6c, 0xC4, 0xD4, 0x6c, 0x6c, 0xC4, 0xD4},
				s: []byte{0xDD},
				r: 2,
			},
			want:    []byte{0x6F, 0x6D, 0xC7, 0xD5, 0x6c, 0x6c, 0xC4, 0xD4, 0x6c, 0x6c, 0xC4, 0xD4, 0x6c, 0x6c, 0xC4, 0xD4},
			wantErr: false,
		},
		{
			name: "z less than len(i)",
			args: args{
				i: []byte{0x6c, 0x6c, 0xC4, 0xD4, 0xFE, 0xAC},
				s: []byte{0xDD},
				r: 2,
			},
			want:    []byte{0x6F, 0x6D, 0xC7, 0xD5, 0xFE, 0xAC},
			wantErr: false,
		},

		{
			name: "z greater than len(i)",
			args: args{
				i: []byte{0x6c, 0x6c},
				s: []byte{0xDD, 0xDD},
				r: 2,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "r greater than 8",
			args: args{
				i: []byte{0x6c, 0x6c, 0xC4, 0xD4},
				s: []byte{0xDD},
				r: 9,
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "r equal to 8",
			args: args{
				i: []byte{0x6c, 0x6c, 0xC4, 0xD4},
				s: []byte{0xDD},
				r: 8,
			},
			want:    []byte{0xDD, 0x6c, 0xC4, 0xD4},
			wantErr: false,
		},
		//{
		//	name: "uneven r",
		//	args: args{
		//		i: []byte{0x6c, 0x6c, 0xC4, 0xD4},
		//		s: []byte{0xDD},
		//		r: 7,
		//	},
		//	want:    []byte{0x5D, 0x6D, 0xC4, 0xD4},
		//	wantErr: false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lsbEncode(tt.args.s, tt.args.i, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("lsbEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lsbEncode() got = %X, want %X", got, tt.want)
			}
		})
	}
}
