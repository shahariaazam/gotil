// Package guti contains packages
package guti

import (
	"encoding/csv"
	"math"
	"os"
	"reflect"
)

const epsilon = 1e-6

// IsExist searches for an item in a slice and returns true if it is found, and false otherwise.
// It supports searching for items of various types, including integers, floats, strings, booleans,
// and objects. It uses reflection to determine the type of the items in the slice, and to compare them
// to the search item. If the second argument is not a slice, it will panic. If the search item is not
// of the same type as the items in the slice, it will be skipped.
//
// Example usage:
//
//	intSlice := []int{1, 2, 3, 4, 5}
//	fmt.Println(guti.IsExist(3, intSlice)) // prints "true"
//	fmt.Println(guti.IsExist(6, intSlice)) // prints "false"
//
//	strSlice := []string{"foo", "bar", "baz"}
//	fmt.Println(guti.IsExist("qux", strSlice)) // prints "false"
//	fmt.Println(guti.IsExist("foo", strSlice)) // prints "true"
//
//	objectSlice := []struct {
//		Name string
//		Age  int
//	}{
//		{Name: "Alice", Age: 25},
//		{Name: "Bob", Age: 25},
//		{Name: "Charlie", Age: 35},
//	}
//	fmt.Println(guti.IsExist(struct {
//		Name string
//		Age  int
//	}{Name: "Bob", Age: 25}, objectSlice)) // prints "true"
//
//	boolSlice := []bool{true, false}
//	fmt.Println(guti.IsExist(true, boolSlice)) // prints "true"
//
//	emptySlice := []int{}
//	fmt.Println(guti.IsExist(1, emptySlice)) // prints "false"
//
// Playground: https://go.dev/play/p/jHua3iwd6xT
func IsExist(what interface{}, in interface{}) bool {
	s := reflect.ValueOf(in)

	if s.Kind() != reflect.Slice {
		panic("IsExist: Second argument must be a slice")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Kind() != reflect.TypeOf(what).Kind() {
			continue
		}

		switch s.Index(i).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if s.Index(i).Int() == reflect.ValueOf(what).Int() {
				return true
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if s.Index(i).Uint() == reflect.ValueOf(what).Uint() {
				return true
			}
		case reflect.Float32, reflect.Float64:
			if math.Abs(s.Index(i).Float()-reflect.ValueOf(what).Float()) < epsilon {
				return true
			}
		case reflect.String:
			if s.Index(i).String() == reflect.ValueOf(what).String() {
				return true
			}
		case reflect.Bool:
			if s.Index(i).Bool() == reflect.ValueOf(what).Bool() {
				return true
			}
		default:
			if reflect.DeepEqual(what, s.Index(i).Interface()) {
				return true
			}
		}
	}

	return false
}

// Filter returns a new list containing the elements of the input list that
// satisfy the given predicate function. The predicate function takes an input
// element of the list and returns true if the element should be included in the
// output list, and false otherwise. The input list can contain elements of any
// type, and the predicate function should take an argument of type interface{}.
//
// Example usage:
//
//	data := []interface{}{1, 2, 3, 4, 5}
//	isEven := func(x interface{}) bool { return x.(int)%2 == 0 }
//	result := Filter(data, isEven)  // result = []interface{}{2, 4}
//
// Playground: https://go.dev/play/p/haueBKmeb3e
func Filter(data []interface{}, predicate func(interface{}) bool) []interface{} {
	result := []interface{}{}
	for _, d := range data {
		if predicate(d) {
			result = append(result, d)
		}
	}
	return result
}

// Any returns true if at least one element of the input list satisfies the given predicate function,
// and false otherwise. The predicate function takes an input element of the list and returns true
// if the element satisfies the predicate, and false otherwise. The input list can contain elements
// of any type, and the predicate function should take an argument of type interface{}.
//
// Example usage:
//
//	data := []interface{}{1, 2, 3, 4, 5}
//	isEven := func(x interface{}) bool { return x.(int)%2 == 0 }
//	result := Any(data, isEven)  // result = true
//
// Playground: https://go.dev/play/p/mVzWG6tTp_2
func Any(data []interface{}, predicate func(interface{}) bool) bool {
	for _, d := range data {
		if predicate(d) {
			return true
		}
	}
	return false
}

// Reduce applies a reducing function to a list and returns a single value.
// The reducing function takes two arguments, an accumulator and a value, and returns
// a new accumulator. The initial value of the accumulator is provided as an argument.
// The function can reduce lists of any type, including integers, floats, strings,
// and custom types. If the initial value is not of the same type as the elements of
// the list, it will panic. The function returns the final value of the accumulator.
//
// Example usage:
//
//	data := []interface{}{1, 2, 3, 4, 5}
//
//	reduceFunc := func(acc interface{}, value interface{}) interface{} {
//		return acc.(int) + value.(int)
//	}
//
//	initial := 0
//	result := guti.Reduce(data, reduceFunc, initial)
//	fmt.Println(result) // should print 15
//
// Playground: https://go.dev/play/p/A7ZQrVp_uIk
func Reduce(data []interface{}, reduce func(interface{}, interface{}) interface{}, initial interface{}) interface{} {
	acc := initial
	for _, d := range data {
		acc = reduce(acc, d)
	}
	return acc
}

