package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

var (
	ErrUnsupportedValidation = errors.New("unsupported validation")
	ErrMinError              = errors.New("value should be greater than")
	ErrMaxError              = errors.New("value should be less than")
	ErrInError               = errors.New("value is not in set")
	ErrLenError              = errors.New("length should be equal to ")
	ErrRegexpError           = errors.New("string must match pattern")
	ErrUnsupportedType       = errors.New("unsupported type")
)

func validateIntField(field int, ruleName, ruleValue string, fieldName string) (*ValidationError, error) {
	if ruleName == "len" || ruleName == "regexp" {
		return &ValidationError{Field: fieldName, Err: fmt.Errorf("%w", ErrUnsupportedValidation)}, nil
	}

	if ruleName == "min" {
		minValue, err := strconv.Atoi(ruleValue)
		if err != nil {
			return nil, err
		}

		if field < minValue {
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("%w %d", ErrMinError, minValue)}, nil
		}
	}

	if ruleName == "max" {
		maxValue, err := strconv.Atoi(ruleValue)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", ruleName, err)
		}

		if field > maxValue {
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("%w %d", ErrMaxError, maxValue)}, nil
		}
	}

	if ruleName == "in" {
		inValues := strings.Split(ruleValue, ",")

		for _, inValue := range inValues {
			value, err := strconv.Atoi(inValue)
			if err != nil {
				return nil, err
			}

			if value == field {
				return nil, nil
			}
		}

		return &ValidationError{Field: fieldName, Err: fmt.Errorf("%w %s", ErrInError, ruleValue)}, nil
	}

	return nil, nil
}

func validateStringField(field, ruleName, ruleValue string, fieldName string) (*ValidationError, error) {
	if ruleName == "min" || ruleName == "max" {
		return &ValidationError{Field: fieldName, Err: fmt.Errorf("%w %s", ErrUnsupportedValidation, ruleName)}, nil
	}

	if ruleName == "len" {
		requiredLength, err := strconv.Atoi(ruleValue)
		if err != nil {
			return nil, err
		}

		if requiredLength != len(field) {
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("%w %d", ErrLenError, requiredLength)}, nil
		}
	}

	if ruleName == "in" {
		inValues := strings.Split(ruleValue, ",")
		for _, inValue := range inValues {
			if inValue == field {
				return nil, nil
			}
		}

		return &ValidationError{Field: fieldName, Err: fmt.Errorf("%w %s", ErrInError, ruleValue)}, nil
	}

	if ruleName == "regexp" {
		reg, err := regexp.Compile(ruleValue)
		if err != nil {
			return nil, err
		}

		if !reg.MatchString(field) {
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("%w %s", ErrRegexpError, ruleValue)}, nil
		}

		return nil, nil
	}

	return nil, nil
}

func validateFiled(rule string, field reflect.Value, fieldName string) ([]ValidationError, error) {
	parts := strings.Split(rule, ":")
	ruleName, ruleValue := parts[0], parts[1]

	kind := field.Type().Kind()
	// Не понимаю проблему с линтером. Говорит не обработаны остальные кейсы
	// хотя дефолт есть. Попытался переписать на иф елз, линтер ругается на то что
	// надо вернуть свитч
	//nolint:exhaustive
	switch kind {
	case reflect.Int:
		val, ok := field.Interface().(int)

		if !ok {
			return nil, fmt.Errorf("invalid int")
		}

		validationError, systemError := validateIntField(val, ruleName, ruleValue, fieldName)

		if systemError != nil {
			return nil, systemError
		}

		var result []ValidationError
		if validationError != nil {
			return append(result, *validationError), nil
		}

		return nil, nil
	case reflect.String:
		val, ok := field.Interface().(string)
		if !ok {
			return nil, fmt.Errorf("invalid string")
		}
		validationError, systemError := validateStringField(val, ruleName, ruleValue, fieldName)
		if systemError != nil {
			return nil, systemError
		}

		var result []ValidationError
		if validationError != nil {
			return append(result, *validationError), nil
		}
		return nil, nil
	case reflect.Slice:
		elemType := field.Type().Elem()
		if elemType.Kind() != reflect.String && elemType.Kind() != reflect.Int {
			return nil, ErrUnsupportedType
		}

		var result []ValidationError
		for idx := 0; idx < field.Len(); idx++ {
			elem := field.Index(idx)
			validationErrors, systemError := validateFiled(rule, elem, fieldName+"["+strconv.Itoa(idx)+"]")
			if systemError != nil {
				return nil, systemError
			}
			result = append(result, validationErrors...)
		}
		return result, nil
	default:
		return nil, ErrUnsupportedType
	}
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var result string
	for _, err := range v {
		result += fmt.Sprintf("%s: %s\n", err.Field, err.Err.Error())
	}

	return result
}

func (v ValidationErrors) len() int {
	return len(v)
}

func Validate(v interface{}) error {
	st := reflect.TypeOf(v)

	if st.Kind() != reflect.Struct {
		return ErrUnsupportedType
	}

	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	result := ValidationErrors{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		validationStr := fieldType.Tag.Get("validate")
		if validationStr == "" {
			continue
		}

		rules := strings.Split(validationStr, "|")

		for _, rule := range rules {
			validationErrors, systemError := validateFiled(rule, field, fieldType.Name)
			if systemError != nil {
				return systemError
			}

			result = append(result, validationErrors...)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
