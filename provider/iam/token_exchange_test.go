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
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	util "github.com/IBM/ibmcloud-volume-interface/lib/utils"
	"github.com/IBM/ibmcloud-volume-interface/lib/utils/reasoncode"

	"github.com/IBM/ibmcloud-volume-interface/config"
	"github.com/IBM/secret-utils-lib/pkg/k8s_utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	mux              *http.ServeMux
	server           *httptest.Server
	logger           *zap.Logger
	lowPriority      zap.LevelEnablerFunc
	consoleDebugging zapcore.WriteSyncer
)

func TestMain(m *testing.M) {
	// Logging
	lowPriority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	consoleDebugging = zapcore.Lock(os.Stdout)
	logger = zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority), zap.AddCaller())

	os.Exit(m.Run())
}

func httpSetup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
}

func Test_ExchangeRefreshTokenForAccessToken_Success(t *testing.T) {
	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
		zap.AddCaller(),
	)
	httpSetup()

	// IAM endpoint
	mux.HandleFunc("/oidc/token",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{"access_token": "at_success","refresh_token": "rt_success", "expiration": 456, "uaa_token": "uaa_success"}`)
		},
	)

	authConfig := &AuthConfiguration{
		IamURL:          server.URL,
		IamClientID:     "test",
		IamClientSecret: "secret",
	}

	tes := new(tokenExchangeService)
	tes.httpClient, _ = config.GeneralCAHttpClient()
	tes.authConfig = authConfig
	tes.secretprovider = new(FakeSecretProvider)
	r, err := tes.ExchangeRefreshTokenForAccessToken("testrefreshtoken", logger)
	assert.Nil(t, err)
	if assert.NotNil(t, r) {
		assert.Equal(t, (*r).Token, "at_success")
	}
}

func Test_ExchangeRefreshTokenForAccessToken_FailedDuringRequest(t *testing.T) {
	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
		zap.AddCaller(),
	)

	httpSetup()

	mux.HandleFunc("/oidc/token",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"errorMessage": "did not work",
				"errorCode": "bad news",
				"errorDetails" : "more details",
				"requirements" : { "error": "requirements error", "code":"requirements code" }
				}`)
		},
	)

	authConfig := &AuthConfiguration{
		IamURL:          server.URL,
		IamClientID:     "test",
		IamClientSecret: "secret",
	}

	tes := new(tokenExchangeService)
	tes.httpClient, _ = config.GeneralCAHttpClient()
	tes.authConfig = authConfig
	tes.secretprovider = new(FakeSecretProvider)
	r, err := tes.ExchangeRefreshTokenForAccessToken("badrefreshtoken", logger)
	assert.Nil(t, r)
	if assert.NotNil(t, err) {
		assert.Equal(t, "IAM token exchange request failed: did not work", err.Error())
		assert.Equal(t, reasoncode.ReasonCode("ErrorFailedTokenExchange"), util.ErrorReasonCode(err))
	}
}

func Test_ExchangeRefreshTokenForAccessToken_FailedDuringRequest_no_message(t *testing.T) {
	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
		zap.AddCaller(),
	)

	httpSetup()

	mux.HandleFunc("/oidc/token",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{}`)
		},
	)

	authConfig := &AuthConfiguration{
		IamURL:          server.URL,
		IamClientID:     "test",
		IamClientSecret: "secret",
	}

	tes := new(tokenExchangeService)
	tes.httpClient, _ = config.GeneralCAHttpClient()
	tes.authConfig = authConfig
	tes.secretprovider = new(FakeSecretProvider)
	r, err := tes.ExchangeRefreshTokenForAccessToken("badrefreshtoken", logger)
	assert.Nil(t, r)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Unexpected IAM token exchange response", err.Error())
		assert.Equal(t, reasoncode.ReasonCode("ErrorUnclassified"), util.ErrorReasonCode(err))
	}
}

func Test_ExchangeRefreshTokenForAccessToken_FailedNoIamUrl(t *testing.T) {
	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
		zap.AddCaller(),
	)

	httpSetup()

	mux.HandleFunc("/oidc/token",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{}`)
		},
	)

	authConfig := &AuthConfiguration{
		IamURL:          "",
		IamClientID:     "test",
		IamClientSecret: "secret",
	}

	tes := new(tokenExchangeService)
	tes.httpClient, _ = config.GeneralCAHttpClient()
	tes.authConfig = authConfig
	tes.secretprovider = new(FakeSecretProvider)
	r, err := tes.ExchangeRefreshTokenForAccessToken("testrefreshtoken", logger)
	assert.Nil(t, r)

	if assert.NotNil(t, err) {
		assert.Equal(t, "IAM token exchange request failed", err.Error())
		assert.Equal(t, reasoncode.ReasonCode("ErrorUnclassified"), util.ErrorReasonCode(err))
		assert.Equal(t, []string{"Post \"/oidc/token\": unsupported protocol scheme \"\""},
			util.ErrorDeepUnwrapString(err))
	}
}

