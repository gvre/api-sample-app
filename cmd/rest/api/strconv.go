package api

import (
	"fmt"
	"strconv"
)

// toInt returns the int value of the variable, if it exists in the provided map.
// If it does not or it is not a valid int, an error is returned.
func toInt(vars map[string]string, v string) (int, error) {
	sval, ok := vars[v]
	if !ok {
		return 0, fmt.Errorf("variable %q was not found", v)
	}

	ival, err := strconv.Atoi(sval)
	if err != nil {
		return 0, err
	}

	return ival, nil
}
