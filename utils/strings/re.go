package strings

import (
	"fmt"
	"regexp"
)

func GetRegexGroups(regEx, data string) (paramsMap map[string]string) {
	/* Credits: https://stackoverflow.com/a/39635221 */
	paramsMap = make(map[string]string)

	compRegEx, err := regexp.Compile(regEx)
	if err != nil {
		return paramsMap
	}

	match := compRegEx.FindStringSubmatch(data)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}

	return
}

func GetRegexGroup(regEx, data string, group string) (string, error) {
	/* Credits: https://stackoverflow.com/a/39635221 */
	reGroups := GetRegexGroups(regEx, data)
	if rData, ok := reGroups[group]; ok {
		return rData, nil
	}

	return "", fmt.Errorf("group %q was not found", group)
}
