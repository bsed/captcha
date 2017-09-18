package captcha

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

// 常量：可以设置的字符的范围
const (
	C_DIGIT = 10
	C_UPPER = 36
	C_LOWER = 62
	// idLen: caphcha id 的长度
	idLen = 20
)

var rand_mod byte = C_UPPER

// idChars 可以使用的字符
var idChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
var lowerCaseDigitsChar = []byte("abcdefghijklmnopqrstuvwxyz0123456789")

var rngKey [32]byte

func init() {
	if _, err := io.ReadFull(rand.Reader, rngKey[:]); err != nil {
		panic("captcha: error reading random source: " + err.Error())
	}
}

/*
** UUID generation
 */

func RandomUUIDBytes() []byte {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		panic(fmt.Sprintf("RandomUUIDBytes# len:%d, err:%v", n, err))
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return uuid
}

func RandomUUID() string {
	uuid := RandomUUIDBytes()
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func RandomUUIDNoDashes() string {
	uuid := RandomUUIDBytes()
	return fmt.Sprintf("%x%x%x%x%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

const (
	imageSeedPurpose = 0x01
	audioSeedPurpose = 0x02
)

func SetCharacterRange(rang byte) {
	if rang == C_DIGIT || rang == C_UPPER || rang == C_LOWER {
		rand_mod = rang
	}
}

func deriveSeed(purpose byte, id string, digits []byte) (out [16]byte) {
	var buf [sha256.Size]byte
	h := hmac.New(sha256.New, rngKey[:])
	h.Write([]byte{purpose})
	io.WriteString(h, id)
	h.Write([]byte{0})
	h.Write(digits)
	sum := h.Sum(buf[:0])
	copy(out[:], sum)
	return
}

func RandomDigits(length int) []byte {
	b := randomBytesMod(length, rand_mod)
	if rand_mod == C_UPPER {
		for i, c := range b {
			if c == 0 {
				b[i] = c + 2
			}
			if c == 1 {
				b[i] = c + 1
			}

			if 9 < c && c < 36 {
				b[i] = c + 26
			}
			if b[i] == 47 || b[i] == 44 || b[i] == 50 {
				b[i]++
			}

		}
	}
	return b
	//return randomBytesMod(length, rand_mod)
}

func randomBytes(length int) (b []byte) {
	b = make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic("captcha: error reading random source: " + err.Error())
	}
	return
}

func randomBytesMod(length int, mod byte) (b []byte) {
	if length == 0 {
		return nil
	}
	if mod == 0 {
		panic("captcha: bad mod argument for randomBytesMod")
	}
	maxrb := 255 - byte(256%int(mod))
	b = make([]byte, length)
	i := 0
	for {
		r := randomBytes(length + (length / 4))
		for _, c := range r {
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = c % mod
			i++
			if i == length {
				return
			}
		}
	}

}

func RandomId() string {
	b := randomBytesMod(idLen, byte(len(idChars)))
	for i, c := range b {
		b[i] = idChars[c]
	}
	return string(b)
}

func RandomLDId(length int) []byte {
	b := randomBytesMod(length, byte(len(lowerCaseDigitsChar)))
	for i, c := range b {
		if 10 < c && c < 36 {
			b[i] = c + 26
		}
	}
	return b
}
