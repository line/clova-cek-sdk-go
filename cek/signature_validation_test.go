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
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/line/clova-cek-sdk-go/cek"
)

func TestValidateSignature(t *testing.T) {
	// set publicKeyStr value for testing
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testPublicKey, err := ioutil.ReadFile(filepath.Join(dir, "testdata", "public.pem"))
	if err != nil {
		t.Fatal(err)
	}
	defer cek.SetPublicKeyStr(string(testPublicKey))()

	testBody := `{"version":"1.0","session":{"new":true,"sessionAttributes":{},"sessionId":"a29cfead-c5ba-474d-8745-6c1a6625f0c5","user":{"userId":"U399a1e08a8d474521fc4bbd8c7b4148f","accessToken":"XHapQasdfsdfFsdfasdflQQ7"}},"context":{"System":{"application":{"applicationId":"com.yourdomain.extension.pizzabot"},"user":{"userId":"U399a1e08a8d474521fc4bbd8c7b4148f","accessToken":"XHapQasdfsdfFsdfasdflQQ7"},"device":{"deviceId":"096e6b27-1717-33e9-b0a7-510a48658a9b","display":{"size":"l100","orientation":"landscape","dpi":96,"contentLayer":{"width":640,"height":360}}}}},"request":{"type":"LaunchRequest"}}`
	b, err := generateSignature([]byte(testBody))
	if err != nil {
		t.Fatal(err)
	}
	signature := base64.StdEncoding.EncodeToString(b)

	testCases := []struct {
		body         string
		signature    string
		wantResponse int
	}{
		// valid signature
		{
			body:         testBody,
			signature:    signature,
			wantResponse: http.StatusOK,
		},
		// invalid signature
		{
			body:         testBody,
			signature:    "invalidsignature",
			wantResponse: http.StatusBadRequest,
		},
		// valid signature, but body is modified
		{
			body:         strings.NewReplacer("LaunchRequest", "IntentRequest").Replace(testBody),
			signature:    signature,
			wantResponse: http.StatusBadRequest,
		},
	}

	ext := cek.NewExtension("com.yourdomain.extension.pizzabot")
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		_, err := ext.ParseRequest(r)
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
	for i, testCase := range testCases {
		req, err := http.NewRequest("POST", server.URL, bytes.NewReader([]byte(testCase.body)))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("SignatureCEK", testCase.signature)
		res, err := httpClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != testCase.wantResponse {
			t.Errorf("Status %d: %d", i, res.StatusCode)
		}
	}
}

func generateSignature(body []byte) ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	private, err := ioutil.ReadFile(filepath.Join(dir, "testdata", "private.pem"))
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(private)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	hash := crypto.SHA256.New()
	hash.Write(body)
	hashed := hash.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, key.(*rsa.PrivateKey), crypto.SHA256, hashed)
}
