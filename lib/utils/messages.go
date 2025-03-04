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
	"fmt"
	"strings"
)

// Message Wrapper Message/Error Class
type Message struct {
	Code          string
	Type          string
	RequestID     string
	Description   string
	BackendError  string
	InternalError string
	RC            int
	Action        string
}

// Error Implement the Error() interface method
func (msg Message) Error() string {
	return msg.Info()
}

// Info ...
func (msg Message) Info() string {
	if strings.Contains(msg.BackendError, "Trace Code:") {
		return fmt.Sprintf("{%s}", msg.BackendError)
	} else {
		return fmt.Sprintf("%s.%s", msg.Description, msg.BackendError)
	}
}
