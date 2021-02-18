package parser

import (
	"encoding/json"

	"github.com/rhomber/pkger/here"
)

type File struct {
	Abs  string // full path on disk to file
	Path here.Path
	Here here.Info
}

func (f File) String() string {
	b, _ := json.MarshalIndent(f, "", " ")
	return string(b)
}
