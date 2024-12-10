package hw02unpackstring

import (
	"errors"
	"github.com/rivo/uniseg"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	gr := uniseg.NewGraphemes(str)

	result := make([][]rune, 0)
	needShield := false
	prevNum := false

	for i := 0; gr.Next(); i++ {
		runes := gr.Runes()
		isChar := len(runes) == 1

		if i == 0 {
			if isChar && unicode.IsDigit(runes[0]) {
				return "", ErrInvalidString
			}

			result = append(result, runes)
			continue
		}

		if !isChar {
			result = append(result, runes)
			prevNum = false
			continue
		}

		sym := runes[0]

		if sym == '\\' {
			prevNum = false

			if needShield {
				result = append(result, runes)
				needShield = false
				continue
			}

			needShield = true
			continue
		}

		if unicode.IsDigit(sym) {
			amount, err := strconv.Atoi(string(sym))

			if err != nil {
				return "", err
			}

			if needShield {
				result = append(result, runes)
				needShield = false
				continue
			} else if prevNum {
				return "", ErrInvalidString
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
		} else {
			if needShield {
				return "", ErrInvalidString
			}

			result = append(result, runes)
			prevNum = false
		}
	}

	resStr := ""

	for _, r := range result {
		resStr = resStr + string(r)
	}

	return resStr, nil

}
