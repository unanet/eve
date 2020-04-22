package service

type m = map[string]interface{}

// mergeKeys stomps on the keys in the left map if they exist in the right map
func mergeKeys(left, right m) m {
	for key, rightVal := range right {
		left[key] = rightVal
	}
	return left
}
