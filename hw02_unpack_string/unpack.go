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

	result := make([]rune, 0)

	gr := uniseg.NewGraphemes(str)
	gr.Next()

	r := gr.Runes()

	if len(r) == 1 {
		if unicode.IsDigit(r[0]) {
			return "", ErrInvalidString
		}
	}

	result = append(result, r...)
	pChar := r
	for gr.Next() {
		r = gr.Runes()

		if len(r) > 1 {
			result = append(result, r...)
			pChar = r
			continue
		}

		if unicode.IsDigit(r[0]) {
			amount, err := strconv.Atoi(string(r))
			isChar := len(pChar) == 1

			if err != nil {
				return "", err
			}

			if isChar && unicode.IsDigit(pChar[0]) {
				return "", ErrInvalidString
			}

			if amount == 0 {
				result = result[:len(result)-len(pChar)]
				continue
			}

			for range amount - 1 {
				result = append(result, pChar...)
			}
		} else {
			result = append(result, r...)
		}

		pChar = r
	}

	//for gr.Next() {
	//	r = gr.Runes()
	//	t := string(r)
	//	fmt.Println(t)
	//
	//	if len(r) > 1 {
	//		result = append(result, r...)
	//		pChar = r
	//		continue
	//	}
	//
	//	if unicode.IsDigit(r[0]) {
	//		amount, err := strconv.Atoi(string(r))
	//		isChar := len(pChar) == 1
	//
	//		if err != nil {
	//			return "", err
	//		}
	//
	//		if isChar && unicode.IsDigit(pChar[0]) {
	//			return "", ErrInvalidString
	//		}
	//
	//		//if isChar {
	//		//	if pChar[0] == '\\' {
	//		//		result = append(result, r...)
	//		//		pChar = append(pChar, r...)
	//		//		continue
	//		//	}
	//		//} else {
	//		//	if pChar[0] == '\\' && pChar[1] == '\\' {
	//		//
	//		//	}
	//		//}
	//
	//		if amount == 0 {
	//			result = result[:len(result)-len(pChar)]
	//			continue
	//		}
	//
	//		for range amount - 1 {
	//			result = append(result, pChar...)
	//		}
	//	} else if string(r[0]) == `\` {
	//		if pChar[0] == '\\' {
	//			result = append(result, r...)
	//		}
	//
	//		pChar = r
	//	} else {
	//		result = append(result, r...)
	//	}
	//
	//	pChar = r
	//}

	return string(result), nil
}
