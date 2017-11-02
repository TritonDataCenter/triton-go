package seq

import (
	"encoding/json"
	"fmt"
)

func unmarshal(d []byte, i interface{}) {
	err := json.Unmarshal(d, i)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal '%s': %s", string(d), err))
	}
}
func marshal(i interface{}) []byte {
	b, err := json.Marshal(i)
	if err != nil {
		panic("Failed to marshal")
	}
	return b
}
func marshalIndent(i interface{}) []byte {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		panic("Failed to marshal")
	}
	return b
}
func isSlice(object interface{}) bool {
	t := fmt.Sprintf("%T", object)
	return t[0] == '['
}
func objectToMap(object interface{}) interface{} {

	b := marshal(object)

	if isSlice(object) {

		var out []interface{}
		unmarshal(b, &out)
		return out
	} else {
		var out map[string]interface{}
		unmarshal(b, &out)
		return out
	}
}
