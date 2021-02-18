package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"app/actions"

	"github.com/rhomber/pkger"
)

func main() {
	pkger.Include("/web")
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

const (
	unused, pathAsset = "", "/assets"
)

func run() error {
	if err := actions.WalkTemplates(os.Stdout); err != nil {
		return err
	}

	err := pkger.Walk(pathAsset, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func dynamic() error {
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