// Map applies a transformation function to each element of a slice and returns a new slice with the
// transformed elements. The transform function takes an element of the input slice as input and returns
// a transformed value. The input slice can contain elements of any type, but the transform function must
// be able to handle each element type appropriately. The returned slice has the same length as the input
// slice, and each element is the result of applying the transform function to the corresponding input element.
// The input slice is not modified by the function.
//
// Example usage:
//
//	input := []interface{}{1, 2, 3, 4, 5}
//
//	transform := func(d interface{}) interface{} {
//		return d.(int) * 2
//	}
//	output := Map(input, transform)
//	fmt.Println(output) // [2 4 6 8 10]
//
// Playground: https://go.dev/play/p/ZguMfToP0Xh
func Map(data []interface{}, transform func(interface{}) interface{}) []interface{} {
	result := []interface{}{}
	for _, d := range data {
		result = append(result, transform(d))
	}
	return result
}

// IndexOf returns the index of the first occurrence of a given element in a list. If the element is not found, it returns -1.
// The data parameter is a slice of interface{} type which can hold any type of data. The element parameter is the element whose index is to be searched in the slice.
// This function returns an integer value that represents the index of the first occurrence of the given element in the slice.
//
// Example usage:
//
//	data := []interface{}{"apple", "banana", "cherry"}
//
//	element := "banana"
//
//	index := guti.IndexOf(data, element)
//	fmt.Println("Index of", element, "is", index)	// should output: Index of banana is 1
//
// Playground: https://go.dev/play/p/K7X-4_RbJPG
func IndexOf(data []interface{}, element interface{}) int {
	for i, d := range data {
		if d == element {
			return i
		}
	}
	return -1
}

// ContainsAll returns true if all elements in the first slice are present in the second slice, otherwise returns false.
func ContainsAll(s1, s2 []interface{}) bool {
	for _, e1 := range s1 {
		found := false
		for _, e2 := range s2 {
			if e1 == e2 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Reverse returns a new slice with the elements of the given slice in reverse order.
func Reverse(slice []interface{}) []interface{} {
	result := make([]interface{}, len(slice))
	for i, j := 0, len(slice)-1; i <= j; i, j = i+1, j-1 {
		result[i], result[j] = slice[j], slice[i]
	}
	return result
}

// FilterNil returns a new slice with all nil elements removed from the given slice.
func FilterNil(slice []interface{}) []interface{} {
	result := make([]interface{}, 0, len(slice))
	for _, v := range slice {
		if v != nil {
			result = append(result, v)
		}
	}
	return result
}

// MapReduce takes a slice of items and applies a mapper function to each item to get a slice of results. It then applies a reducer function to the slice of results to get a single result
func MapReduce(items interface{}, mapper func(interface{}) interface{}, reducer func(interface{}, interface{}) interface{}) interface{} {
	mappedItems := make([]interface{}, 0)
	itemsValue := reflect.ValueOf(items)

	for i := 0; i < itemsValue.Len(); i++ {
		mappedItems = append(mappedItems, mapper(itemsValue.Index(i).Interface()))
	}

	reducedResult := mappedItems[0]
	for i := 1; i < len(mappedItems); i++ {
		reducedResult = reducer(reducedResult, mappedItems[i])
	}

	return reducedResult
}

// Batch takes a slice of items and a batch size, and returns a slice of slices, where each inner slice contains at most batchSize items from the input slice
func Batch(items interface{}, batchSize int) [][]interface{} {
	var batches [][]interface{}
	itemsValue := reflect.ValueOf(items)
	batchSize = int(math.Min(float64(batchSize), float64(itemsValue.Len())))

	for i := 0; i < itemsValue.Len(); i += batchSize {
		end := int(math.Min(float64(i+batchSize), float64(itemsValue.Len())))
		batches = append(batches, ConvertSliceInterfaceToSlice(itemsValue.Slice(i, end)))
	}

	return batches
}

// ConvertSliceInterfaceToSlice takes a reflect.Value of a slice of unknown type and returns a new slice of interface{} type
func ConvertSliceInterfaceToSlice(slice reflect.Value) []interface{} {
	s := make([]interface{}, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		s[i] = slice.Index(i).Interface()
	}
	return s
}

// SaveAsCSV save data to csv
func SaveAsCSV(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	// Get the type of the data slice and write the header row
	dataType := reflect.TypeOf(data).Elem()
	headerRow := make([]string, dataType.NumField())
	for i := 0; i < dataType.NumField(); i++ {
		headerRow[i] = dataType.Field(i).Name
	}
	writer.Write(headerRow)

	// Write each row of data to the CSV file
	dataValue := reflect.ValueOf(data)
	for i := 0; i < dataValue.Len(); i++ {
		row := make([]string, dataType.NumField())
		for j := 0; j < dataType.NumField(); j++ {
			fieldValue := dataValue.Index(i).Field(j)
			row[j] = fieldValue.Interface().(string)
		}
		writer.Write(row)
	}

	writer.Flush()

	return nil
}
