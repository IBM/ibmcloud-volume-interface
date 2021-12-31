/**
 * Copyright 2020 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package auth ...
package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/IBM/ibmcloud-volume-interface/provider/iam"
	"github.com/IBM/ibmcloud-volume-interface/provider/local"

	"github.com/IBM/ibmcloud-volume-interface/lib/provider"
)

const (
	// IMSToken is an IMS user ID and token
	IMSToken = provider.AuthType("IMS_TOKEN")
	// IAMAccessToken ...
	IAMAccessToken = provider.AuthType("IAM_ACCESS_TOKEN")
)

// ForRefreshToken ...
func (ccf *ContextCredentialsFactory) ForRefreshToken(refreshToken string, logger *zap.Logger) (provider.ContextCredentials, error) {
	accessToken, err := ccf.TokenExchangeService.ExchangeRefreshTokenForAccessToken(refreshToken, logger)
	if err != nil {
		// Must preserve provider error code in the ErrorProviderAccountTemporarilyLocked case
		logger.Error("Unable to retrieve access token from refresh token", local.ZapError(err))
		return provider.ContextCredentials{}, err
	}

	imsToken, err := ccf.TokenExchangeService.ExchangeAccessTokenForIMSToken(*accessToken, logger)
	if err != nil {
		// Must preserve provider error code in the ErrorProviderAccountTemporarilyLocked case
		logger.Error("Unable to retrieve IAM token from access token", local.ZapError(err))
		return provider.ContextCredentials{}, err
	}

	return forIMSToken("", imsToken), nil
}

// ForIAMAPIKey ...
func (ccf *ContextCredentialsFactory) ForIAMAPIKey(iamAccountID, apiKey string, logger *zap.Logger) (provider.ContextCredentials, error) {
	imsToken, err := ccf.TokenExchangeService.ExchangeIAMAPIKeyForIMSToken(apiKey, logger)
	if err != nil {
		// Must preserve provider error code in the ErrorProviderAccountTemporarilyLocked case
		logger.Error("Unable to retrieve IMS credentials from IAM API key", local.ZapError(err))
		return provider.ContextCredentials{}, err
	}

	return forIMSToken(iamAccountID, imsToken), nil
}

// ForIAMAccessToken ...
func (ccf *ContextCredentialsFactory) ForIAMAccessToken(apiKey string, logger *zap.Logger) (provider.ContextCredentials, error) {
	//iamAccessToken, err := ccf.TokenExchangeService.ExchangeIAMAPIKeyForAccessToken(apiKey, logger)
	iamAccessToken, err := fetchIAMAccessToken(logger)
	if err != nil {
		logger.Error("Unable to retrieve IAM access token from IAM API key", local.ZapError(err))
		return provider.ContextCredentials{}, err
	}
	iamAccountID, err := ccf.TokenExchangeService.GetIAMAccountIDFromAccessToken(iam.AccessToken{Token: iamAccessToken.Token}, logger)
	if err != nil {
		logger.Error("Unable to retrieve IAM access token from IAM API key", local.ZapError(err))
		return provider.ContextCredentials{}, err
	}

	return forIAMAccessToken(iamAccountID, iamAccessToken), nil
}

// fetchIAMAccessToken ...
func fetchIAMAccessToken(logger *zap.Logger) (*iam.AccessToken, error) {
	logger.Info("Fetching IAM token")
	vaultToken, err := fetchVaultToken()
	if err != nil {
		logger.Error("Unable to read vault token", local.ZapError(err))
		return nil, err
	}
	data := url.Values{}
	endpoint := "https://iam.cloud.ibm.com/identity/token"
	data.Set("cr_token", vaultToken)
	// TODO - receive profile ID as input
	data.Set("profile_id", "")
	data.Set("grant_type", "urn:ibm:params:oauth:grant-type:cr-token")
	client := &http.Client{}
	r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		logger.Error("Unable to create http request", local.ZapError(err))
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil {
		logger.Error("Unable to send request to fetch iam token", local.ZapError(err))
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("Unable to read response received for fetch iam token request", local.ZapError(err))
		return nil, err
	}

	var tokenResponse = struct {
		IAMToken       string `json:"access_token"`
		RefreshToken   string `json:"refresh_token"`
		TokenType      string `json:"token_type"`
		ExpiryDuration int    `json:"expires_in"`
		Expiration     int    `json:"expiration"`
		Scope          string `json:"scope"`
	}{}

	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		logger.Error("Unable to unmarshal response body received for fetch iam token request", local.ZapError(err))
		return nil, err
	}

	if tokenResponse.IAMToken == "" {
		logger.Error("Empty iam token", zap.String("Response", string(body)))
		return nil, errors.New("Failed to fetch iam token")
	}

	logger.Info("Successfully fetched iam token")
	return &iam.AccessToken{Token: tokenResponse.IAMToken}, nil
}

// fetchVaultToken reads vault token mounted into the secrets and returns the same...
func fetchVaultToken() (string, error) {
	// get path from env
	data, err := os.ReadFile("/var/run/secrets/tokens/vault-token")
	if err != nil {
		return "", err
	}
	if string(data) == "" {
		return "", errors.New("unable to fetch vault token")
	}
	return string(data), nil
}

// UpdateAPIKey ...
func (ccf *ContextCredentialsFactory) UpdateAPIKey(apiKey string, logger *zap.Logger) error {
	logger.Info("Updating api key")
	if ccf.TokenExchangeService == nil {
		logger.Error("Failed to update api key in context credentials")
		return errors.New("failed to update api key")
	}
	err := ccf.TokenExchangeService.UpdateAPIKey(apiKey, logger)
	return err
}

// forIMSToken ...
func forIMSToken(iamAccountID string, imsToken *iam.IMSToken) provider.ContextCredentials {
	return provider.ContextCredentials{
		AuthType:     IMSToken,
		IAMAccountID: iamAccountID,
		UserID:       strconv.Itoa(imsToken.UserID),
		Credential:   imsToken.Token,
	}
}

// forIAMAccessToken ...
func forIAMAccessToken(iamAccountID string, iamAccessToken *iam.AccessToken) provider.ContextCredentials {
	return provider.ContextCredentials{
		AuthType:     IAMAccessToken,
		IAMAccountID: iamAccountID,
		Credential:   iamAccessToken.Token,
	}
}
