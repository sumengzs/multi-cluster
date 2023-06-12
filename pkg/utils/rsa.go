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
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"k8s.io/client-go/util/keyutil"
)

const (
	DefaultLength = 2048
)

// RSAEncryptByPublicKey 加密
func RSAEncryptByPublicKey(plainText []byte, key *rsa.PublicKey) ([]byte, error) {
	partLen := key.N.BitLen()/8 - 11
	chunks := split(plainText, partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		buf, err := rsa.EncryptPKCS1v15(rand.Reader, key, chunk)
		if err != nil {
			return []byte{}, err
		}
		buffer.Write(buf)
	}

	return []byte(base64.RawStdEncoding.EncodeToString(buffer.Bytes())), nil
}

// RSADecryptByPrivateKey 解密
func RSADecryptByPrivateKey(cipherText []byte, key *rsa.PrivateKey) ([]byte, error) {
	partLen := key.N.BitLen() / 8
	raw, err := base64.RawStdEncoding.DecodeString(string(cipherText))
	chunks := split(raw, partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, key, chunk)
		if err != nil {
			return []byte{}, err
		}
		buffer.Write(decrypted)
	}

	return buffer.Bytes(), err
}
func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

func NewPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, DefaultLength)
}
func EncodePrivateKeyPEM(key *rsa.PrivateKey) []byte {
	block := pem.Block{
		Type:  keyutil.RSAPrivateKeyBlockType,
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	return pem.EncodeToMemory(&block)
}
