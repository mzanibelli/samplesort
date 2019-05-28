package crypto_test

import (
	"bytes"
	"samplesort/crypto"
	"testing"
)

func TestCheck(t *testing.T) {
	t.Run("it should return true if the SHA256 sum matches",
		func(t *testing.T) {
			b := bytes.NewBuffer([]byte("hello world"))
			sum := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
			if !crypto.Check(b, sum) {
				t.Errorf("checksum mismatch")
			}
		})
}
