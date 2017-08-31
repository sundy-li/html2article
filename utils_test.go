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
			name: "test0",
			args: args{
				str: "test ",
			},
			want: "test ",
		},
		{
			name: "test1",
			args: args{
				str: " test ",
			},
			want: " test ",
		},
		{
			name: "test2",
			args: args{
				str: "test 2  ",
			},
			want: "test 2 ",
		},
		{
			name: "test3",
			args: args{
				str: "test 3  \n    ",
			},
			want: "test 3 ",
		},
		{
			name: "test4",
			args: args{
				str: "test4",
			},
			want: "test4",
		},
		{
			name: "test5",
			args: args{
				str: "test5  test5  \n test   ",
			},
			want: "test5 test5 test ",
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

func TestDistance(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want int
	}{
		{"1", "abc", "ab", 1},
		{"2", "abc", "abd", 1},
		{"3", "ab", "abcef", 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := distance(tt.a, tt.b); got != tt.want {
				t.Errorf("Distance() = %v, want %v", got, tt.want)
			}
		})
	}
}
