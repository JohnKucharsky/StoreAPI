package shared

func IntersectUniq[T comparable](list1 []T, list2 []T) []T {
	var result []T
	seen := map[T]struct{}{}

	for _, elem := range list2 {
		seen[elem] = struct{}{}
	}

	for _, elem := range list1 {
		if _, ok := seen[elem]; ok {
			result = append(result, elem)
		}
	}

	return result
}
