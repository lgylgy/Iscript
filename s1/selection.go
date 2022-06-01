package s1

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

func Select(values []string, number, seed int) ([]string, error) {
	if len(values) < number {
		return nil, fmt.Errorf("not enough values")
	}
	sort.Strings(values)
	rand.Seed(int64(seed))
	rand.Shuffle(len(values), func(i, j int) {
		values[i], values[j] = values[j], values[i]
	})
	return values[:number], nil
}

func Split(value string, n int) []string {
	result := []string{}
	length := len(value)
	offset := int(math.Ceil(float64(length) / float64(n)))
	for i := 0; i < length; {
		j := i + offset
		if j > length {
			j = length
		}
		result = append(result, value[i:j])
		i = j
	}
	return result
}
