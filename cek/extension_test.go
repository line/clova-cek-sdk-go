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
	"bytes"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/line/clova-cek-sdk-go/cek"
)

var testRequestBodies = []string{`{
  "version": "1.0",
  "session": {
    "new": false,
    "sessionAttributes": {},
    "sessionId": "a29cfead-c5ba-474d-8745-6c1a6625f0c5",
    "user": {
      "userId": "U399a1e08a8d474521fc4bbd8c7b4148f",
      "accessToken": "XHapQasdfsdfFsdfasdflQQ7"
    }
  },
  "context": {
    "System": {
      "application": {
        "applicationId": "com.yourdomain.extension.pizzabot"
      },
      "user": {
        "userId": "U399a1e08a8d474521fc4bbd8c7b4148f",
        "accessToken": "XHapQasdfsdfFsdfasdflQQ7"
      },
      "device": {
        "deviceId": "096e6b27-1717-33e9-b0a7-510a48658a9b",
        "display": {
          "size": "l100",
          "orientation": "landscape",
          "dpi": 96,
          "contentLayer": {
            "width": 640,
            "height": 360
          }
        }
      }
    }
  },
  "request": {
    "type": "EventRequest",
    "requestId": "f09874hiudf-sdf-4wku-flksdjfo4hjsdf",
    "timestamp": "2018-06-11T09:19:23Z",
    "event" : {
      "namespace":"ClovaSkill",
      "name":"SkillEnabled",
      "payload": null
    }
  }
}`, `{
  "version": "1.0",
  "session": {
    "new": false,
    "sessionAttributes": {},
    "sessionId": "a29cfead-c5ba-474d-8745-6c1a6625f0c5",
    "user": {
      "userId": "U399a1e08a8d474521fc4bbd8c7b4148f",
      "accessToken": "XHapQasdfsdfFsdfasdflQQ7"
    }
  },
  "context": {
    "System": {
      "application": {
        "applicationId": "com.yourdomain.extension.pizzabot"
      },
      "user": {
        "userId": "U399a1e08a8d474521fc4bbd8c7b4148f",
        "accessToken": "XHapQasdfsdfFsdfasdflQQ7"
      },
      "device": {
        "deviceId": "096e6b27-1717-33e9-b0a7-510a48658a9b",
        "display": {
          "size": "l100",
          "orientation": "landscape",
          "dpi": 96,
          "contentLayer": {
            "width": 640,
            "height": 360
          }
        }
      }
    }
  },
  "request": {
    "type": "IntentRequest",
    "intent": {
      "name": "OrderPizza",
      "slots": {
        "pizzaType": {
          "name": "pizzaType",
          "value": "ペパロニ"
        }
      }
    }
  }
}`, `{
  "version": "1.0",
  "session": {
    "new": true,
    "sessionAttributes": {},
    "sessionId": "a29cfead-c5ba-474d-8745-6c1a6625f0c5",
    "user": {
      "userId": "U399a1e08a8d474521fc4bbd8c7b4148f",
      "accessToken": "XHapQasdfsdfFsdfasdflQQ7"
    }
  },
  "context": {
    "System": {
      "application": {
        "applicationId": "com.yourdomain.extension.pizzabot"
      },
      "user": {
        "userId": "U399a1e08a8d474521fc4bbd8c7b4148f",
        "accessToken": "XHapQasdfsdfFsdfasdflQQ7"
      },
      "device": {
        "deviceId": "096e6b27-1717-33e9-b0a7-510a48658a9b",
        "display": {
          "size": "l100",
          "orientation": "landscape",
          "dpi": 96,
          "contentLayer": {
            "width": 640,
            "height": 360
          }
        }
      }
    }
  },
  "request": {
    "type": "LaunchRequest"
  }
}`, `{
  "version": "1.0",
  "session": {
    "new": false,
    "sessionAttributes": {},
    "sessionId": "a29cfead-c5ba-474d-8745-6c1a6625f0c5",
    "user": {
      "userId": "U399a1e08a8d474521fc4bbd8c7b4148f",
      "accessToken": "XHapQasdfsdfFsdfasdflQQ7"
    }
  },
  "context": {
    "System": {
      "application": {
        "applicationId": "com.yourdomain.extension.pizzabot"
      },
      "user": {
        "userId": "U399a1e08a8d474521fc4bbd8c7b4148f",
        "accessToken": "XHapQasdfsdfFsdfasdflQQ7"
      },
      "device": {
        "deviceId": "096e6b27-1717-33e9-b0a7-510a48658a9b",
        "display": {
          "size": "l100",
          "orientation": "landscape",
          "dpi": 96,
          "contentLayer": {
            "width": 640,
            "height": 360
          }
        }
      }
    }
  },
  "request": {
    "type": "SessionEndedRequest"
  }
}`}

