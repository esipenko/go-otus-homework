package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	chars := []rune(str)

	//gr := uniseg.NewGraphemes(str)

	result := make([]rune, 0)
	var pChar rune

	//if unicode.IsDigit(chars[0]) {
	//	return "", ErrInvalidString
	//}

	for idx, r := range chars {
		if unicode.IsDigit(r) {
			amount, err := strconv.Atoi(string(r))

			if err != nil {
				return "", err
			}

			if idx == 0 {
				return "", ErrInvalidString
			}

			if unicode.IsDigit(pChar) {
				return "", ErrInvalidString
			}

			if amount == 0 {
				result = result[:len(result)-1]
				continue
			}

			for range amount - 1 {
				result = append(result, result[len(result)-1])
			}
		} else {
			result = append(result, r)
		}

		pChar = r
	}

	return string(result), nil
	// Iterate over grapheme clusters
	//for gr.Next() {
	//	r := gr.Runes()
	//
	//	if ()
	//	if len(r) > 1 {
	//		result = result
	//		fmt.Print(string(r[0]), " ")
	//
	//	}
	//	fmt.Println(gr.Runes())
	//}
	//return "", nil
}
