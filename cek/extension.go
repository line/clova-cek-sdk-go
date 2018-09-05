// Copyright 2018 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package cek

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Extension type
type Extension struct {
	ID        string
	debugMode bool
}

// ExtensionOption type
type ExtensionOption func(*Extension)

// NewExtension function
func NewExtension(extensionID string, options ...ExtensionOption) *Extension {
	ext := &Extension{
		ID: extensionID,
	}
	for _, option := range options {
		option(ext)
	}
	return ext
}

// WithDebugMode function
func WithDebugMode(ext *Extension) {
	ext.debugMode = true
}

// ParseRequest method
func (e *Extension) ParseRequest(r *http.Request) (*RequestMessage, error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if !e.debugMode {
		if err := validateSignature(r.Header.Get("SignatureCEK"), body); err != nil {
			return nil, fmt.Errorf("invalid signature: %s", err.Error())
		}
	}

	message := &RequestMessage{}
	if err := json.Unmarshal(body, message); err != nil {
		return nil, err
	}
	if message.Context != nil && message.Context.System != nil && message.Context.System.Application != nil &&
		message.Context.System.Application.ApplicationID == e.ID {
		return message, nil
	}
	return nil, fmt.Errorf("invalid application")
}
