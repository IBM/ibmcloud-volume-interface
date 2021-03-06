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
	"testing"

	"github.com/IBM/ibmcloud-volume-interface/provider/iam"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func init() {
	logger, _ = zap.NewDevelopment()
}

func TestNewContextCredentialsFactory(t *testing.T) {
	authConfig := &iam.AuthConfiguration{
		IamClientID:     "test-client-id",
		IamClientSecret: "test-client-secret",
	}

	contextCredentials, err := NewContextCredentialsFactory(authConfig)
	assert.NoError(t, err)
	assert.NotNil(t, contextCredentials)
}
