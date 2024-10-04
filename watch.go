package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os/exec"
	"time"

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
			fmt.Printf("Watching %s\n", path)
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
	throttled := false
	wait := 200 * time.Millisecond
	t := time.NewTimer(wait)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			throttled = false
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
			if throttled {
				break
			}
			handleEvent(config, e)
			throttled = true
			t.Reset(wait)
		}
	}
}

func handleEvent(config Config, e fsnotify.Event) {
	if e.Op.Has(fsnotify.Chmod) || e.Op.Has(fsnotify.Rename) {
		return
	}

	fmt.Println("-----------------------------")
	fmt.Println(e)
	fmt.Println()

	args := makeArgs(e.Name, config.Args)
	cmd := exec.Command(config.Cmd, args...)
	cmd.Stdout = config.Out
	if err := cmd.Run(); err != nil {
		fmt.Println("An error occurred:", err)
	}
}

func makeArgs(sub string, template []string) []string {
	args := make([]string, len(template))

	for i, v := range template {
		if v == "{}" {
			args[i] = sub
		} else {
			args[i] = template[i]
		}
	}

	return args
}
