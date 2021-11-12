package testItems

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func RandomInt(min, max int) int {
	if min == max {
		return min
	} else if min > max {
		min, max = max, min
	}

	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func RandomFloat64(min, max float64) float64 {
	if min == max {
		return min
	} else if min > max {
		min, max = max, min
	}

	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}

func RandomSliceFloat64(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}

func lengthAfterDot(input string) (length int) {
	index := strings.IndexByte(input, '.')
	if index > -1 {
		return len(input) - index - 1
	}
	return 0
}

func ConvertInterfaceToFloat64(input interface{}) (output float64, numDecPlaces int, isFloat64 bool) {

	switch value := input.(type) {

	case float32:
		return float64(value), lengthAfterDot(fmt.Sprint(input)), true
	case float64:
		return value, lengthAfterDot(fmt.Sprint(input)), true
	case string:
		output, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0.0, 0, false
		}
		return output, lengthAfterDot(fmt.Sprint(output)), true
	case int:
		return float64(value), 0, true
	case int16:
		return float64(value), 0, true
	case int32:
		return float64(value), 0, true
	case int64:
		return float64(value), 0, true
	case int8:
		return float64(value), 0, true
	case uint:
		return float64(value), 0, true
	case uint16:
		return float64(value), 0, true
	case uint32:
		return float64(value), 0, true
	case uint64:
		return float64(value), 0, true
	case uint8:
		return float64(value), 0, true
	default:
		return 0.0, 0, false
	}
}
