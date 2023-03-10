/*
Copyright 2021 Intel(R)
SPDX-License-Identifier: Apache-2.0
*/

package sgx

import (
	"encoding/pem"
	"errors"
	"runtime"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

const (
	rsaKeySize = 3072
)

// NewPrivateKey generate an rsa private key
func NewPrivateKey() (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		return nil, err
	}

	runtime.SetFinalizer(key, func(k *rsa.PrivateKey) {
		// Zero key after usage
		*k = rsa.PrivateKey{}
	})
	return key, nil
}

// EncodeKey returns PEM encoding of give private key
func EncodeKey(key *rsa.PrivateKey) []byte {
	if key == nil {
		return []byte{}
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
}

// EncodePublicKey returns PEM encoding of given public key
func EncodePublicKey(key interface{}) ([]byte, error) {
	if key == nil {
		return []byte{}, nil
	}
	bytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: bytes,
	}), nil
}

// DecodeKey returns the decoded private key of given encodedKey
func DecodeKey(encodedKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(encodedKey)

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	wipe(block.Bytes)
	if err != nil {
		return nil, err
	}

	runtime.SetFinalizer(key, func(k *rsa.PrivateKey) {
		// Zero key after usage
		*k = rsa.PrivateKey{}
	})

	return key, nil
}

// EncodeCert returns PEM encoding of given cert
func EncodeCert(cert *x509.Certificate) []byte {
	if cert == nil {
		return []byte{}
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
}

// DecodeCert return the decoded certificate of given encodedCert
func DecodeCert(encodedCert []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(encodedCert)

	cert, err := x509.ParseCertificate(block.Bytes)
	wipe(block.Bytes)

	if err != nil {
		return nil, err
	}
	runtime.SetFinalizer(cert, func(c *x509.Certificate) {
		*c = x509.Certificate{}
	})

	return cert, nil
}

// DecodeCert return the decoded csr of given encodedCertRequest
func DecodeCertRequest(encodedCertRequest []byte) (*x509.CertificateRequest, error) {
	block, _ := pem.Decode(encodedCertRequest)
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return nil, errors.New("PEM block is not a CERTIFICATE REQUEST")
	}
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, err
	}
	return csr, nil
}

func wipe(arr []byte) {
	for i := range arr {
		arr[i] = 0
	}
}
