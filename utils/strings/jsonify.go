package strings

import (
	"encoding/json"
	"fmt"
)

func JsonifyStrict(obj interface{}) (string, error) {
	bs, err := json.MarshalIndent(obj, "", "  ")
	return string(bs), err
}

func JsonifyLax(obj interface{}) string {
	bs, err := json.Marshal(obj)
	if err != nil {
		return fmt.Sprintf("%+v", obj)
	}
	return string(bs)
}
