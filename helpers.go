package utils

import (
	"digi-model-engine/utils/exceptions"
	"encoding/json"
	"math"
	"reflect"
	"runtime"
	"strings"
)

func UnmarshalJSON(data []byte, target any) error {
	if err := json.Unmarshal(data, target); err != nil {
		exceptions.InternalServerError(err)
	}
	return nil
}

// retrieves the filename where a function is defined.
func GetFilename(function interface{}) string {
	funcPtr := reflect.ValueOf(function).Pointer()
	funcName := runtime.FuncForPC(funcPtr).Name()

	// Split the handlerName to get the package and function name
	parts := strings.Split(funcName, ".")

	// The last part contains the function name and filename
	functionAndFile := parts[len(parts)-1]

	// Split again to separate function name and filename
	parts = strings.Split(functionAndFile, "/")

	// The last part is the filename
	return parts[len(parts)-1]
}

func StructToMap(u1 any) (m map[string]interface{}) {
	b, _ := json.Marshal(&u1)
	_ = json.Unmarshal(b, &m)
	return
}

func RemoveDuplicates[T string | int | float64](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// Utility function to find the maximum value in a float64 slice.
func MaxFloat64(slice []float64) float64 {
	max := slice[0]
	for _, val := range slice {
		if val > max {
			max = val
		}
	}
	return max
}

// Utility function to find the maximum value in an int slice.
func MaxInt(slice []int) int {
	max := slice[0]
	for _, val := range slice {
		if val > max {
			max = val
		}
	}
	return max
}

func CalculateTax(annualSalary float64, taxScheme string) float64 {
	if annualSalary <= 250000 {
		return 0.0
	}
	if annualSalary <= 500000 {
		exceeding := annualSalary - 250000
		return 0.05 * exceeding
	}
	if annualSalary <= 750000 {
		percentage := 0.0
		exceeding := annualSalary - 500000
		if taxScheme == "new" {
			percentage = 10.0
		} else {
			percentage = 20.0
		}
		return 12500 + ((percentage) / 100 * exceeding)
	}
	if annualSalary <= 1000000 {
		exceeding := annualSalary - 750000
		if taxScheme == "new" {
			return 37500.0 + ((15.0) / 100 * exceeding)
		} else {
			return 62500.0 + ((20.0) / 100 * exceeding)
		}
	}
	if annualSalary <= 1250000 {
		exceeding := annualSalary - 1000000
		if taxScheme == "new" {
			return 75_000.0 + ((20.0) / 100 * exceeding)
		} else {
			return 1_12500.0 + ((30.0) / 100 * exceeding)
		}
	}
	if annualSalary <= 1500000 {
		exceeding := annualSalary - 1250000
		if taxScheme == "new" {
			return 125000.0 + ((25.0) / 100 * exceeding)
		} else {
			return 187500.0 + ((30.0) / 100 * exceeding)
		}
	}
	exceeding := annualSalary - 15_00000
	if taxScheme == "new" {
		return 1_87500.0 + ((30.0) / 100 * exceeding)
	} else {
		return 2_62500.0 + ((30.0) / 100 * exceeding)
	}
}

func Round(number, decimalPlaces float64) float64 {
	mulFactor := math.Pow(10, decimalPlaces)
	return math.Round(number*mulFactor) / mulFactor
}
