package main

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		desc    string
		cmdline []string
		want    Config
	}{
		{
			desc:    "command only",
			cmdline: []string{"ls"},
			want:    Config{".", NewSet(".git", "node_modules"), "ls", []string{}, nil},
		},
		{
			desc:    "command with args",
			cmdline: []string{"ls", "-l", "-h"},
			want:    Config{".", NewSet(".git", "node_modules"), "ls", []string{"-l", "-h"}, nil},
		},
		{
			desc:    "-d",
			cmdline: []string{"-d", "/home", "ls"},
			want:    Config{"/home", NewSet(".git", "node_modules"), "ls", []string{}, nil},
		},
		{
			desc:    "-directory",
			cmdline: []string{"--directory", "/home", "ls"},
			want:    Config{"/home", NewSet(".git", "node_modules"), "ls", []string{}, nil},
		},
		{
			desc:    "-e single exclude",
			cmdline: []string{"-e", "a", "ls"},
			want:    Config{".", NewSet(".git", "node_modules", "a"), "ls", []string{}, nil},
		},
		{
			desc:    "-e multiple excludes",
			cmdline: []string{"-e", "a,.b,./c", "ls"},
			want:    Config{".", NewSet(".git", "node_modules", "a", ".b", "c"), "ls", []string{}, nil},
		},
		{
			desc:    "--exclude single exclude",
			cmdline: []string{"--exclude", "a", "ls"},
			want:    Config{".", NewSet(".git", "node_modules", "a"), "ls", []string{}, nil},
		},
		{
			desc:    "--exclude multiple excludes",
			cmdline: []string{"--exclude", "a,.b,./c", "ls"},
			want:    Config{".", NewSet(".git", "node_modules", "a", ".b", "c"), "ls", []string{}, nil},
		},
		{
			desc:    "All options",
			cmdline: []string{"-d", "/a", "-e", "b,c", "ls", "-l", "-r"},
			want:    Config{"/a", NewSet(".git", "node_modules", "b", "c"), "ls", []string{"-l", "-r"}, nil},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := LoadConfig(tC.cmdline, nil)
			if err != nil {
				t.Fatalf("Error loading config: %v", err)
			}

			if !reflect.DeepEqual(got, tC.want) {
				t.Fatalf("Config incorrect, got %q, want %q", got, tC.want)
			}
		})
	}
}

func TestBadConfig(t *testing.T) {
	testCases := []struct {
		desc    string
		cmdline []string
		want    error
	}{
		{
			desc:    "No command",
			cmdline: []string{},
			want:    ErrNoCommand,
		},
		{
			desc:    "No argument to -e",
			cmdline: []string{"-e"},
			want:    ErrMissingExclude,
		},
		{
			desc:    "No argument to --exclude",
			cmdline: []string{"--exclude"},
			want:    ErrMissingExclude,
		},
		{
			desc:    "No argument to -d",
			cmdline: []string{"-d"},
			want:    ErrMissingDir,
		},
		{
			desc:    "No argument to --directory",
			cmdline: []string{"--directory"},
			want:    ErrMissingDir,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, got := LoadConfig(tC.cmdline, nil)
			if got == nil {
				t.Fatalf("Error should have occurred loading config")
			}
			if got != tC.want {
				t.Fatalf("Wrong error loading config, got %q, want %q", got, tC.want)
			}
		})
	}
}

func TestSetMembership(t *testing.T) {
	s := NewSet()
	s.Add("bacon")

	if !s.Contains("bacon") {
		t.Fatalf("Set does not contain added key: %+v", s)
	}
}

func TestSetNotMember(t *testing.T) {
	s := NewSet("bacon")

	if s.Contains("egg") {
		t.Fatalf("Set should not contain key 'egg': %+v", s)
	}
}

func TestSetAddOnCreate(t *testing.T) {
	elements := []string{"bacon", "eggs", "sausage"}
	s := NewSet(elements...)

	for _, e := range elements {
		if !s.Contains(e) {
			t.Errorf("Set does not contain key %s: %+v", e, s)
		}
	}
}
