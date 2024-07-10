package customslices

import "reflect"

// собираем батчи сгруппированные по переданному полю
func GroupByElements(input interface{}, fieldName string) [][]interface{} {
	source := castSlice(input)

	if source == nil {
		return nil
	}

	groupMap := make(map[interface{}][]interface{})

	for _, elem := range source {
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

	var batches [][]interface{}
	for _, v := range groupMap {
		batches = append(batches, v)
	}

	return batches
}

func DivideSlice(source interface{}, batchSize int) [][]interface{} {
	var batches [][]interface{}

	input := castSlice(source)
	if input == nil {
		return nil
	}
	for i := 0; i < len(input); i += batchSize {
		end := i + batchSize
		if end > len(input) {
			end = len(input)
		}

		batch := make([]interface{}, end-i)
		for j := i; j < end; j++ {
			batch[j-i] = input[j]
		}

		batches = append(batches, batch)
	}

	return batches
}

// проверяем, является ли пустой интерфейс слайсом, и приводим к слайсу интерфейсов
func castSlice(input interface{}) []interface{} {
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Slice {
		return nil
	}

	// Создание нового слайса []interface{} с нужной длиной и емкостью
	result := make([]interface{}, v.Len())

	// Приведение каждого элемента слайса к interface{}
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		result[i] = elem.Interface()
	}

	return result
}
