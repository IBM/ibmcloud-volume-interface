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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/IBM/ibmcloud-volume-interface/provider/iam"
	"github.com/IBM/secret-utils-lib/pkg/k8s_utils"
	"github.com/stretchr/testify/assert"
)

func TestNewContextCredentialsFactory(t *testing.T) {
	// Pass without k8s client
	authConfig := &iam.AuthConfiguration{
		IamURL:          "url",
		IamClientID:     "test",
		IamClientSecret: "secret",
	}

	_, err := NewContextCredentialsFactory(authConfig, nil)
	assert.NotNil(t, err)

	// Pass with k8s client
	k8sClient, _ := k8s_utils.FakeGetk8sClientSet()
	pwd, _ := os.Getwd()
	file := filepath.Join(pwd, "..", "..", "etc", "libconfig.toml")
	_ = k8s_utils.FakeCreateSecret(k8sClient, "DEFAULT", file)
	_, err = NewContextCredentialsFactory(authConfig, &k8sClient)
	fmt.Println(err)
	assert.Nil(t, err)
}
