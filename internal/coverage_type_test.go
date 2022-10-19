package internal

import (
	"testing"
)

func TestClient(t *testing.T) {
	testCases := []struct {
		Name                string
		CoverageType        string
		ExpectedServerIndex int
		HasError            bool
	}{
		{
			Name:                "nominal scenarios for Global coverage",
			CoverageType:        "g",
			ExpectedServerIndex: 1,
			HasError:            false,
		},
		{
			Name:                "nominal scenarios for Japan coverage",
			CoverageType:        "jp",
			ExpectedServerIndex: 0,
			HasError:            false,
		},
		{
			Name:                "non-nominal scenarios for unsupported coverage",
			CoverageType:        "n",
			ExpectedServerIndex: -1,
			HasError:            true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			serverIndex, err := GetServerIndex(testCase.CoverageType)
			if hasError := (err != nil); hasError != testCase.HasError {
				t.Errorf("got %v, expected %v, err %v", hasError, testCase.HasError, err)
			}
			if testCase.ExpectedServerIndex != serverIndex {
				t.Errorf("got %v, expected %v", serverIndex, testCase.ExpectedServerIndex)
			}
		})
	}
}
