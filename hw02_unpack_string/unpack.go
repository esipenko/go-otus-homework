package hw02unpackstring

import (
	"errors"
	"github.com/rivo/uniseg"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func isSlash(r []rune) bool {
	return len(r) == 1 && r[0] == '\\'
}

func isDigit(r []rune) bool {
	return len(r) == 1 && unicode.IsDigit(r[0])
}

func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	gr := uniseg.NewGraphemes(str)

	result := make([][]rune, 0)
	needShield := false
	prevNum := false

	gr.Next()
	runes := gr.Runes()

	if isDigit(runes) {
		return "", ErrInvalidString
	}

	if isSlash(runes) {
		needShield = true
	} else {
		result = append(result, runes)
	}

	for i := 1; gr.Next(); i++ {
		runes := gr.Runes()

		d := isDigit(runes)
		s := isSlash(runes)

		// Принципиально важны только слеш и число
		if !d && !s {
			if needShield {
				return "", ErrInvalidString
			}

			result = append(result, runes)
			prevNum = false
			continue
		}

		sym := runes[0]

		// Если ловим слеш, хотим что-то экранировать
		if s {
			prevNum = false

			if needShield {
				result = append(result, runes)
				needShield = false
				continue
			}

			needShield = true
			continue
		}

		// Тут уже только число
		if needShield {
			result = append(result, runes)
			needShield = false
			continue
		}

		if prevNum {
			return "", ErrInvalidString
		}

		amount, err := strconv.Atoi(string(sym))

		if err != nil {
			return "", err
		}

		prevNum = true

		if amount == 0 {
			result = result[:len(result)-1]
			continue
		}

		if amount == 1 {
			continue
		}

		for range amount - 1 {
			result = append(result, result[len(result)-1])
		}
	}

	resStr := ""

	for _, r := range result {
		resStr = resStr + string(r)
	}

	return resStr, nil

}
