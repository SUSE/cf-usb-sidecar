package tls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

const (
	// PrefixDriectory Directory Prefix where certs are generated
	PrefixDriectory = "/tmp"
	// PublicKeyNameSuffix is the name of the public key file to be generated
	PublicKeyNameSuffix = ".CAcert.pem"
	// PrivateKeyNameSuffix is the name of the private key file to be generated
	PrivateKeyNameSuffix = ".CAkey.pem"
	// CAName is the name of the CA and host for the generated key pair
	CommonName = "0.0.0.0"
	// Years for CA cert to be valid
	Years = 4
)

// Certificate contains information
type Certificate struct {
	Public  string
	Private string
}

// GenCert generate a self signed cert
func GenCert() (*Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	template := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			CommonName:   CommonName,
			Organization: []string{CommonName},
		},
		NotBefore:             now.Add(-5 * time.Minute).UTC(), // 5 minutes before now
		NotAfter:              now.AddDate(Years, 0, 0).UTC(),  // valid for Years
		BasicConstraintsValid: true,                            // required to set as a CA
		IsCA:         true, // set as a CA
		SubjectKeyId: []byte{1, 2, 3, 4},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}

	certFileName := CommonName + PublicKeyNameSuffix
	certFile, err := os.Create(filepath.Join(PrefixDriectory, certFileName))
	if err != nil {
		return nil, err
	}
	defer certFile.Close()
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})

	keyFileName := CommonName + PrivateKeyNameSuffix
	keyFile, err := os.OpenFile(filepath.Join(PrefixDriectory, keyFileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}
	defer keyFile.Close()
	pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	public := certFile.Name()
	private := keyFile.Name()

	return &Certificate{Public: public, Private: private}, nil
}
