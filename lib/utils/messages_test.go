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

// Package util ...
package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageError(t *testing.T) {
	message := Message{
		Code:        "ProvisioningFailed",
		Type:        "Invalid",
		Description: "Failed to create file share with the storage provider",
		RC:          500,
	}
	assert.NotNil(t, message.Error())
	assert.Equal(t, "{Code:ProvisioningFailed, Description:Failed to create file share with the storage provider., RC:500}", message.Error())

	message = Message{
		Code:         "ProvisioningFailed",
		Description:  "",
		BackendError: "Trace Code:03ef81ec-e20b-4ebc-ae96-356631b7e8f1, Code:shares_profile_capacity_iops_invalid, Description:The capacity or IOPS specified in the request is not valid for the 'dp2' profile, RC:400 Bad Request",
		RC:           400,
	}
	assert.NotNil(t, message.Error())
	assert.Equal(t, "{Trace Code:03ef81ec-e20b-4ebc-ae96-356631b7e8f1, Code:shares_profile_capacity_iops_invalid, Description:The capacity or IOPS specified in the request is not valid for the 'dp2' profile, RC:400 Bad Request}", message.Error())
}
