/*
Copyright 2023 The Multi Cluster Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func TestRSAEncryption(t *testing.T) {
	plainText := []byte("Hello, World!")
	key, err := rsa.GenerateKey(rand.Reader, DefaultLength)
	if err != nil {
		t.Fatal("Failed to generate RSA key:", err)
	}

	// encrypt
	cipherText, err := RSAEncryptByPublicKey(plainText, &key.PublicKey)
	if err != nil {
		t.Fatal("Failed to encrypt:", err)
	}

	// decrypt
	decryptedText, err := RSADecryptByPrivateKey(cipherText, key)
	if err != nil {
		t.Fatal("Failed to decrypt:", err)
	}

	if string(decryptedText) != string(plainText) {
		t.Error("Decrypted text does not match original plain text")
	}
	t.Log("Encryption and decryption successful")
}

func TestEncodePrivateKeyPEM(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, DefaultLength)
	if err != nil {
		t.Fatal("Failed to generate RSA key:", err)
	}

	pemBytes := EncodePrivateKeyPEM(key)
	if len(pemBytes) == 0 {
		t.Error("Failed to encode private key to PEM")
	}

	t.Logf("Private key encoded to PEM:\n%s", string(pemBytes))
}
