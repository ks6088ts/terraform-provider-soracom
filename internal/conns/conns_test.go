package conns

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClient(t *testing.T) {
	testCases := []struct {
		Name     string
		Config   Config
		HasError bool
	}{
		{
			Name: "nominal scenarios",
			Config: Config{
				Profile: "default",
			},
			HasError: false,
		},
		{
			Name: "non-nominal scenarios with invalid Profile",
			Config: Config{
				Profile: "invalid",
			},
			HasError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := testCase.Config.Client()
			if hasError := (err != nil); hasError != testCase.HasError {
				t.Errorf("got %v, expected %v, err %v", hasError, testCase.HasError, err)
			}
		})
	}
}

func TestGetProfileDir(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("failed to get HOME dir")
	}

	testCases := []struct {
		name              string
		soracomProfileDir string
		want              string
		hasError          bool
	}{
		{
			name:              "nominal scenarios with SORACOM_PROFILE_DIR defined",
			soracomProfileDir: "/test",
			want:              "/test",
			hasError:          false,
		},
		{
			name:              "nominal scenarios without SORACOM_PROFILE_DIR",
			soracomProfileDir: "",
			want:              filepath.Join(homeDir, ".soracom"),
			hasError:          false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := os.Setenv("SORACOM_PROFILE_DIR", testCase.soracomProfileDir)
			if err != nil {
				t.Errorf("failed to set environment variable")
			}

			got, err := getProfileDir()
			if hasError := (err != nil); hasError != testCase.hasError {
				t.Errorf("got %v, expected %v, err %v", hasError, testCase.hasError, err)
			}
			if got != testCase.want {
				t.Errorf("want = %v, but got = %v", testCase.want, got)
			}
		})
	}
}

func TestGetProfile(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("failed to get working directory, err=%v", err)
	}

	if err := os.Setenv("SORACOM_PROFILE_DIR", filepath.Join(wd, "..", "..", "test")); err != nil {
		t.Errorf("failed to set environment variable")
	}

	want := profile{
		Sandbox:               false,
		CoverageType:          "jp",
		RegisterPaymentMethod: false,
	}

	got, err := getProfile("profile")
	if err != nil {
		t.Errorf("failed to get profile, err=%v, got=%v", err, got)
	}

	if got.Sandbox != want.Sandbox {
		t.Errorf("failed to get sandbox")
	}
	if got.CoverageType != want.CoverageType {
		t.Errorf("failed to get coverage type")
	}
	if got.RegisterPaymentMethod != want.RegisterPaymentMethod {
		t.Errorf("failed to get registerPaymentMethod")
	}
}
