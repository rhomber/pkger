package mem_test

import (
	"testing"

	"github.com/rhomber/pkger/pkging"
	"github.com/rhomber/pkger/pkging/mem"
	"github.com/rhomber/pkger/pkging/pkgtest"
)

func Test_Pkger(t *testing.T) {
	pkgtest.All(t, func(ref *pkgtest.Ref) (pkging.Pkger, error) {
		return mem.New(ref.Info)
	})
}
