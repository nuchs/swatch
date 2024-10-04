package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		usage(err)
		os.Exit(1)
	}
}

func run(cmdline []string, out io.Writer) error {
	config, err := LoadConfig(cmdline, out)
	if err != nil {
		return fmt.Errorf("Failed to load config: %w", err)
	}

	fsys := os.DirFS(config.Dir)
	if err := Watch(config, fsys); err != nil {
		return fmt.Errorf("Failed to watch directories: %w", err)
	}

	return nil
}

func usage(err error) {
	fmt.Println("ERROR:", err)
	fmt.Println(`
Swatch is a simple watcher program. It watches for changes in a directory and
runs the specified command when they occur. 

usage:
    swatch [-d <directory>] [-e <exclude list>] command [args...]

options:
    -d, --directory   The directory to watch. By default the current directory 
                      is watched.
    -e, --exclude     A comma separated list of directories to exclude from 
                      being watched. If no value is specified it will exclude 
                      the .git & node_modules directories.
    command           The command to run when a file changes. Must come after 
                      the options.
    args              Arguments passed to command
  `)
}
