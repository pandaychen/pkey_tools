package pkey_tools

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/cloudflare/cfssl/helpers"
	"github.com/gravitational/trace"
)

const (
	RSA_NAME = "RSA"
)

type PrivKey struct {
	key *rsa.PrivateKey
}

const (
	// TLSKeyAlgo is default TLS algo used for K8s X509 certs
	TLSKeyAlgo = "rsa"

	// TLSKeySize is default TLS key size used for K8s X509 certs
	TLSKeySize = 2048

	// RSAPrivateKeyPEMBlock is the name of the PEM block where private key is stored
	RSAPrivateKeyPEMBlock = "RSA PRIVATE KEY"

	// CertificatePEMBlock is the name of the PEM block where certificate is stored
	CertificatePEMBlock = "CERTIFICATE"

	// LicenseKeyPair is a name of the license key pair
	LicenseKeyPair = "license"

	// LoopbackIP is IP of the loopback interface
	LoopbackIP = "127.0.0.1"

	// LicenseKeyBits used when generating private key for license certificate
	LicenseKeyBits = 2048

	// LicenseOrg is the default name of license subject organization
	LicenseOrg = "xxxx.io"

	// LicenseTimeFormat represents format of expiration time in license payload
	LicenseTimeFormat = "2006-01-02 15:04:05"
)

// NewPrivateKey generates and returns private key
func NewRSAPrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, LicenseKeyBits)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return privateKey, nil
}

func NewPrivKey(keyPEM []byte) (*PrivKey, error) {
	if keyPEM != nil {
		key, err := helpers.ParsePrivateKeyPEMWithPassword(keyPEM, nil)
		if err != nil {
			return nil, trace.Wrap(err)
		}
		rkey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, trace.BadParameter("RSA keys supported, got %T", key)
		}
		return &PrivKey{key: rkey}, nil
	} else {
		rkey, err := NewRSAPrivateKey()
		if err != nil {
			return nil, trace.Wrap(err)
		}
		return &PrivKey{key: rkey}, nil
	}
}

// Algo returns the requested key algorithm represented as a string.
func (kr *PrivKey) Algo() string {
	return RSA_NAME
}

// Size returns the requested key size.
func (kr *PrivKey) Size() int {
	return kr.key.N.BitLen()
}

func (kr *PrivKey) GetPublicKey() *rsa.PublicKey {
	return &kr.key.PublicKey
}

// Generate generates a key as specified in the request. Currently,
// only ECDSA and RSA are supported.
func (kr *PrivKey) Generate() (crypto.PrivateKey, error) {
	return kr.key, nil
}

// SigAlgo returns an appropriate X.509 signature algorithm given the
// key request's type and size.
func (kr *PrivKey) SigAlgo() x509.SignatureAlgorithm {
	switch {
	case kr.Size() >= 4096:
		return x509.SHA512WithRSA
	case kr.Size() >= 3072:
		return x509.SHA384WithRSA
	case kr.Size() >= 2048:
		return x509.SHA256WithRSA
	default:
		return x509.SHA1WithRSA
	}
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

/*
func main() {
	prikey, _ := NewPrivKey(nil)
	fmt.Println(string(PrivateKeyToBytes(prikey.key)))
	fmt.Println(string(PublicKeyToBytes(prikey.GetPublicKey())))
	fmt.Println(prikey.Algo(),prikey.GetPublicKey())
}
*/