func Test_ExchangeRefreshTokenForAccessToken_FailedRequesting_empty_body(t *testing.T) {
	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
		zap.AddCaller(),
	)

	httpSetup()

	mux.HandleFunc("/oidc/token",
		func(w http.ResponseWriter, r *http.Request) {
			// Leave response empty
		},
	)

	authConfig := &AuthConfiguration{
		IamURL:          server.URL,
		IamClientID:     "test",
		IamClientSecret: "secret",
	}

	tes := new(tokenExchangeService)
	tes.httpClient, _ = config.GeneralCAHttpClient()
	tes.authConfig = authConfig
	tes.secretprovider = new(FakeSecretProvider)
	r, err := tes.ExchangeRefreshTokenForAccessToken("badrefreshtoken", logger)
	assert.Nil(t, r)

	if assert.NotNil(t, err) {
		assert.Equal(t, "IAM token exchange request failed", err.Error())
		assert.Equal(t, reasoncode.ReasonCode("ErrorUnclassified"), util.ErrorReasonCode(err))
		assert.Equal(t, []string{"empty response body"},
			util.ErrorDeepUnwrapString(err))
	}
}

func TestNewTokenExchangeServiceWithClient(t *testing.T) {
	authConfig := &AuthConfiguration{
		IamURL:          server.URL,
		IamClientID:     "test",
		IamClientSecret: "secret",
	}

	newClient := &http.Client{
		Timeout: time.Second * 10,
	}
	tes, err := NewTokenExchangeServiceWithClient(authConfig, newClient)
	assert.NoError(t, err)
	assert.NotNil(t, tes)
}

func TestExchangeIAMAPIKeyForIMSToken(t *testing.T) {
	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
		zap.AddCaller(),
	)

	httpSetup()

	mux.HandleFunc("/oidc/token",
		func(w http.ResponseWriter, r *http.Request) {
			// Leave response empty
		},
	)

	authConfig := &AuthConfiguration{
		IamURL:          server.URL,
		IamClientID:     "test",
		IamClientSecret: "secret",
	}

	tes := new(tokenExchangeService)
	tes.httpClient, _ = config.GeneralCAHttpClient()
	tes.authConfig = authConfig
	tes.secretprovider = new(FakeSecretProvider)
	ims, err := tes.ExchangeIAMAPIKeyForIMSToken("badrefreshtoken", logger)
	assert.Nil(t, ims)
	assert.Error(t, err)
}

func Test_ExchangeAccessTokenForIMSToken_Success(t *testing.T) {
	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
		zap.AddCaller(),
	)
	httpSetup()

	// IAM endpoint
	mux.HandleFunc("/oidc/token",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprint(w, `{"access_token": "at_success","refresh_token": "rt_success", "ims_user_id": 123, "ims_token": "ims_token_1"}`)
		},
	)

	authConfig := &AuthConfiguration{
		IamURL:          server.URL,
		IamClientID:     "test",
		IamClientSecret: "secret",
	}

	tes := new(tokenExchangeService)
	tes.httpClient, _ = config.GeneralCAHttpClient()
	tes.authConfig = authConfig
	tes.secretprovider = new(FakeSecretProvider)
	r, err := tes.ExchangeAccessTokenForIMSToken(AccessToken{Token: "testaccesstoken"}, logger)
	assert.Nil(t, err)
	if assert.NotNil(t, r) {
		assert.Equal(t, (*r).UserID, 123)
		assert.Equal(t, (*r).Token, "ims_token_1")
	}
}

