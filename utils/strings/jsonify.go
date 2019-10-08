package strings

import "encoding/json"

func Jsonify(obj interface{}) (string, error) {
	bs, err := json.MarshalIndent(obj, "", "  ")
	return string(bs), err
}
