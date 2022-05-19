package iufunc

func DoOrElse(condition bool, okValue interface{}, elseValue interface{}) interface{} {
	if condition {
		return okValue
	} else {
		return elseValue
	}
}
