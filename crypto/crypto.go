package crypto

import (
	"crypto/sha256"
	"fmt"
	"io"
)

func Check(r io.Reader, sum string) bool {
	h := sha256.New()
	n, err := io.Copy(h, r)
	switch {
	case err != nil:
		return false
	case n == 0:
		return false
	default:
		return sum == fmt.Sprintf("%x", h.Sum(nil))
	}
}
