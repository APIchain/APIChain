package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb/storage"
)

var (
	filename string
	child    bool
)

func init() {
	flag.StringVar(&filename, "filename", filepath.Join(os.TempDir(), "goleveldb_filelock_test"), "Filename used for testing")
	flag.BoolVar(&child, "child", false, "This is the child")
}

func runChild() error {
	var args []string
	args = append(args, os.Args[1:]...)
	args = append(args, "-child")
	cmd := exec.Command(os.Args[0], args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	r := bufio.NewReader(&out)
	for {
		line, _, e1 := r.ReadLine()
		if e1 != nil {
			break
		}
		log.Debug("[Child]", string(line))
	}
	return err
}

func main() {
	flag.Parse()

	log.Debug("Using path: %s\n", filename)
	if child {
		log.Debug("Child flag set.")
	}

	stor, err := storage.OpenFile(filename, false)
	if err != nil {
		log.Debug("Could not open storage: %s", err)
		os.Exit(10)
	}

	if !child {
		log.Debug("Executing child -- first test (expecting error)")
		err := runChild()
		if err == nil {
			log.Debug("Expecting error from child")
		} else if err.Error() != "exit status 10" {
			log.Debug("Got unexpected error from child:", err)
		} else {
			log.Debug("Got error from child: %s (expected)\n", err)
		}
	}

	err = stor.Close()
	if err != nil {
		log.Debug("Error when closing storage: %s", err)
		os.Exit(11)
	}

	if !child {
		log.Debug("Executing child -- second test")
		err := runChild()
		if err != nil {
			log.Debug("Got unexpected error from child:", err)
		}
	}

	os.RemoveAll(filename)
}
