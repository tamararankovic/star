package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

type Queue struct {
	Files []string
}

func ts() string {
	return string(time.Now().Unix())
}

func (q *Queue) Store(data []byte) error {
	f, err := os.Create("/cache")
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		f.Close()
		return err
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (q *Queue) Load() error {
	files, err := ioutil.ReadDir("/cache")
	if err != nil {
		return err
	}

	for _, file := range files {
		q.Files = append(q.Files, file.Name())
	}

	sort.Strings(q.Files)

	return nil
}
