package main

import (
	"testing"
)

func TestCompiler(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should return the correct ast",
			args: args{
				input: "(add 2 (subtract 4 2))",
			},
			want: "add(2, subtract(4, 2));",
		},
		{
			name: "should return the correct ast",
			args: args{
				input: "(add 2 (subtract 4 2) (add 2 3))",
			},
			want: "add(2, subtract(4, 2), add(2, 3));",
		},
		{
			name: "should return the correct ast",
			args: args{
				input: "(add 2 (subtract 4 2) (add 2 3) (subtract 4 2))",
			},
			want: "add(2, subtract(4, 2), add(2, 3), subtract(4, 2));",
		},
		{
			name: "should return the correct ast",
			args: args{
				input: "(add 2 (subtract 4 2) (add 2 3) (subtract 4 2) (add 2 3))",
			},
			want: "add(2, subtract(4, 2), add(2, 3), subtract(4, 2), add(2, 3));",
		},
		{
			name: "should return the correct ast",
			args: args{
				input: "(add 2 (subtract 4 2) (add 2 3) (subtract 4 2) (add 2 3) (subtract 4 2))",
			},
			want: "add(2, subtract(4, 2), add(2, 3), subtract(4, 2), add(2, 3), subtract(4, 2));",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Compiler(tt.args.input); got != tt.want {
				t.Errorf("Compiler() = %v, want %v", got, tt.want)
			}
		})
	}
}