func Test_ExchangeAccessTokenForIMSToken_FailedDuringRequest(t *testing.T) {
	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
		zap.AddCaller(),
	)

	httpSetup()

	mux.HandleFunc("/oidc/token",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"errorMessage": "did not work",
				"errorCode": "bad news",
				"errorDetails" : "more details",
				"requirements" : { "error": "requirements error", "code":"requirements code" }
				}`)
		},
	)

	authConfig := &AuthConfiguration{
		IamURL:          server.URL,
		IamClientID:     "test",
		IamClientSecret: "secret",
	}

	tes := new(tokenExchangeService)
	tes.httpClient, _ = config.GeneralCAHttpClient()
	tes.authConfig = authConfig
	tes.secretprovider = new(FakeSecretProvider)
	r, err := tes.ExchangeAccessTokenForIMSToken(AccessToken{Token: "badaccesstoken"}, logger)
	assert.Nil(t, r)
	if assert.NotNil(t, err) {
		assert.Equal(t, "IAM token exchange request failed: did not work", err.Error())
		assert.Equal(t, reasoncode.ReasonCode("ErrorFailedTokenExchange"), util.ErrorReasonCode(err))
		assert.Equal(t, []string{"more details requirements code: requirements error"},
			util.ErrorDeepUnwrapString(err))
	}
}

func Test_ExchangeAccessTokenForIMSToken_FailedAccountLocked(t *testing.T) {
	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
		zap.AddCaller(),
	)

	httpSetup()

	mux.HandleFunc("/oidc/token",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `
				{"errorMessage": "OpenID Connect exception",
				"errorCode": "BXNIM0400E",
				"errorDetails" : "Failed external authentication.",
				"requirements" : { "error": "Account has been locked for 30 minutes", "code":"SoftLayer_Exception_User_Customer_AccountLocked" }
				}`)
		},
	)

	authConfig := &AuthConfiguration{
		IamURL:          server.URL,
		IamClientID:     "test",
		IamClientSecret: "secret",
	}

	tes := new(tokenExchangeService)
	tes.httpClient, _ = config.GeneralCAHttpClient()
	tes.authConfig = authConfig
	tes.secretprovider = new(FakeSecretProvider)
	r, err := tes.ExchangeAccessTokenForIMSToken(AccessToken{Token: "badaccesstoken"}, logger)
	assert.Nil(t, r)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Infrastructure account is temporarily locked", err.Error())
		assert.Equal(t, reasoncode.ReasonCode("ErrorProviderAccountTemporarilyLocked"), util.ErrorReasonCode(err))
		assert.Equal(t, []string{"IAM token exchange request failed: OpenID Connect exception", "Failed external authentication. SoftLayer_Exception_User_Customer_AccountLocked: Account has been locked for 30 minutes"}, util.ErrorDeepUnwrapString(err))
	}
}

func Test_ExchangeAccessTokenForIMSToken_FailedDuringRequest_no_message(t *testing.T) {
	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
		zap.AddCaller(),
	)

	httpSetup()

	mux.HandleFunc("/oidc/token",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{}`)
		},
	)

	authConfig := &AuthConfiguration{
		IamURL:          server.URL,
		IamClientID:     "test",
		IamClientSecret: "secret",
	}

	tes := new(tokenExchangeService)
	tes.httpClient, _ = config.GeneralCAHttpClient()
	tes.authConfig = authConfig
	tes.secretprovider = new(FakeSecretProvider)
	r, err := tes.ExchangeAccessTokenForIMSToken(AccessToken{Token: "badrefreshtoken"}, logger)
	assert.Nil(t, r)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Unexpected IAM token exchange response", err.Error())
		assert.Equal(t, reasoncode.ReasonCode("ErrorUnclassified"), util.ErrorReasonCode(err))
	}
}

