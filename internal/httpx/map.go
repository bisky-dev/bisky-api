package httpx

func MapSliceE[T any, R any](items []T, mapper func(T) (R, error)) ([]R, error) {
	result := make([]R, 0, len(items))
	for _, item := range items {
		mapped, err := mapper(item)
		if err != nil {
			return nil, err
		}
		result = append(result, mapped)
	}
	return result, nil
}
