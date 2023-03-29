package utils

import "fmt"

func GetOr[T any](name string, value *T, def *T, final T) (T, error) {
	if value == nil && def == nil {
		return final, fmt.Errorf("cannot get '%s'", name)
	}
	if value != nil {
		return *value, nil
	}
	return *def, nil
}

func GetOrElse[T any](value *T, def T) T {
	if value == nil {
		return def
	}

	return *value
}
