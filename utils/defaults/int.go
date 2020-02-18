package defaults

func GetOrDefaultInt(existingValue *int, defaultValue int) int {
	if existingValue == nil || *existingValue == 0 {
		return defaultValue
	}

	return *existingValue
}
