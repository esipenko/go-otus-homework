package hw02unpackstring

import (
	"errors"
	"github.com/rivo/uniseg"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	gr := uniseg.NewGraphemes(str)
	res := ""

	tryingShield := false
	hasPrevNum := false
	prevStr := ""

	for gr.Next() {
		runes := gr.Runes()

		isDigit := len(runes) == 1 && unicode.IsDigit(runes[0])
		isSlash := len(runes) == 1 && runes[0] == '\\'

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

		amount, err := strconv.Atoi(string(sym))

		// Тут уже только число
		if hasPrevNum || prevStr == "" || err != nil {
			return "", ErrInvalidString
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
