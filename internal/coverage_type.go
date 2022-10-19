package internal

import (
	"fmt"
)

// GetServerIndex returns serverIndex which corresponds to `ServerConfigurations`
func GetServerIndex(coverageType string) (int, error) {
	switch coverageType {
	case "jp":
		return 0, nil
	case "g":
		return 1, nil
	}
	return -1, fmt.Errorf("Invalid coverageType = %v specified", coverageType)
}
