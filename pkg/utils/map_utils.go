package utils

// GetMapKeys get map keys
func GetMapKeys(data map[string]interface{}) []string {
	var keys = make([]string, 0)
	for key := range data {
		keys = append(keys, key)
	}
	return keys
}

// MergeMap merge map
func MergeMap(mObj ...map[string]interface{}) map[string]interface{} {
	newObj := map[string]interface{}{}
	for _, m := range mObj {
		if m != nil {
			for k, v := range m {
				newObj[k] = v
			}
		}
	}
	return newObj
}
