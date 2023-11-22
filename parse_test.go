package main

import (
	"reflect"
	"testing"
)

func Test_splitWithStringContext(t *testing.T) {
	type args struct {
		s   string
		sep string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "no strings",
			args: args{"true and false or 123", "and"},
			want: []string{"true ", " false or 123"},
		},
		{
			name: "strings do not contain sep",
			args: args{`"lari" and "ollie"`, "and"},
			want: []string{`"lari" `, ` "ollie"`},
		},
		{
			name: "strings do not contain sep",
			args: args{`"lari" and "ollie" and`, "and"},
			want: []string{`"lari" `, ` "ollie" `},
		},
		{
			name: "ignores separator in strings",
			args: args{`"lari and alex" and "ollie and lari"`, "and"},
			want: []string{`"lari and alex" `, ` "ollie and lari"`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitWithWrappingContext(tt.args.s, tt.args.sep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitWithWrappingContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
