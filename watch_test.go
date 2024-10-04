package main

import (
	"reflect"
	"testing"
)

func TestArgumentSubstitution(t *testing.T) {
	testCases := []struct {
		desc string
		args []string
		want []string
	}{
		{
			desc: "Empty",
			args: []string{},
			want: []string{},
		},
		{
			desc: "No subs",
			args: []string{"1", "2", "3"},
			want: []string{"1", "2", "3"},
		},
		{
			desc: "Substite placeholder",
			args: []string{"{}", "2", "{}"},
			want: []string{"X", "2", "X"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := makeArgs("X", tC.args)
			if !reflect.DeepEqual(got, tC.want) {
				t.Fatalf("Got %q, want %q", got, tC.want)
			}
		})
	}
}

func TestFindDirs(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}
