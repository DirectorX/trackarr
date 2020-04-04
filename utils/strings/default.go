package strings

func NewOrExisting(new *string, existing string) string {
	if new == nil || *new == "" {
		return existing
	}

	return *new
}
