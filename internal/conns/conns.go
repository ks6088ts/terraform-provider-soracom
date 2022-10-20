package conns

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	soracom "github.com/ks6088ts/soracom-sdk-go/generated/api"
	"github.com/ks6088ts/terraform-provider-soracom/internal"
)

type Config struct {
	Profile string
}

type profile struct {
	Sandbox               bool    `json:"sandbox"`
	CoverageType          string  `json:"coverageType"`
	Email                 *string `json:"email,omitempty"`
	Password              *string `json:"password,omitempty"`
	AuthKeyId             *string `json:"authKeyId,omitempty"`
	AuthKey               *string `json:"authKey,omitempty"`
	Username              *string `json:"username,omitempty"`
	OperatorId            *string `json:"operatorId,omitempty"`
	Endpoint              *string `json:"endpoint,omitempty"`
	RegisterPaymentMethod bool    `json:"registerPaymentMethod"`
}

func getProfileDir() (string, error) {
	d := os.Getenv("SORACOM_PROFILE_DIR")
	if d != "" {
		return d, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".soracom"), nil
}

func getProfile(profileName string) (*profile, error) {

	d, err := getProfileDir()
	if err != nil {
		return nil, err
	}
	profilePath := filepath.Join(d, profileName+".json")

	b, err := os.ReadFile(profilePath)
	if err != nil {
		return nil, err
	}

	var p profile
	err = json.Unmarshal(b, &p)
	if err != nil {
		return nil, err
	}

	// supply default values for older versions (which support 'jp' coverage type only)
	if p.CoverageType == "" {
		p.CoverageType = "jp"
	}

	return &p, nil
}

type SoracomClient struct {
	Client       *soracom.APIClient
	AuthResponse *soracom.AuthResponse
	serverIndex  int
}

func (c *SoracomClient) GetContext(ctx context.Context) context.Context {
	contextApiKeys := map[string]soracom.APIKey{
		"api_key": {
			Key: *c.AuthResponse.ApiKey,
		},
		"api_token": {
			Key: *c.AuthResponse.Token,
		},
	}
	authCtx := context.WithValue(ctx, soracom.ContextAPIKeys, contextApiKeys)
	return context.WithValue(authCtx, soracom.ContextServerIndex, c.serverIndex)
}

func (c *Config) soracomClient(p *profile) (interface{}, error) {
	client := soracom.NewAPIClient(soracom.NewConfiguration())

	// call auth api to get api key and api token
	authRequest := soracom.AuthRequest{
		AuthKeyId: p.AuthKeyId,
		AuthKey:   p.AuthKey,
	}
	resp, r, err := client.AuthApi.Auth(context.Background()).AuthRequest(authRequest).Execute()
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, fmt.Errorf("/auth failed with status code = %v", r.StatusCode)
	}

	serverIndex, err := internal.GetServerIndex(p.CoverageType)
	if err != nil {
		return nil, fmt.Errorf("Failed to get serverIndex, err = %v", err)
	}

	return &SoracomClient{
		Client:       client,
		AuthResponse: resp,
		serverIndex:  serverIndex,
	}, nil
}

// Client configures and returns a fully initialized SoracomClient
func (c *Config) Client() (interface{}, error) {
	profile, err := getProfile(c.Profile)
	if err != nil {
		return nil, err
	}
	if profile.Sandbox {
		return nil, fmt.Errorf("Sandbox is not supported.")
	}
	return c.soracomClient(profile)
}
