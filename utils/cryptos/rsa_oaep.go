package cryptos

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// refer to https://github.com/java-crypto/cross_platform_crypto/tree/main/RsaEncryptionOaepSha256String

func RSAEncryptionOaepSha256ToBase64(publicKeyByte []byte, dataToEncryptBytes []byte) (base64ContentBytes []byte, err error) {
	block, _ := pem.Decode(publicKeyByte)
	if block == nil {
		return nil, errors.New("privateKey err or not exist")
	}
	parseResultPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey, ok := parseResultPublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("publicKey err")
	}
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, dataToEncryptBytes, nil)
	if err != nil {
		return nil, err
	}

	base64ContentBytes = make([]byte, base64.StdEncoding.EncodedLen(len(encryptedBytes)))
	base64.StdEncoding.Encode(base64ContentBytes, encryptedBytes)
	return base64ContentBytes, nil
}

func RSADecryptionOaepSha256FromBase64(privateKeyByte []byte, ciphertextBase64Bytes []byte) (contentBytes []byte, err error) {
	block, _ := pem.Decode(privateKeyByte)
	if block == nil {
		return nil, errors.New("privateKey err or not exist")
	}
	parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	privateKey, ok := parseResult.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("privateKey err")
	}
	decryptedBytes, err := privateKey.Decrypt(nil, base64Decoding(ciphertextBase64Bytes), &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		return nil, err
	}
	return decryptedBytes, nil
}

func base64Decoding(input []byte) []byte {
	data, _ := base64.StdEncoding.DecodeString(string(input))
	return data
}
