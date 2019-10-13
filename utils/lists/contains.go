package lists

import "strings"

/* Public */

func StringListContains(list []string, key string, caseSensitive bool) bool {
	for _, listKey := range list {
		switch caseSensitive {
		case false:
			if strings.ToLower(listKey) == strings.ToLower(key) {
				return true
			}
		default:
			if listKey == key {
				return true
			}
		}
	}

	return false
}
