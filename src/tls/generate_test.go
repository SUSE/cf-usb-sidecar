package tls

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenCert(t *testing.T) {
	cert, err := GenCert()
	assert.NoError(t, err)
	assert.NotNil(t, cert)
	// check the path was generated correctly
	certName := CommonName + PrivateKeyNameSuffix
	assert.Equal(t, filepath.Join(PrefixDriectory, certName), cert.Private)
	keyName := CommonName + PublicKeyNameSuffix
	assert.Equal(t, filepath.Join(PrefixDriectory, keyName), cert.Public)
	// check the certs exist at the path
	_, err = os.Stat(cert.Private)
	assert.NoError(t, err)
	_, err = os.Stat(cert.Public)
	assert.NoError(t, err)
	// remove the test generated certs
	assert.NoError(t, os.Remove(cert.Private))
	assert.NoError(t, os.Remove(cert.Public))
}
