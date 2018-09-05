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

package cek_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/line/clova-cek-sdk-go/cek"
)

func TestResponseBuilder(t *testing.T) {
	testMessages := []*cek.ResponseMessage{
		cek.NewResponseBuilder().
			OutputSpeech(cek.NewOutputSpeechBuilder().
				AddSpeechText("Hi, nice to meet you", cek.SpeechInfoLangEN).
				Build()).
			Build(),
		cek.NewResponseBuilder().
			OutputSpeech(cek.NewOutputSpeechBuilder().
				AddSpeechText("歌を歌ってみます。", cek.SpeechInfoLangJA).
				AddSpeechURL("https://DUMMY_DOMAIN/song.mp3").
				Build()).
			ShouldEndSession(true).
			Build(),
		cek.NewResponseBuilder().
			OutputSpeech(cek.NewOutputSpeechBuilder().
				SpeechSet(
					&cek.SpeechInfo{
						Lang:  cek.SpeechInfoLangJA,
						Type:  cek.SpeechInfoTypePlainText,
						Value: "天気予報です。",
					},
					&cek.Verbose{
						Type: cek.OutputSpeechVerboseTypeSpeechList,
						Values: cek.SpeechInfoArray([]*cek.SpeechInfo{
							{
								Lang:  cek.SpeechInfoLangJA,
								Type:  cek.SpeechInfoTypePlainText,
								Value: "週末まで全国に梅雨…猛暑和らぐ。",
							},
							{
								Lang:  cek.SpeechInfoLangJA,
								Type:  cek.SpeechInfoTypePlainText,
								Value: "明日全国的に梅雨…ところによって局地的に激しい雨に注意。",
							},
						}),
					}).
				Build()).
			ShouldEndSession(true).
			Build(),
		cek.NewResponseBuilder().
			SessionAttributes(map[string]string{
				"RequestedIntent": "OrderPizza",
				"pizzaType":       "ペパロニピザ",
			}).
			OutputSpeech(cek.NewOutputSpeechBuilder().
				AddSpeechText("何枚注文しますか?", cek.SpeechInfoLangJA).
				Build()).
			Build(),
		cek.NewResponseBuilder().
			SessionAttributes(map[string]string{
				"RequestedIntent": "OrderPizza",
				"pizzaType":       "ペパロニピザ",
			}).
			OutputSpeech(cek.NewOutputSpeechBuilder().
				AddSpeechText("何枚注文しますか?", cek.SpeechInfoLangJA).
				Build()).
			Reprompt(cek.NewOutputSpeechBuilder().
				AddSpeechText("お言葉がなければ、注文をキャンセルしてよろしいですか?", cek.SpeechInfoLangJA).
				Build()).
			Build(),
	}
	wantBodies := []string{`{
  "version": "1.0",
  "sessionAttributes": {},
  "response": {
    "outputSpeech": {
      "type": "SimpleSpeech",
      "values": {
          "type": "PlainText",
          "lang": "en",
          "value": "Hi, nice to meet you"
      }
    },
    "card": {},
    "directives": [],
    "shouldEndSession": false
  }
}`, `{
  "version": "1.0",
  "sessionAttributes": {},
  "response": {
    "outputSpeech": {
      "type": "SpeechList",
      "values": [
        {
          "type": "PlainText",
          "lang": "ja",
          "value": "歌を歌ってみます。"
        },
        {
          "type": "URL",
          "lang": "" ,
          "value": "https://DUMMY_DOMAIN/song.mp3"
        }
      ]
    },
    "card": {},
    "directives": [],
    "shouldEndSession": true
  }
}
`, `{
  "version": "1.0",
  "sessionAttributes": {},
  "response": {
    "outputSpeech": {
      "type": "SpeechSet",
      "brief": {
        "type": "PlainText",
        "lang": "ja",
        "value": "天気予報です。"
      },
      "verbose": {
        "type": "SpeechList",
        "values": [
          {
              "type": "PlainText",
              "lang": "ja",
              "value": "週末まで全国に梅雨…猛暑和らぐ。"
          },
          {
              "type": "PlainText",
              "lang": "ja",
              "value": "明日全国的に梅雨…ところによって局地的に激しい雨に注意。"
          }
        ]
      }
    },
    "card": {},
    "directives": [],
    "shouldEndSession": true
  }
}`, `{
  "version": "1.0",
  "sessionAttributes": {
    "RequestedIntent": "OrderPizza",
    "pizzaType": "ペパロニピザ"
  },
  "response": {
    "outputSpeech": {
      "type": "SimpleSpeech",
      "values": {
          "type": "PlainText",
          "lang": "ja",
          "value": "何枚注文しますか?"
      }
    },
    "card": {},
    "directives": [],
    "shouldEndSession": false
  }
}`, `{
  "version": "1.0",
  "sessionAttributes": {
    "RequestedIntent": "OrderPizza",
    "pizzaType": "ペパロニピザ"
  },
  "response": {
    "outputSpeech": {
      "type": "SimpleSpeech",
      "values": {
          "type": "PlainText",
          "lang": "ja",
          "value": "何枚注文しますか?"
      }
    },
    "card": {},
    "directives": [],
    "reprompt" : {
      "outputSpeech" : {
        "type" : "SimpleSpeech",
        "values" : {
          "type" : "PlainText",
          "lang" : "ja",
          "value" : "お言葉がなければ、注文をキャンセルしてよろしいですか?"
        }
      }
    },
    "shouldEndSession": false
  }
}`}

	for i, message := range testMessages {
		gotBody, err := json.Marshal(message)
		if err != nil {
			t.Fatal(err)
		}

		var got, want interface{}
		if err := json.Unmarshal([]byte(gotBody), &got); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal([]byte(wantBodies[i]), &want); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Error %d: %v, want %v", i, got, want)
		}
	}
}
