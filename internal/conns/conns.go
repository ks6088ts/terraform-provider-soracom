package conns

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	soracom "github.com/ks6088ts/soracom-sdk-go/generated/api"
	"github.com/ks6088ts/soracom-sdk-go/generated/sandbox"
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

	if err := validateProfile(&p); err != nil {
		return nil, err
	}

	// supply default values for older versions (which support 'jp' coverage type only)
	if p.CoverageType == "" {
		p.CoverageType = "jp"
	}

	return &p, nil
}

func validateProfile(p *profile) error {
	// validation logic: handles authentication methods other than auth key
	if p.AuthKey == nil || p.AuthKeyId == nil {
		return fmt.Errorf("authentication by auth key is only supported. please specify AuthKey and AuthKeyId in profile")
	}
	return nil
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

func (c *Config) sandboxClient(p *profile) (interface{}, error) {
	randStr := strconv.FormatInt(time.Now().UnixNano(), 10)
	sandboxInitRequest := *sandbox.NewSandboxInitRequest(*p.AuthKey, *p.AuthKeyId, fmt.Sprintf("sdk-test+%s@soracom.jp", randStr), fmt.Sprintf("Password%s", randStr))
	sandboxInitRequest.CoverageTypes = append(sandboxInitRequest.CoverageTypes, "g", "jp")
	configuration := sandbox.NewConfiguration()
	client := sandbox.NewAPIClient(configuration)
	resp, r, err := client.OperatorApi.SandboxInitializeOperator(context.Background()).SandboxInitRequest(sandboxInitRequest).Execute()
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 201 {
		return nil, fmt.Errorf("Invalid HTTP response: %v", r.StatusCode)
	}

	soracomConfiguration := soracom.NewConfiguration()
	// fixme: adhoc impl to overwrite endpoint url
	soracomConfiguration.Servers = soracom.ServerConfigurations{
		soracom.ServerConfiguration{
			URL: sandbox.NewConfiguration().Servers[0].URL, // fixme: verify index access for safety
		},
	}

	return &SoracomClient{
		Client: soracom.NewAPIClient(soracomConfiguration),
		AuthResponse: &soracom.AuthResponse{
			ApiKey:     resp.ApiKey,
			Token:      resp.Token,
			OperatorId: resp.OperatorId,
			UserName:   nil,
		},
		serverIndex: 0,
	}, nil
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
		return c.sandboxClient(profile)
	}
	return c.soracomClient(profile)
}
