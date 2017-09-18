package captcha

import (
	"bytes"
	"errors"
	"io"
	"time"
)

const (
	DefaultLen = 4
	CollectNum = 100
	Expiration = 10 * time.Minute
)

var (
	ErrNotFound = errors.New("captcha: id not found")
	globalStore = NewMemoryStore(CollectNum, Expiration)
)

func SetCustomStore(s Store) {
	globalStore = s
}

func New() string {
	return NewLen(DefaultLen)
}

func NewLen(length int) (id string) {
	id = RandomId()
	// or bb := RandomLDId(length)
	bb := RandomDigits(length)
	globalStore.Set(id, bb)
	return
}

func Reload(id string) bool {
	old := globalStore.Get(id, false)
	if old == nil {
		return false
	}
	globalStore.Set(id, RandomDigits(len(old)))
	return true
}

func WriteImage(w io.Writer, id string, width, height int) error {
	d := globalStore.Get(id, false)
	if d == nil {
		return ErrNotFound
	}
	_, err := NewImage(id, d, width, height).WriteTo(w)
	return err
}

func Verify(id string, digits []byte) bool {
	if digits == nil || len(digits) == 0 {
		return false
	}
	reald := globalStore.Get(id, true)
	if reald == nil {
		return false
	}
	return bytes.Equal(digits, reald)
}

func VerifyString(id string, digits string) bool {
	if digits == "" {
		return false
	}
	ns := make([]byte, len(digits))
	for i, d := range digits {
		ns[i] = Rune2Digit(d)
	}
	return Verify(id, ns)
}
