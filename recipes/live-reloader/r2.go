package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type RecursiveWatcher struct {
	Root    fs.FileInfo
	Watcher *fsnotify.Watcher
	Events  chan fsnotify.Event
	Errors  chan error
}

/*
func (w *RecursiveWatcher) ReadDir(name string) ([]fs.DirEntry, error) {
	fullPath := filepath.Join(w.Root.Name(), name)
	files, err := os.ReadDir(fullPath)
	return files, err
}
*/

func (r2 *RecursiveWatcher) Start() {

	addWatchForEveryDir := func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = r2.Watcher.Add(path)
			if err != nil {
				r2.Errors <- err
			}
		}
		return nil
	}
	filepath.WalkDir(r2.Root.Name(), addWatchForEveryDir)

	go func() {
		for ev := range r2.Watcher.Events {
			if ev.Op.Has(fsnotify.Create) {
				info, err := os.Stat(ev.Name)
				if err == nil {
					if info.IsDir() {
						//	add the directory as a watcher
						filepath.WalkDir(ev.Name, addWatchForEveryDir)
					}
				}
			}
			if ev.Op.Has(fsnotify.Remove) {
				info, err := os.Stat(ev.Name)
				if err == nil {
					if info.IsDir() {
						//	check if there is an associated watcher
						//	check if there are children watchers
						//	delete them all
						err := r2.Watcher.Remove(ev.Name)
						if err != nil {
							r2.Errors <- err
						}
					}
				}
			}
			r2.Events <- ev
		}
	}()

	//	we pass errors unmolested from watcher to our root object
	//	in the future, we may need to intervene
	go func() {
		for er := range r2.Watcher.Errors {
			r2.Errors <- er
		}
	}()

}

func (r2 *RecursiveWatcher) Shutdown() {
	r2.Watcher.Close()
	close(r2.Events)
}

func NewRecursiveWatcher(dir string) (*RecursiveWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("couldn't create a RecursiveWatcher. Couldn't create an fsnotify.Watcher: %w", err)
	}
	err = watcher.Add(dir)
	if err != nil {
		return nil, fmt.Errorf("couldn't create a RecursiveWatcher. Couldn't add %s: %w", dir, err)
	}
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("couldn't create a RecursiveWatcher. Couldn't open %s: %w", dir, err)
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("couldn't create a RecursiveWatcher. %q is not a a directory", dir)
	}

	r2 := RecursiveWatcher{
		Root:    fileInfo,
		Watcher: watcher,
		Events:  make(chan fsnotify.Event),
		Errors:  make(chan error),
	}

	r2.Start()
	return &r2, nil
}
