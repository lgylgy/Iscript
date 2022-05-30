package s1

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
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
	return strings.SplitN(value, "", n)
}
