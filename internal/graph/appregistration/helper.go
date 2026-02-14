package appregistration

func applyString(value *string, setter func(*string)) bool {
	if value == nil {
		return false
	}

	setter(value)
	return true
}

func applyBool(value *bool, setter func(*bool)) bool {
	if value == nil {
		return false
	}

	setter(value)
	return true
}

func applyStringSlice(value *[]string, setter func([]string)) bool {

	if value == nil {
		return false
	}

	setter(*value)
	return true
}
