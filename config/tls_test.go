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

package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGeneralCAHttpClient(t *testing.T) {
	t.Log("Testing GeneralCAHttpClient")

	client, _ := GeneralCAHttpClient()

	assert.NotNil(t, client)
}

func TestGeneralCAHttpClientWithTimeout(t *testing.T) {
	t.Log("Testing GeneralCAHttpClientWithTimeout")

	client, _ := GeneralCAHttpClientWithTimeout(120)

	assert.NotNil(t, client)
	assert.Equal(t, client.Timeout, time.Duration(120))
}
