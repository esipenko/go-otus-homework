package hw02unpackstring

import (
	"errors"
	"github.com/rivo/uniseg"
	"strconv"
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

	result := make([][]rune, 0)
	tryingShield := false
	prevNum := false

	gr.Next()
	runes := gr.Runes()

	if checkDigit(runes) {
		return "", ErrInvalidString
	}

	if checkSlash(runes) {
		tryingShield = true
	} else {
		result = append(result, runes)
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

			result = append(result, runes)
			prevNum = false
			continue
		}

		sym := runes[0]

		if tryingShield {
			result = append(result, runes)
			tryingShield = false
			continue
		}

		// Если ловим слеш, хотим что-то экранировать
		if isSlash {
			prevNum = false
			tryingShield = true
			continue
		}

		// Тут уже только число
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

		for i := 1; i < amount; i++ {
			result = append(result, result[len(result)-1])
		}
	}

	resStr := ""

	for _, r := range result {
		resStr += string(r)
	}

	return resStr, nil

}
