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

// Package iam ...
package iam

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	sp "github.com/IBM/ibmcloud-volume-interface/provider/secretprovider"
	"github.com/IBM/ibmcloud-volume-interface/provider/token"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	endpoint = flag.String("sidecarendpoint", "/csi/provider.sock", "Storage secret sidecar endpoint")
)

const (
	// tokenExpirydiff ...
	tokenExpirydiff = 600
)

// ComputeIdentityProvider ...
type ComputeIdentityProvider struct {
	defaultProfileID string
	iksEnabled       bool
	vaultToken       string
	logger           *zap.Logger
}

// NewComputeIdentityProvider ...
func NewComputeIdentityProvider(profileID string, iksEnabled bool, logger *zap.Logger) (TokenProvider, error) {
	vaultToken, err := ReadVaultToken()
	if err != nil {
		logger.Error("Error initializing compute identity provider", zap.Error(err))
		return nil, err
	}
	computeIdentityProvider := &ComputeIdentityProvider{
		defaultProfileID: profileID,
		iksEnabled:       iksEnabled,
		vaultToken:       vaultToken,
		logger:           logger,
	}
	return computeIdentityProvider, nil
}

// GetIAMToken ...
func (cp *ComputeIdentityProvider) GetIAMToken(profileID string, freshTokenRequired bool) (string, uint64, error) {
	cp.logger.Info("In GetIAMToken(), fetching iam token via compute identity method")

	var tokenExpiryTime uint64
	// If IKS is enabled, call goes to sidecar
	if cp.iksEnabled {
		conn, err := grpc.Dial(*endpoint, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDialer(unixConnect)) //nolint:staticcheck
		if err != nil {
			cp.logger.Error("Unable to setup grpc session", zap.Error(err))
			return "", tokenExpiryTime, err
		}
		c := sp.NewIAMTokenProviderClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		defer conn.Close()
		resp, err := c.GetIAMToken(ctx, &sp.Request{ProfileId: profileID, IsFreshTokenRequired: freshTokenRequired})
		if err != nil {
			cp.logger.Error("Error fetching iam token from grpc call", zap.Error(err))
			return "", tokenExpiryTime, err
		}
		return resp.Iamtoken, resp.Tokenlifetime, nil
	}

	// checking if the vault token is valid, if invalid reading it again
	if !token.IsTokenValid(cp.logger, cp.vaultToken) {
		vaultToken, err := ReadVaultToken()
		if err != nil {
			cp.logger.Error("Error reading vault token", zap.Error(err))
			return "", tokenExpiryTime, err
		}
		cp.vaultToken = vaultToken
	}

	iamToken, err := SendGetTokenRequest(cp.logger, profileID, cp.vaultToken)
	if err != nil {
		cp.logger.Error("Error fetching iam token", zap.Error(err))
		return "", tokenExpiryTime, err
	}
	tokenExpiryTime, err = token.FetchTokenLifeTime(cp.logger, iamToken, tokenExpirydiff)
	if err != nil {
		cp.logger.Error("Error fetch token lifetime", zap.Error(err))
		return "", tokenExpiryTime, err
	}
	cp.logger.Info("Successfully fetched iam token")
	return iamToken, tokenExpiryTime, nil
}

// GetDefaultIAMToken ...
func (cp *ComputeIdentityProvider) GetDefaultIAMToken(freshTokenRequired bool) (string, uint64, error) {
	return cp.GetIAMToken(cp.defaultProfileID, freshTokenRequired)
}

// SendGetTokenRequest ...
func SendGetTokenRequest(logger *zap.Logger, profileID, vaultToken string) (string, error) {
	logger.Info("Sending get token request to iam")
	data := url.Values{}
	data.Set("cr_token", vaultToken)
	data.Set("profile_id", profileID)
	data.Set("grant_type", "urn:ibm:params:oauth:grant-type:cr-token")

	client := &http.Client{}
	r, err := http.NewRequest("POST", "https://iam.cloud.ibm.com/identity/token", strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		logger.Error("Error creating http request", zap.Error(err))
		return "", err
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil {
		logger.Error("Error sending token exchange request", zap.Error(err))
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("Error reading token exchange response", zap.Error(err))
		return "", err
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
		logger.Error("Error unmarshaling token response", zap.Error(err))
		return "", err
	}

	if tokenResponse.IAMToken == "" {
		logger.Error("Empty token", zap.String("Response received", string(body)))
		return "", errors.New("failed to fetch iam token")
	}

	logger.Info("Successfully fetched iam token")
	return tokenResponse.IAMToken, nil
}

// ReadVaultToken ...
func ReadVaultToken() (string, error) {
	tokenPath := os.Getenv("VAULT_TOKEN_PATH")
	if tokenPath == "" {
		return "", errors.New("VAULT_TOKEN_PATH not set")
	}

	byteData, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", err
	}
	return string(byteData), nil
}

func unixConnect(addr string, t time.Duration) (net.Conn, error) {
	unixAddr, err := net.ResolveUnixAddr("unix", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUnix("unix", nil, unixAddr)
	return conn, err
}
