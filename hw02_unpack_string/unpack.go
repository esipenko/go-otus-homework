package hw02unpackstring

import (
	"errors"
	"github.com/rivo/uniseg"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func checkSlash(r []rune) bool {
	return len(r) == 1 && r[0] == '\\'
}

func checkDigit(r []rune) bool {
	return len(r) == 1 && unicode.IsDigit(r[0])
}

func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	gr := uniseg.NewGraphemes(str)

	res := ""

	tryingShield := false
	hasPrevNum := false
	prevStr := ""

	gr.Next()
	runes := gr.Runes()

	if checkDigit(runes) {
		return "", ErrInvalidString
	}

	if checkSlash(runes) {
		tryingShield = true
	} else {
		prevStr = string(runes)
		res += prevStr
	}

	for i := 1; gr.Next(); i++ {
		runes := gr.Runes()

		isDigit := checkDigit(runes)
		isSlash := checkSlash(runes)

		// Принципиально важны только слеш и число
		if !isDigit && !isSlash {
			if tryingShield {
				return "", ErrInvalidString
			}

			prevStr = string(runes)
			res += prevStr

			hasPrevNum = false
			continue
		}

		sym := runes[0]

		if tryingShield {
			tryingShield = false

			prevStr = string(runes)
			res += prevStr

			continue
		}

		// Если ловим слеш, хотим что-то экранировать
		if isSlash {
			hasPrevNum = false
			tryingShield = true
			continue
		}

		// Тут уже только число
		if hasPrevNum {
			return "", ErrInvalidString
		}

		amount, err := strconv.Atoi(string(sym))

		if err != nil {
			return "", err
		}

		hasPrevNum = true

		if amount == 0 {
			end := len(res) - len(prevStr)
			res = res[:end]
			continue
		}

		res += strings.Repeat(prevStr, amount-1)
	}

	return res, nil

}
