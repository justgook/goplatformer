package util

func PointerSliceClone[S ~[]*E, E any](s S) S {
	output := make([]*E, len(s))
	for j := range s {
		aaa := *s[j]
		output[j] = &aaa
	}

	return output
}
