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
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

var publicKeyStr = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwiMvQNKD/WQcX9KiWNMb
nSR+dJYTWL6TmqqwWFia69TyiobVIfGfxFSefxYyMTcFznoGCpg8aOCAkMxUH58N
0/UtWWvfq0U5FQN9McE3zP+rVL3Qul9fbC2mxvazxpv5KT7HEp780Yew777cVPUv
3+I73z2t0EHnkwMesmpUA/2Rp8fW8vZE4jfiTRm5vSVmW9F37GC5TEhPwaiIkIin
KCrH0rXbfe3jNWR7qKOvVDytcWgRHJqRUuWhwJuAnuuqLvqTyAawqEslhKZ5t+1Z
0GN8b2zMENSuixa1M9K0ZKUw3unzHpvgBlYmXRGPTSuq/EaGYWyckYz8CBq5Lz2Q
UwIDAQAB
-----END PUBLIC KEY-----
`

func validateSignature(signature string, body []byte) error {
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		return errors.New("failed to parse PEM block containing the public key")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	hash := crypto.SHA256.New()
	if _, err := hash.Write(body); err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(key.(*rsa.PublicKey), crypto.SHA256, hash.Sum(nil), sig)
}
