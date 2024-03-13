package cards

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

//Paginate returns a 0-indexed page of anything and the page count
func Paginate[T any](items []T, page, count int) (paged []T, pageCount int) {
	if page < 0 || len(items) <= 0 {
		return
	}
	start := page*count
	if start >= len(items) {
		return
	}
	pageCount = int(math.Ceil(float64(len(items)) / float64(count)))
	if pageCount <= 0 {
		pageCount = 1
	}
	for i := 0; i < count; i++ {
		if start+i >= len(items) {
			return
		}
		paged = append(paged, items[start+i])
	}
	return
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
    letterIdxBits = 6                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
var src = rand.NewSource(time.Now().UnixNano())
func RandomString(n int) string {
    b := make([]byte, n)
    // A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
    for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = src.Int63(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }

    return *(*string)(unsafe.Pointer(&b))
}
func RandomStringUpper(n int) string {
	return strings.ToUpper(RandomString(n))
}

func GetColor(hex string) int {
	if hex == "" {
		return 0
	}
	if hex[0] == '#' {
		hex = string(hex[1:])
	}
	switch hex {
	case "red":
		hex = "FF0000"
	case "green":
		hex = "00FF00"
	case "blue":
		hex = "0000FF"
	case "purple":
		hex = "800080"
	case "yellow":
		hex = "FFFF00"
	case "white":
		hex = "FFFFFF"
	case "black":
		hex = "000000"
	}
	dec, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return 0
	}
	return int(dec)
}