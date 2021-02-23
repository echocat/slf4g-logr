package logr

import (
	"fmt"
)

func KeysAndValuesToFields(keysAndValues ...interface{}) map[string]interface{} {
	if len(keysAndValues)%2 != 0 {
		panic("illegal amount of arguments for keysAndValues provided;" +
			" expected always a value for a key, but one value seems to be missing")
	}
	l := len(keysAndValues)
	result := make(map[string]interface{}, l/2)
	for i := 0; i < l; i += 2 {
		if keysAndValues[i] == nil {
			panic(fmt.Sprintf("provided keyAndValue pair contains a nil key at index %d", i))
		}
		result[fmt.Sprint(keysAndValues[i])] = keysAndValues[i+1]
	}

	return result
}
