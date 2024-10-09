package checker

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Checker interface {
	Check(strct any) error
}

// WithCheck Добавляет функционал для проверки заполнения структуры юзкейса требуемыми хэндлерами.
type WithCheck struct{}

var ErrInitError = errors.New("запуск остановлен с ошибкой")

// Check Инициирует проверку заполнения структуры юзкейса.
// Собирает перечень незаполненных хэндлеров и выдаёт их перечень в тексте ошибке.
func (ucc WithCheck) Check(checkingStruct any) error {
	emptyParams := make([]string, 0)
	v := reflect.ValueOf(checkingStruct)
	t := reflect.TypeOf(checkingStruct)
	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Tag.Get("check") == "optional" {
			continue
		}

		if (v.Field(i).Kind() == reflect.Pointer ||
			v.Field(i).Kind() == reflect.Interface) && v.Field(i).IsNil() {
			emptyParams = append(emptyParams, v.Field(i).String())
		}
	}

	if len(emptyParams) == 0 {
		return nil
	}

	return fmt.Errorf(
		"%s %w. обнаружены непроинициированные обязательные свойства структуры: %s",
		v.Type().String(),
		ErrInitError,
		strings.Join(emptyParams, ", "),
	)
}
