package system

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
)

type Certificate struct {
	CertFile string
	KeyFile  string
}

func (cert Certificate) GenerateIfNotExist() error {
	certFileExists, err := core.FileExists(cert.CertFile)
	if err != nil {
		return err
	}
	keyFileExists, err := core.FileExists(cert.KeyFile)
	if err != nil {
		return err
	}
	if !(certFileExists && keyFileExists) {
		err := cert.Generate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Certificate) Generate() error {
	now := time.Now()
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2024),
		Subject: pkix.Name{
			OrganizationalUnit: []string{"IPCManView"},
			CommonName:         "IPCManView",
			Country:            []string{"US"},
		},
		NotBefore:   now,
		NotAfter:    now.AddDate(10, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &certPrivKey.PublicKey, certPrivKey)
	if err != nil {
		return err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err := os.WriteFile(c.CertFile, certPEM.Bytes(), 0600); err != nil {
		return err
	}

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	if err := os.WriteFile(c.KeyFile, certPrivKeyPEM.Bytes(), 0600); err != nil {
		return err
	}

	return nil
}
