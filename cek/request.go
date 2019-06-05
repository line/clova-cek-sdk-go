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
	"errors"
)

// PlayerActivity type
type PlayerActivity string

// PlayerActivity constants
const (
	PlayerActivityIDLE    PlayerActivity = "IDLE"
	PlayerActivityPLAYING PlayerActivity = "PLAYING"
	PlayerActivityPAUSED  PlayerActivity = "PAUSED"
	PlayerActivitySTOPPED PlayerActivity = "STOPPED"
)

// Orientation type
type Orientation string

// Orientation constants
const (
	OrientationLandscape Orientation = "landscape"
	OrientationPortrait  Orientation = "portrait"
)

// Size type
type Size string

// Size constants
const (
	SizeNone   Size = "none"
	SizeS100   Size = "s100"
	SizeM100   Size = "m100"
	SizeL100   Size = "l100"
	SizeXL100  Size = "xl100"
	SizeCustom Size = "custom"
)

// RequestType type
type RequestType string

// RequestType constants
const (
	RequestTypeEvent        RequestType = "EventRequest"
	RequestTypeIntent       RequestType = "IntentRequest"
	RequestTypeLaunch       RequestType = "LaunchRequest"
	RequestTypeSessionEnded RequestType = "SessionEndedRequest"
)

// SlotValueType type
type SlotValueType string

// SlotValueType constants
const (
	SlotValueTypeTime             SlotValueType = "TIME"
	SlotValueTypeDate             SlotValueType = "DATE"
	SlotValueTypeDateTime         SlotValueType = "DATETIME"
	SlotValueTypeTimeInterval     SlotValueType = "TIME.INTERVAL"
	SlotValueTypeDateInterval     SlotValueType = "DATE.INTERVAL"
	SlotValueTypeDateTimeInterval SlotValueType = "DATETIME.INTERVAL"
)

// RequestMessage type
type RequestMessage struct {
	Context *Context `json:"context"`
	Request Request  `json:"request"`
	Session *Session `json:"session"`
	Version string   `json:"version"`
}

type rawRequest struct {
	Type    RequestType `json:"type"`
	Request Request     `json:"-"`
}

func (r *rawRequest) UnmarshalJSON(b []byte) error {
	type alias rawRequest
	raw := alias{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	var request Request
	switch raw.Type {
	case RequestTypeEvent:
		request = &EventRequest{}
	case RequestTypeIntent:
		request = &IntentRequest{}
	case RequestTypeLaunch:
		request = &LaunchRequest{}
	case RequestTypeSessionEnded:
		request = &SessionEndedRequest{}
	default:
		return errors.New("invalid request type")
	}
	if err := json.Unmarshal(b, request); err != nil {
		return err
	}
	r.Request = request
	return nil
}

// UnmarshalJSON method for RequestMessage
func (m *RequestMessage) UnmarshalJSON(b []byte) error {
	type alias RequestMessage
	raw := struct {
		Request rawRequest `json:"request"`
		*alias
	}{
		alias: (*alias)(m),
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	m.Request = raw.Request.Request
	return nil
}

// Context type
type Context struct {
	AudioPlayer *AudioPlayer `json:"AudioPlayer,omitempty"`
	System      *System      `json:"System"`
}

// AudioPlayer type
type AudioPlayer struct {
	OffsetInMilliseconds int            `json:"offsetInMilliseconds,omitempty"`
	PlayerActivity       PlayerActivity `json:"playerActivity"`
	Stream               interface{}    `json:"stream,omitempty"`
	TotalInMilliseconds  int            `json:"totalInMilliseconds,omitempty"`
}

// System type
type System struct {
	Application *Application `json:"application"`
	Device      *Device      `json:"device"`
	User        *User        `json:"user"`
}

// Application type
type Application struct {
	ApplicationID string `json:"applicationId"`
}

// Device type
type Device struct {
	DeviceID string   `json:"deviceId"`
	Display  *Display `json:"display"`
}

// Display type
type Display struct {
	ContentLayer *ContentLayer `json:"contentLayer,omitempty"`
	DPI          int           `json:"dpi,omitempty"`
	Orientation  Orientation   `json:"orientation,omitempty"`
	Size         Size          `json:"size"`
}

// ContentLayer type
type ContentLayer struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// User type
type User struct {
	UserID      string `json:"userId"`
	AccessToken string `json:"accessToken,omitempty"`
}

// Request interface
type Request interface {
	Request()
}

// EventRequest type
type EventRequest struct {
	Event     *Event `json:"event"`
	RequestID string `json:"requestId"`
	Timestamp string `json:"timestamp"`
}

// Event type
type Event struct {
	Name      string      `json:"name"`
	Namespace string      `json:"namespace"`
	Payload   interface{} `json:"payload"`
}

// IntentRequest type
type IntentRequest struct {
	Intent *Intent `json:"intent"`
}

// Intent type
type Intent struct {
	Name  string           `json:"name"`
	Slots map[string]*Slot `json:"slots"`
}

// Slot type
type Slot struct {
	Name      string        `json:"name"`
	Value     string        `json:"value"`
	ValueType SlotValueType `json:"valueType,omitempty"`
	Unit      string        `json:"unit,omitempty"`
}

// LaunchRequest type
type LaunchRequest struct {
}

// SessionEndedRequest type
type SessionEndedRequest struct {
}

// Request method for implementing Request interface
func (*EventRequest) Request() {}

// Request method for implementing Request interface
func (*IntentRequest) Request() {}

// Request method for implementing Request interface
func (*LaunchRequest) Request() {}

// Request method for implementing Request interface
func (*SessionEndedRequest) Request() {}

// Session type
type Session struct {
	New               bool              `json:"new"`
	SessionAttributes map[string]string `json:"sessionAttributes"`
	SessionID         string            `json:"sessionId"`
	User              *User             `json:"user"`
}