func Test_ExchangeAccessTokenForIMSToken_FailedRequesting_empty_body(t *testing.T) {
	logger := zap.New(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
		zap.AddCaller(),
	)

	httpSetup()

	mux.HandleFunc("/oidc/token",
		func(w http.ResponseWriter, r *http.Request) {
			// Leave response empty
		},
	)

	authConfig := &AuthConfiguration{
		IamURL:          server.URL,
		IamClientID:     "test",
		IamClientSecret: "secret",
	}

	tes := new(tokenExchangeService)
	tes.httpClient, _ = config.GeneralCAHttpClient()
	tes.authConfig = authConfig
	tes.secretprovider = new(FakeSecretProvider)
	r, err := tes.ExchangeAccessTokenForIMSToken(AccessToken{Token: "badrefreshtoken"}, logger)
	assert.Nil(t, r)

	if assert.NotNil(t, err) {
		assert.Equal(t, "IAM token exchange request failed", err.Error())
		assert.Equal(t, reasoncode.ReasonCode("ErrorUnclassified"), util.ErrorReasonCode(err))
		assert.Equal(t, []string{"empty response body"}, util.ErrorDeepUnwrapString(err))
	}
}

func Test_ExchangeIAMAPIKeyForAccessToken(t *testing.T) {
	var testCases = []struct {
		name               string
		expectedError      error
		expectedReasonCode string
		secretProvider     *FakeSecretProvider
	}{
		{
			name:          "Successfully fetched token",
			expectedError: nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			logger := zap.New(
				zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), consoleDebugging, lowPriority),
				zap.AddCaller(),
			)
			httpSetup()

			authConfig := &AuthConfiguration{
				IamURL: server.URL,
			}
			tes := new(tokenExchangeService)
			tes.httpClient, _ = config.GeneralCAHttpClient()
			tes.authConfig = authConfig
			tes.secretprovider = new(FakeSecretProvider)
			_, actualError := tes.ExchangeIAMAPIKeyForAccessToken("apikey1", logger)
			if testCase.expectedError == nil {
				assert.Nil(t, actualError)
			} else {
				assert.NotNil(t, actualError)
			}
		})
	}
}

func Test_NewTokenExchangeService(t *testing.T) {
	authConfig := &AuthConfiguration{
		IamURL: server.URL,
	}

	// Test without any k8s client
	_, err := NewTokenExchangeService(authConfig, nil)
	assert.NotNil(t, err)

	// Test with any k8s client
	k8sClient, _ := k8s_utils.FakeGetk8sClientSet()
	pwd, _ := os.Getwd()
	file := filepath.Join(pwd, "..", "..", "etc", "libconfig.toml")
	err = k8s_utils.FakeCreateSecret(k8sClient, "DEFAULT", file)
	_, err = NewTokenExchangeService(authConfig, &k8sClient)
	assert.Nil(t, err)
}

const (
	// fakeEndpoint ...
	fakeEndpoint = "https://fakehost.com"
)

// FakeSecretProvider ...
type FakeSecretProvider struct {
}

// GetDefaultIAMToken ...
func (fs *FakeSecretProvider) GetDefaultIAMToken(freshTokenRequired bool, reasonForCall ...string) (string, uint64, error) {
	if !freshTokenRequired {
		return "token", 1000, nil
	}
	return "", 0, errors.New("fake error")
}

// GetIAMToken ...
func (fs *FakeSecretProvider) GetIAMToken(secret string, freshTokenRequired bool, reasonForCall ...string) (string, uint64, error) {
	if freshTokenRequired {
		return "token", 1000, nil
	}
	return "", 0, errors.New("fake error")
}

// GetRIAASEndpoint ...
func (fs *FakeSecretProvider) GetRIAASEndpoint(readConfig bool) (string, error) {
	return fakeEndpoint, nil
}

// GetPrivateRIAASEndpoint ...
func (fs *FakeSecretProvider) GetPrivateRIAASEndpoint(readConfig bool) (string, error) {
	return fakeEndpoint, nil
}

// GetContainerAPIRoute ...
func (fs *FakeSecretProvider) GetContainerAPIRoute(readConfig bool) (string, error) {
	return fakeEndpoint, nil
}

// GetPrivateContainerAPIRoute ...
func (fs *FakeSecretProvider) GetPrivateContainerAPIRoute(readConfig bool) (string, error) {
	return fakeEndpoint, nil
}

// GetResourceGroupID ...
func (fs *FakeSecretProvider) GetResourceGroupID() string {
	return "resource-group-id"
}
