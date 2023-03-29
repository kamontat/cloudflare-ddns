package utils

func ContainValue[
	T string |
		int | int32 | int64 |
		float32 | float64 |
		bool](arr []T, search T) bool {
	return ContainObject(arr, search, func(a, b T) bool {
		return a == b
	})
}

func MatchAllValue[
	T string |
		int | int32 | int64 |
		float32 | float64 |
		bool](arr []T, search T) bool {
	return MatchAllObject(arr, search, func(a, b T) bool {
		return a == b
	})
}

func ContainObject[T any](arr []T, search T, eq func(a, b T) bool) bool {
	for _, element := range arr {
		if eq(element, search) {
			return true
		}
	}
	return false
}

func MatchAllObject[T any](arr []T, search T, eq func(a, b T) bool) bool {
	for _, element := range arr {
		if !eq(element, search) {
			return false
		}
	}
	return true
}
