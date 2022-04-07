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
	"testing"

	"github.com/golang-jwt/jwt/v4"

	"github.com/IBM/ibmcloud-volume-interface/config"
	sp "github.com/IBM/secret-utils-lib/pkg/secret_provider"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_GetIAMAccountIDFromAccessToken(t *testing.T) {
	logger, _ := zap.NewDevelopment(zap.AddCaller())

	fakeAccountID := "12345"
	fakeSigningKey := []byte("aabbccdd")

	fakeToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"account": map[string]interface{}{"bss": fakeAccountID}}).SignedString(fakeSigningKey)

	testcases := []struct {
		name              string
		token             string
		expectedAccountID string
		expectedError     error
	}{
		{
			name:              "fake_token",
			token:             fakeToken,
			expectedAccountID: fakeAccountID,
		},
		{
			name:          "invalid token",
			token:         "invalid",
			expectedError: errors.New("not nil"),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			httpSetup()

			authConfig := &AuthConfiguration{
				IamURL:          server.URL,
				IamClientID:     "test",
				IamClientSecret: "secret",
			}

			tes := new(tokenExchangeService)
			tes.httpClient, _ = config.GeneralCAHttpClient()
			tes.secretprovider = new(sp.FakeSecretProvider)
			tes.authConfig = authConfig
			accountID, err := tes.GetIAMAccountIDFromAccessToken(AccessToken{Token: testcase.token}, logger)
			if testcase.expectedError != nil {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, testcase.expectedAccountID, accountID)
			}
		})
	}
}
