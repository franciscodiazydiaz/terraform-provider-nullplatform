package nullplatform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

const TOKEN_PATH = "/token"

type TokenRequest struct {
	Apikey string `json:"apikey"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type NullClient struct {
	Client *retryablehttp.Client
	ApiURL string
	ApiKey string
	Token  Token
}

type NullErrors struct {
	Message string `json:"message"`
	Id      int    `json:"id"`
}

type NullOps interface {
	GetToken() diag.Diagnostics

	CreateScope(*Scope) (*Scope, error)
	PatchScope(string, *Scope) error
	GetScope(string) (*Scope, error)

	PatchNRN(string, *PatchNRN) error
	GetNRN(string) (*NRN, error)

	GetApplication(appId string) (*Application, error)

	CreateService(*Service) (*Service, error)
	GetService(string) (*Service, error)
	PatchService(string, *Service) error
	DeleteService(string) error

	CreateLink(*Link) (*Link, error)
	PatchLink(string, *Link) error
	DeleteLink(string) error
	GetLink(string) (*Link, error)

	CreateParameter(param *Parameter, importIfCreated bool) (*Parameter, error)
	PatchParameter(parameterId string, param *Parameter) error
	GetParameter(parameterId string) (*Parameter, error)
	DeleteParameter(parameterId string) error
	GetParameterList(nrn string) (*ParameterList, error)

	CreateParameterValue(paramId int, paramValue *ParameterValue) (*ParameterValue, error)
	GetParameterValue(parameterId string, parameterValueId string) (*ParameterValue, error)
	DeleteParameterValue(parameterId string, parameterValueId string) error
}

func NewNullClient(apiURL, apiKey string, maxRetries int) *NullClient {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = maxRetries
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 30 * time.Second
	retryClient.CheckRetry = customRetryPolicy
	retryClient.Logger = nil // Disable logging

	return &NullClient{
		Client: retryClient,
		ApiURL: apiURL,
		ApiKey: apiKey,
	}
}

func customRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// Network error or other errors that resulted in no response
	if err != nil {
		return true, err
	}

	if resp.StatusCode == http.StatusConflict || resp.StatusCode == http.StatusTooManyRequests {
		return true, nil
	}

	// Don't retry for other cases
	return false, nil
}

func (c *NullClient) MakeRequest(method, path string, body any) (*http.Response, error) {
	url := fmt.Sprintf("https://%s%s", c.ApiURL, path)
	req, err := retryablehttp.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token.AccessToken))

	return c.Client.Do(req)
}

func (c *NullClient) GetToken() diag.Diagnostics {
	treq := TokenRequest{
		Apikey: c.ApiKey,
	}

	jsonBody, err := json.Marshal(treq)
	if err != nil {
		return diag.FromErr(err)
	}

	url := fmt.Sprintf("https://%s%s", c.ApiURL, TOKEN_PATH)

	// Use the client's Post method directly
	res, err := c.Client.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return diag.FromErr(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return diag.FromErr(fmt.Errorf("error creating resource, got %d, api key was %s", res.StatusCode, c.ApiKey))
	}

	tRes := &Token{}
	err = json.NewDecoder(res.Body).Decode(tRes)
	if err != nil {
		return diag.FromErr(err)
	}

	if tRes.AccessToken == "" {
		return diag.FromErr(fmt.Errorf("no access token for null platform token rsp is: %s", tRes))
	}

	c.Token = (*tRes)

	return nil
}
