package utils

import (
	"fmt"
	"log"
	"strconv"
)

// ByteCountSI converts input string like "982814103" (size in bytes, number represented as a string)
// to a more human-readable string like "969.9 MB"
func ByteCountSI(input string) string {

	b, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(err)
	}

	return ByteCountSIInt(b)
}

// ByteCountSIInt converts input int like 982814103 (size in bytes)
// to a more human-readable string like "969.9 MB"
func ByteCountSIInt(input int) string {

	const unit = 1000
	if input < unit {
		return fmt.Sprintf("%d B", input)
	}
	div, exp := int64(unit), 0
	for n := input / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(input)/float64(div), "kMGTPE"[exp])
}
