package kobayashi

import (
	"math"
	"strconv"
	"strings"
)

type Unbaser struct {
	base     int
	dict     map[string]int
	selector int
	alphabet map[int]string
}

func NewUnbaser(radix int) *Unbaser {
	var selector int
	if radix > 62 {
		selector = 95
	} else if radix > 54 {
		selector = 62
	} else if radix > 52 {
		selector = 54
	} else {
		selector = 52
	}
	alphabet := make(map[int]string)
	alphabet[52] = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOP"
	alphabet[54] = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQR"
	alphabet[62] = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphabet[95] = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
	d := make(map[string]int)
	for i, v := range strings.Split(alphabet[selector], "") {
		d[v] = i
	}
	return &Unbaser{
		base:     radix,
		selector: selector,
		alphabet: alphabet,
		dict:     d,
	}
}

func (u *Unbaser) Unbase(val string) (int, error) {
	if 2 <= u.base && u.base <= 36 {
		v, err := strconv.ParseInt(val, u.base, 64)
		if err != nil {
			return -1, err
		}
		return int(v), nil
	} else {
		ret := 0
		valArray := arrayReverse(strings.Split(val, ""))

		for i, cipher := range valArray {
			ret += int(math.Pow(float64(u.base), float64(i))) * u.dict[cipher]
		}
		return ret, nil
	}
}

func arrayReverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
