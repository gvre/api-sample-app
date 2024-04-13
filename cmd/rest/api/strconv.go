package api

import (
	"strconv"
)

// toInt returns the int value of the variable or an error if value is not a valid integer.
func toInt(val string) (int, error) {
	ival, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return ival, nil
}
