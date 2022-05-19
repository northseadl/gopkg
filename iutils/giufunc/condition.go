package giufunc

func OkOrElse[T any](condition bool, okValue T, elseValue T) T {
	if condition {
		return okValue
	} else {
		return elseValue
	}
}
