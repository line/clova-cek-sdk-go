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

// ResponseBuilder type
type ResponseBuilder struct {
	sessionAttributes map[string]string
	card              map[string]interface{}
	directives        []*Directive
	outputSpeech      *OutputSpeech
	reprompt          *OutputSpeech
	shouldEndSession  bool
}

// NewResponseBuilder function
func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{
		sessionAttributes: map[string]string{},
		card:              map[string]interface{}{},
		directives:        []*Directive{},
	}
}

// SessionAttributes method
func (b *ResponseBuilder) SessionAttributes(attributes map[string]string) *ResponseBuilder {
	b.sessionAttributes = attributes
	return b
}

// OutputSpeech method
func (b *ResponseBuilder) OutputSpeech(os *OutputSpeech) *ResponseBuilder {
	b.outputSpeech = os
	return b
}

// Reprompt method
func (b *ResponseBuilder) Reprompt(os *OutputSpeech) *ResponseBuilder {
	b.reprompt = os
	return b
}

// ShouldEndSession method
func (b *ResponseBuilder) ShouldEndSession(v bool) *ResponseBuilder {
	b.shouldEndSession = v
	return b
}

// AddDirective method
func (b *ResponseBuilder) AddDirective(d *Directive) *ResponseBuilder {
	b.directives = append(b.directives, d)
	return b
}

// Build method
func (b *ResponseBuilder) Build() *ResponseMessage {
	var reprompt *Reprompt
	if b.reprompt != nil {
		reprompt = &Reprompt{
			OutputSpeech: b.reprompt,
		}
	}
	return &ResponseMessage{
		Response: &Response{
			Card:             b.card,
			Directives:       b.directives,
			OutputSpeech:     b.outputSpeech,
			ShouldEndSession: b.shouldEndSession,
			Reprompt:         reprompt,
		},
		SessionAttributes: b.sessionAttributes,
		Version:           "1.0",
	}
}

// OutputSpeechBuilder type
type OutputSpeechBuilder struct {
	brief    *SpeechInfo
	verbose  *Verbose
	speeches []*SpeechInfo
}

// NewOutputSpeechBuilder function
func NewOutputSpeechBuilder() *OutputSpeechBuilder {
	return &OutputSpeechBuilder{}
}

// AddSpeechText method
func (b *OutputSpeechBuilder) AddSpeechText(text string, lang SpeechInfoLang) *OutputSpeechBuilder {
	b.speeches = append(b.speeches, &SpeechInfo{
		Lang:  lang,
		Type:  SpeechInfoTypePlainText,
		Value: text,
	})
	return b
}

// AddSpeechURL method
func (b *OutputSpeechBuilder) AddSpeechURL(url string) *OutputSpeechBuilder {
	b.speeches = append(b.speeches, &SpeechInfo{
		Lang:  SpeechInfoLangEmpty,
		Type:  SpeechInfoTypeURL,
		Value: url,
	})
	return b
}

// SpeechSet method
func (b *OutputSpeechBuilder) SpeechSet(brief *SpeechInfo, verbose *Verbose) *OutputSpeechBuilder {
	b.brief = brief
	b.verbose = verbose
	return b
}

// Build method
func (b *OutputSpeechBuilder) Build() *OutputSpeech {
	if b.brief != nil && b.verbose != nil {
		return &OutputSpeech{
			Brief:   b.brief,
			Type:    OutputSpeechTypeSpeechSet,
			Verbose: b.verbose,
		}
	}
	if len(b.speeches) == 1 {
		return &OutputSpeech{
			Type:   OutputSpeechTypeSimpleSpeech,
			Values: b.speeches[0],
		}
	}
	return &OutputSpeech{
		Type:   OutputSpeechTypeSpeechList,
		Values: SpeechInfoArray(b.speeches),
	}
}
