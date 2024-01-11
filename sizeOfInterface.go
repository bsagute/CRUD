package main

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/dustin/go-humanize"
)

func deepSizeOf(v interface{}) (uintptr, error) {
	seen := make(map[uintptr]bool)
	return deepSizeOfValue(reflect.ValueOf(v), seen)
}

func deepSizeOfValue(v reflect.Value, seen map[uintptr]bool) (uintptr, error) {
	if v.CanAddr() {
		addr := v.UnsafeAddr()
		if seen[addr] {
			return 0, nil
		}
		seen[addr] = true
	}

	size := v.Type().Size()

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			elemSize, err := deepSizeOfValue(v.Index(i), seen)
			if err != nil {
				return 0, err
			}
			size += elemSize
		}
	case reflect.Map:
		keys := v.MapKeys()
		for _, key := range keys {
			keySize, err := deepSizeOfValue(key, seen)
			if err != nil {
				return 0, err
			}
			size += keySize

			valueSize, err := deepSizeOfValue(v.MapIndex(key), seen)
			if err != nil {
				return 0, err
			}
			size += valueSize
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldSize, err := deepSizeOfValue(v.Field(i), seen)
			if err != nil {
				return 0, err
			}
			size += fieldSize
		}
	}

	return size, nil
}

func main() {
	myMap := make(map[string]interface{})
	myMap["key1"] = "value1"
	myMap["key2"] = 42
	myMap["key3"] = true

	mapSize, err := deepSizeOf(myMap)
	if err != nil {
		fmt.Printf("Error calculating map size: %v\n", err)
		return
	}

	mapSizeInKb := float64(mapSize) / 1024.0

	fmt.Printf("Size of the map: %s\n", humanize.Bytes(mapSize))
	fmt.Printf("Size of the map: %.2f KB\n", mapSizeInKb)
}
