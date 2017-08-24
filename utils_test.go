package html2article

import "testing"

func TestCompress(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple0",
			args: args{
				str: "simple\n",
			},
			want: "simple\n",
		},
		{
			name: "simple1",
			args: args{
				str: " simple ",
			},
			want: " simple ",
		},
		{
			name: "simple2",
			args: args{
				str: "simple 2  ",
			},
			want: "simple 2 ",
		},
		{
			name: "simple3",
			args: args{
				str: "simple 3  \n    ",
			},
			want: "simple 3 ",
		},
		{
			name: "simple4",
			args: args{
				str: "simple4",
			},
			want: "simple4",
		},
		{
			name: "simple5",
			args: args{
				str: "simple4  simple4  \n simple   ",
			},
			want: "simple4 simple4 simple ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Compress(tt.args.str); got != tt.want {
				t.Errorf("Compress() = %v, want %v", got, tt.want)
			}
		})
	}
}
