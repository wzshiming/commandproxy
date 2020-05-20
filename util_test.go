package commandproxy

import (
	"reflect"
	"testing"
)

func TestSplitCommand(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			args: args{
				cmd: "",
			},
			wantErr: true,
		},
		{
			args: args{
				cmd: "a",
			},
			want: []string{"a"},
		},
		{
			args: args{
				cmd: "a b",
			},
			want: []string{"a", "b"},
		},
		{
			args: args{
				cmd: " a b",
			},
			want: []string{"a", "b"},
		},
		{
			args: args{
				cmd: "a b ",
			},
			want: []string{"a", "b"},
		},
		{
			args: args{
				cmd: "a  b",
			},
			want: []string{"a", "b"},
		},
		{
			args: args{
				cmd: "a  \"",
			},
			wantErr: true,
		},
		{
			args: args{
				cmd: "a  '",
			},
			wantErr: true,
		},
		{
			args: args{
				cmd: "a\t\tb",
			},
			want: []string{"a", "b"},
		},
		{
			args: args{
				cmd: `a \\`,
			},
			want: []string{"a", "\\"},
		},
		{
			args: args{
				cmd: `a \
b`,
			},
			want: []string{"a", "b"},
		},
		{
			args: args{
				cmd: `a "b"`,
			},
			want: []string{"a", "b"},
		},
		{
			args: args{
				cmd: `a 'b'`,
			},
			want: []string{"a", "b"},
		},
		{
			args: args{
				cmd: `a "b "`,
			},
			want: []string{"a", "b "},
		},
		{
			args: args{
				cmd: `a """"`,
			},
			want: []string{"a"},
		},

		{
			args: args{
				cmd: `a "'b'"`,
			},
			want: []string{"a", "'b'"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitCommand(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitCommand() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplaceEscape(t *testing.T) {
	type args struct {
		s  string
		re map[byte]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				s:  "",
				re: map[byte]string{},
			},
			want: "",
		},
		{
			args: args{
				s:  "a",
				re: map[byte]string{},
			},
			want: "a",
		},
		{
			args: args{
				s:  "%",
				re: map[byte]string{},
			},
			want: "%",
		},
		{
			args: args{
				s:  "%%",
				re: map[byte]string{},
			},
			want: "%",
		},
		{
			args: args{
				s:  "%%a",
				re: map[byte]string{},
			},
			want: "%a",
		},
		{
			args: args{
				s:  "a%%",
				re: map[byte]string{},
			},
			want: "a%",
		},
		{
			args: args{
				s:  "aa%%aa",
				re: map[byte]string{},
			},
			want: "aa%aa",
		},
		{
			args: args{
				s:  "%a",
				re: map[byte]string{},
			},
			want: "%a",
		},
		{
			args: args{
				s: "%h:%p",
				re: map[byte]string{
					'h': "127.0.0.1",
					'p': "80",
				},
			},
			want: "127.0.0.1:80",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceEscape(tt.args.s, tt.args.re); got != tt.want {
				t.Errorf("ReplaceEscape() = %v, want %v", got, tt.want)
			}
		})
	}
}
