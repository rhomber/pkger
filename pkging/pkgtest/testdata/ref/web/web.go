package web

import (
	"net/http"

	"github.com/rhomber/pkger"
)

func Serve() {
	pkger.Stat("github.com/rhomber/pkger:/README.md")
	dir := http.FileServer(pkger.Dir("/public"))
	http.ListenAndServe(":3000", dir)
}
