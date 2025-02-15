package extslices

// filter given any slice by the given filter function
func Filter[T any](slice []T, filter func(T) bool) []T {
	var filtered []T
	for _, v := range slice {
		if filter(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// map given any slice by the given map function
func Map[T, R any](slice []T, mapper func(T) R) []R {
	mapped := make([]R, len(slice))
	for i, v := range slice {
		mapped[i] = mapper(v)
	}
	return mapped
}