func TestExtension(t *testing.T) {
	testCases := []struct {
		extension    *cek.Extension
		wantResponse int
	}{
		{
			extension:    cek.NewExtension("com.yourdomain.extension.pizzabot"),
			wantResponse: http.StatusBadRequest,
		},
		{
			extension:    cek.NewExtension("com.yourdomain.extension.pizzabot", cek.WithDebugMode),
			wantResponse: http.StatusOK,
		},
	}
	for i, testCase := range testCases {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			_, err := testCase.extension.ParseRequest(r)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}))
		defer server.Close()

		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		req, err := http.NewRequest("POST", server.URL, bytes.NewReader([]byte(testRequestBodies[0])))
		if err != nil {
			t.Fatal(err)
		}
		res, err := httpClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != testCase.wantResponse {
			t.Errorf("Status %d: %d", i, res.StatusCode)
		}
	}
}

func TestParseRequest(t *testing.T) {
	commonContext := &cek.Context{
		System: &cek.System{
			Application: &cek.Application{
				ApplicationID: "com.yourdomain.extension.pizzabot",
			},
			Device: &cek.Device{
				DeviceID: "096e6b27-1717-33e9-b0a7-510a48658a9b",
				Display: &cek.Display{
					ContentLayer: &cek.ContentLayer{Width: 640, Height: 360},
					DPI:          96,
					Orientation:  cek.OrientationLandscape,
					Size:         cek.SizeL100,
				},
			},
			User: &cek.User{
				UserID:      "U399a1e08a8d474521fc4bbd8c7b4148f",
				AccessToken: "XHapQasdfsdfFsdfasdflQQ7",
			},
		},
	}
	wantMessages := []*cek.RequestMessage{
		{
			Context: commonContext,
			Request: &cek.EventRequest{
				Event: &cek.Event{
					Name:      "SkillEnabled",
					Namespace: "ClovaSkill",
				},
				RequestID: "f09874hiudf-sdf-4wku-flksdjfo4hjsdf",
				Timestamp: "2018-06-11T09:19:23Z",
			},
			Session: &cek.Session{
				New:               false,
				SessionAttributes: map[string]string{},
				SessionID:         "a29cfead-c5ba-474d-8745-6c1a6625f0c5",
				User: &cek.User{
					UserID:      "U399a1e08a8d474521fc4bbd8c7b4148f",
					AccessToken: "XHapQasdfsdfFsdfasdflQQ7",
				},
			},
			Version: "1.0",
		},
		{
			Context: commonContext,
			Request: &cek.IntentRequest{
				Intent: &cek.Intent{
					Name: "OrderPizza",
					Slots: map[string]*cek.Slot{
						"pizzaType": {
							Name:  "pizzaType",
							Value: "ペパロニ",
						},
					},
				},
			},
			Session: &cek.Session{
				New:               false,
				SessionAttributes: map[string]string{},
				SessionID:         "a29cfead-c5ba-474d-8745-6c1a6625f0c5",
				User: &cek.User{
					UserID:      "U399a1e08a8d474521fc4bbd8c7b4148f",
					AccessToken: "XHapQasdfsdfFsdfasdflQQ7",
				},
			},
			Version: "1.0",
		},
		{
			Context: commonContext,
			Request: &cek.LaunchRequest{},
			Session: &cek.Session{
				New:               true,
				SessionAttributes: map[string]string{},
				SessionID:         "a29cfead-c5ba-474d-8745-6c1a6625f0c5",
				User: &cek.User{
					UserID:      "U399a1e08a8d474521fc4bbd8c7b4148f",
					AccessToken: "XHapQasdfsdfFsdfasdflQQ7",
				},
			},
			Version: "1.0",
		},
		{
			Context: commonContext,
			Request: &cek.SessionEndedRequest{},
			Session: &cek.Session{
				New:               false,
				SessionAttributes: map[string]string{},
				SessionID:         "a29cfead-c5ba-474d-8745-6c1a6625f0c5",
				User: &cek.User{
					UserID:      "U399a1e08a8d474521fc4bbd8c7b4148f",
					AccessToken: "XHapQasdfsdfFsdfasdflQQ7",
				},
			},
			Version: "1.0",
		},
	}

	var currentTestIdx int
	ext := cek.NewExtension("com.yourdomain.extension.pizzabot", cek.WithDebugMode)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		message, err := ext.ParseRequest(r)
		if err != nil {
			t.Fatal(err)
		}
		want := wantMessages[currentTestIdx]
		if !reflect.DeepEqual(message, want) {
			t.Errorf("Message %v; want %v", message, want)
		}
	}))
	defer server.Close()

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	for i, testRequestBody := range testRequestBodies {
		currentTestIdx = i
		req, err := http.NewRequest("POST", server.URL, bytes.NewReader([]byte(testRequestBody)))
		if err != nil {
			t.Fatal(err)
		}
		res, err := httpClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusOK {
			t.Errorf("Status %d: %d", i, res.StatusCode)
		}
	}
}
