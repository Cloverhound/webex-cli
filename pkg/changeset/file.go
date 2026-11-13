package changeset

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Parse decodes a change set, keeping numbers exact.
func Parse(data []byte) (*ChangeSet, error) {
	var cs ChangeSet
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	d.DisallowUnknownFields()
	if err := d.Decode(&cs); err != nil {
		return nil, fmt.Errorf("reading change set: %w", err)
	}
	return &cs, nil
}

// Load reads a change set file.
func Load(path string) (*ChangeSet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

// Encode renders a change set as indented JSON.
func Encode(cs *ChangeSet) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cs); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Save writes a change set atomically.
func Save(path string, cs *ChangeSet) error {
	data, err := Encode(cs)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".changeset-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	// A change set can hold credentials a write sends, such as a trunk's password, so a
	// new file is owner-only and a replaced one keeps its mode.
	mode := os.FileMode(0o600)
	if info, err := os.Stat(path); err == nil {
		mode = info.Mode().Perm()
	}
	if err := os.Chmod(tmp.Name(), mode); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// Edit loads the change set at path, or starts one with create when there is
// no file, applies change and saves the result. A lock file keeps concurrent
// recorders from losing each other's operations.
func Edit(path string, create func() (*ChangeSet, error), change func(*ChangeSet) error) error {
	unlock, err := lock(path)
	if err != nil {
		return err
	}
	defer unlock()

	cs, err := Load(path)
	if errors.Is(err, os.ErrNotExist) {
		cs, err = create()
	}
	if err != nil {
		return err
	}
	if err := change(cs); err != nil {
		return err
	}
	return Save(path, cs)
}

// A lock older than this belongs to a recorder that died.
const staleLock = 30 * time.Second

func lock(path string) (func(), error) {
	name := path + ".lock"
	deadline := time.Now().Add(10 * time.Second)
	for {
		f, err := os.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err == nil {
			f.Close()
			return func() { os.Remove(name) }, nil
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
		if info, statErr := os.Stat(name); statErr == nil && time.Since(info.ModTime()) > staleLock {
			os.Remove(name)
			continue
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("%s is locked by another recorder; remove %s if none is running", path, name)
		}
		time.Sleep(25 * time.Millisecond)
	}
}
