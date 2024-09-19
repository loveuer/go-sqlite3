package alloc_test

import (
	"math"
	"testing"

	"github.com/loveuer/go-sqlite3/internal/alloc"
)

func TestVirtual(t *testing.T) {
	defer func() { _ = recover() }()
	alloc.Virtual(math.MaxInt+2, math.MaxInt+2)
	t.Error("want panic")
}
