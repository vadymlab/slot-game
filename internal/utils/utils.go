package utils

import (
	"github.com/urfave/cli/v2"
)

// MergeSlices combines multiple slices of CLI flags into a single slice.
// It calculates the total length of the resulting slice to avoid reallocation during appending.
func MergeSlices(slices ...[]cli.Flag) []cli.Flag {
	totalLength := 0
	for _, slice := range slices {
		totalLength += len(slice)
	}

	result := make([]cli.Flag, 0, totalLength)

	for _, slice := range slices {
		result = append(result, slice...)
	}

	return result
}
