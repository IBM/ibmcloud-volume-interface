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

package util

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetErrorType(t *testing.T) {
	err := errors.New("Infrastructure account is temporarily locked")
	newErr := NewError("ErrorProviderAccountTemporarilyLocked", "Infrastructure account is temporarily locked", err)
	assert.NotNil(t, GetErrorType(newErr))
	newErr = NewError("ProvisioningFailed", "ProvisioningFailed", errors.New("ProvisioningFailed"))
	assert.NotNil(t, GetErrorType(newErr))
	newErr = NewErrorWithProperties("ProvisioningFailed", "", map[string]string{"properties": "properties"}, errors.New("ProvisioningFailed"))
	assert.NotNil(t, GetErrorType(newErr))
}
