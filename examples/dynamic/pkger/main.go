package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/rhomber/pkger"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	info, err := pkger.Stat("/go.mod")
	if err != nil {
		return err
	}
	fmt.Println(info)
	if err := pkger.MkdirAll("/foo/bar/baz", 0755); err != nil {
		return err
	}

	f, err := pkger.Create("/foo/bar/baz/biz.txt")
	if err != nil {
		return err
	}
	f.Write([]byte("BIZ!!"))
	if err := f.Close(); err != nil {
		return err
	}

	f, err = pkger.Open("/foo/bar/baz/biz.txt")
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, f)
	return f.Close()
}
