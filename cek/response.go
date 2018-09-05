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

// OutputSpeechType type
type OutputSpeechType string

// OutputSpeechType constants
const (
	OutputSpeechTypeSimpleSpeech OutputSpeechType = "SimpleSpeech"
	OutputSpeechTypeSpeechList   OutputSpeechType = "SpeechList"
	OutputSpeechTypeSpeechSet    OutputSpeechType = "SpeechSet"
)

// OutputSpeechVerboseType type
type OutputSpeechVerboseType string

// OutputSpeechVerboseType constants
const (
	OutputSpeechVerboseTypeSimpleSpeech OutputSpeechVerboseType = "SimpleSpeech"
	OutputSpeechVerboseTypeSpeechList   OutputSpeechVerboseType = "SpeechList"
)

// SpeechInfoLang type
type SpeechInfoLang string

// SpeechInfoLang constants
const (
	SpeechInfoLangEN    SpeechInfoLang = "en"
	SpeechInfoLangJA    SpeechInfoLang = "ja"
	SpeechInfoLangKO    SpeechInfoLang = "ko"
	SpeechInfoLangEmpty SpeechInfoLang = ""
)

// SpeechInfoType type
type SpeechInfoType string

// SpeechInfoType constants
const (
	SpeechInfoTypePlainText SpeechInfoType = "PlainText"
	SpeechInfoTypeURL       SpeechInfoType = "URL"
)

// ResponseMessage type
type ResponseMessage struct {
	Response          *Response         `json:"response"`
	SessionAttributes map[string]string `json:"sessionAttributes"`
	Version           string            `json:"version"`
}

// Response type
type Response struct {
	Card             interface{}   `json:"card"`
	Directives       []*Directive  `json:"directives"`
	OutputSpeech     *OutputSpeech `json:"outputSpeech"`
	Reprompt         *Reprompt     `json:"reprompt,omitempty"`
	ShouldEndSession bool          `json:"shouldEndSession"`
}

// Directive type
type Directive struct {
	Header  *Header     `json:"header"`
	Payload interface{} `json:"payload"`
}

// Header type
type Header struct {
	MessageID string `json:"messageId"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// OutputSpeech type
type OutputSpeech struct {
	Brief   *SpeechInfo      `json:"brief,omitempty"`
	Type    OutputSpeechType `json:"type"`
	Values  SpeechInfoValues `json:"values,omitempty"`
	Verbose *Verbose         `json:"verbose,omitempty"`
}

// SpeechInfoValues type
type SpeechInfoValues interface {
	SpeechInfoValues()
}

// SpeechInfo type
type SpeechInfo struct {
	Lang  SpeechInfoLang `json:"lang"`
	Type  SpeechInfoType `json:"type"`
	Value string         `json:"value"`
}

// SpeechInfoArray type
type SpeechInfoArray []*SpeechInfo

// SpeechInfoValues method for implementing SpeechInfoValues interface
func (*SpeechInfo) SpeechInfoValues() {}

// SpeechInfoValues method for implementing SpeechInfoValues interface
func (SpeechInfoArray) SpeechInfoValues() {}

// Verbose type
type Verbose struct {
	Type   OutputSpeechVerboseType `json:"type"`
	Values SpeechInfoValues        `json:"values"`
}

// Reprompt type
type Reprompt struct {
	OutputSpeech *OutputSpeech `json:"outputSpeech"`
}
