package test

func comparePtrs(s1 *string, s2 *string) bool {
	if s1 == nil && s2 == nil {
		return true
	} else if (s1 == nil && s2 != nil) || (s1 != nil && s2 == nil) || (*s1 != *s2) {
		return false
	}
	return true
}

func compareArrPtr(s1 []*string, s2 []*string) bool {
	if s1 == nil && s2 == nil {
		return true
	} else if (s1 == nil && s2 != nil) || (s1 != nil && s2 == nil) || (len(s1) != len(s2)) {
		return false
	} else {
		size := len(s1)
		for i := 0; i < size; i++ {
			val1 := s1[i]
			val2 := s2[i]
			if !comparePtrs(val1, val2) {
				return false
			}
		}
		return true
	}
}
