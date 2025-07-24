package util

func StrToBool(str string) bool {
	return str == "true"
}

func HasDuplicate(a, b []string) bool {
	for _, itemA := range a {
		for _, itemB := range b {
			if itemA == itemB {
				return true
			}
		}
	}
	return false
}
