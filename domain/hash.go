package domain

import "bytes"

var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var base = len(alphabet)

func encode(id int) string {
	var buf bytes.Buffer
	for id > 0 {
		buf.WriteRune(alphabet[id%base])
		id /= base
	}
	return revString(buf.String())
}

func revString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func decode(s string) int {
	var id int
	for _, r := range s {
		id = id*base + find(r)
	}
	return id
}

// find index of char in alphabet
func find(char rune) int {
	for i, r := range alphabet {
		if r == char {
			return i
		}
	}
	return -1
}
