package main

import (
	"reflect"
	"testing"
	"testing/fstest"
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
		desc     string
		fsys     fstest.MapFS
		excludes Set
		want     []string
	}{
		{
			desc:     "Flat dir",
			fsys:     makeFs(),
			excludes: NewSet(),
			want:     []string{"."},
		},
		{
			desc:     "Multiple dir, no files",
			fsys:     makeFs("dir1/file1", "dir1/file2", "dir2/file3", "dir2/dir21/file4"),
			excludes: NewSet(),
			want:     []string{".", "dir1", "dir2", "dir2/dir21"},
		},
		{
			desc:     "Excludes",
			fsys:     makeFs("dir1/file1", "dir1/file2", "dir2/file3", "dir2/dir21/file4"),
			excludes: NewSet("dir2", "dir333"),
			want:     []string{".", "dir1"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := findDirectories(tC.fsys, tC.excludes)
			if err != nil {
				t.Fatalf("Error finding directories: %v", err)
			}

			if !reflect.DeepEqual(got, tC.want) {
				t.Fatalf("got %q, want %q", got, tC.want)
			}
		})
	}
}

func makeFs(files ...string) fstest.MapFS {
	stub := fstest.MapFile{}
	fsys := make(map[string]*fstest.MapFile, len(files))
	for _, f := range files {
		fsys[f] = &stub
	}
	return fsys
}
