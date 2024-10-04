package main

import (
	"errors"
	"strings"
)

var (
	ErrNoCommand      = errors.New("You must specify a command")
	ErrMissingDir     = errors.New("You must provide an argument to -d")
	ErrMissingExclude = errors.New("You must provide an argument to -e")
)

type Set struct {
	elements map[string]struct{}
}

func (s Set) Add(e string) {
	s.elements[e] = struct{}{}
}

func (s Set) Contains(e string) bool {
	_, exists := s.elements[e]
	return exists
}

func NewSet(elements ...string) Set {
	s := Set{make(map[string]struct{})}
	for _, e := range elements {
		s.Add(e)
	}
	return s
}

type Config struct {
	Dir      string
	Excludes Set
	Cmd      string
	Args     []string
}

func newConfig() Config {
	return Config{
		".",
		NewSet(".git", "node_modules"),
		"",
		make([]string, 0),
	}
}

func LoadConfig(cmdline []string) (Config, error) {
	config := newConfig()
	numArgs := len(cmdline)

processLoop:
	for i := 0; i < numArgs; i++ {
		switch cmdline[i] {
		case "-d", "--directory":
			i++
			if i >= numArgs {
				return config, ErrMissingDir
			}
			config.Dir = cmdline[i]

		case "-e", "--exclude":
			i++
			if i >= numArgs {
				return config, ErrMissingExclude
			}
			for _, exclusion := range strings.Split(cmdline[i], ",") {
				config.Excludes.Add(strings.TrimPrefix(exclusion, "./"))
			}

		default:
			config.Cmd = cmdline[i]
			if (i + 1) < len(cmdline) {
				config.Args = cmdline[i+1:]
			}
			break processLoop
		}
	}

	if config.Cmd == "" {
		return config, ErrNoCommand
	}

	return config, nil
}
