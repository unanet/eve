package service

type StringList []string

func (s StringList) Contains(value string) bool {
	for _, a := range s {
		if a == value {
			return true
		}
	}
	return false
}

type M map[string]interface{}

// mergeKeys stomps on the keys in the left map if they exist in the right map
func mergeKeys(left, right M) M {
	for key, rightVal := range right {
		left[key] = rightVal
	}
	return left
}
