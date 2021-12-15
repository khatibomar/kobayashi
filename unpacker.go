package kobayashi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrMallFormated = fmt.Errorf("Malformed p.a.c.k.e.r symtab")
)

type Unpacker struct {
	unbaser *Unbaser
	payload string
	symtab  []string
	radix   int
	count   int
}

func NewUnpacker() *Unpacker {
	return &Unpacker{}
}

func (u *Unpacker) Detect(body string) bool {
	body = strings.ReplaceAll(" ", "", body)
	re := regexp.MustCompile(`eval\(function\(p,a,c,k,e,[r|d]?`)
	return re.Match([]byte(body))
}
func (u *Unpacker) Unpack(body string) (string, error) {
	var err error
	re := regexp.MustCompile(`}\('(.*)', *(\d+), *(\d+), *'(.*?)'\.split\('\|'\)`)
	matches := re.FindAllStringSubmatch(body, -1)
	if len(matches) == 0 {
		return "", fmt.Errorf("Failed to get direct link")
	}
	u.payload = matches[0][1]
	u.symtab = strings.Split(matches[0][4], "|")
	u.radix, err = strconv.Atoi(matches[0][2])
	if err != nil {
		return "", err
	}
	u.count, err = strconv.Atoi(matches[0][3])
	if err != nil {
		return "", err
	}
	if u.radix != u.count {
		return "", ErrMallFormated
	}
	u.unbaser = NewUnbaser(u.radix)

	re = regexp.MustCompile(`\b\w+\b`)
	result := re.ReplaceAllStringFunc(u.payload, u.LookUp)
	result = strings.ReplaceAll(result, "\\", "")
	return result, nil
}

func (u *Unpacker) LookUp(matches string) string {
	if len(matches) == 0 {
		return ""
	}
	match := string(matches[0])
	v, _ := u.unbaser.Unbase(match)
	ub := u.symtab[v]
	if len(ub) == 0 {
		return match
	}
	return ub
}
