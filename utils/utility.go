package utils

import "strings"

func IndexInArray(arr []string, item string) int {
	for index, it := range arr {
		if strings.HasPrefix(it, item) {
			return index
		}
	}
	return -1
}
