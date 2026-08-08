package cmd

import (
	"fmt"
	"strconv"
)

func parseOptionalInt(args []string, defaultValue, minValue int, name string) (int, error) {
	if len(args) == 0 {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(args[0])
	if err != nil || value < minValue {
		requirement := "非负整数"
		if minValue > 0 {
			requirement = "正整数"
		}
		return 0, fmt.Errorf("%s 必须是%s", name, requirement)
	}
	return value, nil
}
