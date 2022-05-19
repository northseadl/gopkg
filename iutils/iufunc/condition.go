package iufunc

func OkOrElse(condition bool, okValue interface{}, elseValue interface{}) interface{} {
	if condition {
		return okValue
	} else {
		return elseValue
	}
}
