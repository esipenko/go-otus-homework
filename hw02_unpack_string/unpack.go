package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	res := ""

	tryingShield := false
	hasPrevNum := false
	prevStr := ""

	for _, r := range str {
		isDigit := unicode.IsDigit(r)
		isSlash := r == '\\'

		// Принципиально важны только слеш и число
		if !isDigit && !isSlash {
			if tryingShield {
				return "", ErrInvalidString
			}

			prevStr = string(r)
			res += prevStr

			hasPrevNum = false
			continue
		}

		if tryingShield {
			tryingShield = false

			prevStr = string(r)
			res += prevStr

			continue
		}

		// Если ловим слеш, хотим что-то экранировать
		if isSlash {
			hasPrevNum = false
			tryingShield = true
			continue
		}

		amount, err := strconv.Atoi(string(r))

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
