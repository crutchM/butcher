package customslices

import "reflect"

// собираем батчи сгруппированные по переданному полю
func GroupByElements[T any](input []T, fieldName string) [][]T {

	groupMap := make(map[interface{}][]T)

	for _, elem := range input {
		elemVal := reflect.ValueOf(elem)
		if elemVal.Kind() != reflect.Struct {
			return nil
		}
		field := elemVal.FieldByName(fieldName)
		if !field.IsValid() {
			return nil
		}
		fieldValue := field.Interface()
		groupMap[fieldValue] = append(groupMap[fieldValue], elem)
	}

	var batches [][]T
	for _, v := range groupMap {
		batches = append(batches, v)
	}

	return batches
}

func DivideSlice[T any](input []T, batchSize int) [][]T {
	var batches [][]T

	for i := 0; i < len(input); i += batchSize {
		end := i + batchSize
		if end > len(input) {
			end = len(input)
		}

		batch := make([]T, end-i)
		for j := i; j < end; j++ {
			batch[j-i] = input[j]
		}

		batches = append(batches, batch)
	}

	return batches
}
