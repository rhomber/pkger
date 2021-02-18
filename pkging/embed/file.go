package embed

import (
	"github.com/rhomber/pkger/here"
	"github.com/rhomber/pkger/pkging"
)

type File struct {
	Info   *pkging.FileInfo `json:"info"`
	Here   here.Info        `json:"her"`
	Path   here.Path        `json:"path"`
	Data   []byte           `json:"data"`
	Parent here.Path        `json:"parent"`
}
