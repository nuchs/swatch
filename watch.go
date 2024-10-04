package main

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/fsnotify/fsnotify"
)

func Watch(config Config, fsys fs.FS) error {
	dirs, err := findDirectories(fsys, config.Excludes)
	if err != nil {
		return fmt.Errorf("Failed to find directories to watch: %w", err)
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()

	for _, d := range dirs {
		w.Add(d)
	}

	if err := runWatcher(config, w); err != nil {
		return err
	}

	return nil
}

func findDirectories(fsys fs.FS, excludes Set) ([]string, error) {
	var dirs []string

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("%s will not be watched: %v\n", path, err)
			return nil
		}
		if excludes.Contains(path) {
			fmt.Printf("Not watching %s\n", path)
			return fs.SkipDir
		}

		if d.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dirs, nil
}

func runWatcher(config Config, w *fsnotify.Watcher) error {
	for {
		select {
		case err, ok := <-w.Errors:
			if !ok {
				fmt.Println("Error channel closed unexpectedly")
				return err
			}
			fmt.Println("Error:", err)
		case e, ok := <-w.Events:
			if !ok {
				return errors.New("Event channel closed unexpectedly")
			}
			handleEvent(config, e)
		}
	}
}

func handleEvent(config Config, e fsnotify.Event) {
	if e.Op&fsnotify.Chmod != 0 {
		return
	}

	fmt.Println(e)
}
