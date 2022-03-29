package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	firstChar, _ := utf8.DecodeRuneInString(s)
	if unicode.IsDigit(firstChar) {
		return "", ErrInvalidString
	}

	var previousChar rune
	var b strings.Builder

	for _, char := range s {
		if unicode.IsDigit(char) {
			if unicode.IsDigit(previousChar) {
				return "", ErrInvalidString
			}

			count, err := strconv.Atoi(string(char))
			if err != nil {
				return "", err
			}
			if count > 0 {
				repeatedPreviousChar := strings.Repeat(string(previousChar), count)
				b.WriteString(repeatedPreviousChar)
			}
		} else if !unicode.IsDigit(previousChar) && previousChar != 0 {
			b.WriteRune(previousChar)
		}

		previousChar = char
	}
	if !unicode.IsDigit(previousChar) {
		b.WriteRune(previousChar)
	}

	return b.String(), nil
}
