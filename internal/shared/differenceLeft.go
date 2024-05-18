package shared

func DifferenceLeft[T comparable](list1 []T, list2 []T) []T {
	var left []T
	seenRight := map[T]struct{}{}

	for _, elem := range list2 {
		seenRight[elem] = struct{}{}
	}

	for _, elem := range list1 {
		if _, ok := seenRight[elem]; !ok {
			left = append(left, elem)
		}
	}

	return left
}
