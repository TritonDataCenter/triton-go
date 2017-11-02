package seq

import "fmt"

func Test(expected, actual interface{}) *Result {

	eMap := flatten("", objectToMap(expected))
	aMap := flatten("", objectToMap(actual))
	result := diff(eMap, aMap)
	return result
}

type Map map[string]interface{}

func (m Map) Test(actual interface{}) *Result {
	return Test(m, actual)
}

func debug(m map[string]string) {
	for k, v := range m {
		fmt.Println(k, ":", v)
	}
}
