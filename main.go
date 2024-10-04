package main

import (
	"fmt"
	"io/fs"
	"os"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		usage(err)
		os.Exit(1)
	}
}

func run(cmdline []string) error {
	config, err := LoadConfig(cmdline)
	if err != nil {
		return err
	}

	fsys := os.DirFS(config.Dir)
	dirs := findDirectories(fsys, config.Excludes)

	watch(config, dirs)

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

func watch(config Config, dirs []string) {

}

func findDirectories(fsys fs.FS, excludes Set) []string {
	var dirs []string

	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("%s will not be watched: %v\n", path, err)
			return nil
		}
		if excludes.Contains(path) {
			fmt.Printf("Skipping %s\n", path)
			return fs.SkipDir
		}

		if d.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})

	return dirs
}
