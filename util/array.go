package util

// ArrayContains checks if array a[] contains value v or not
func ArrayContains[R comparable](a []R, v R) bool {
	for _, e := range a {
		if e == v {
			return true
		}
	}

	return false
}
