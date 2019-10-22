package maps

/* Public */

func MergeStringMap(primaryMap map[string]string, mergeMap map[string]string) {
	for k, v := range mergeMap {
		primaryMap[k] = v
	}
}
