package server

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

func createCertificate(keyUsage x509.KeyUsage, hostnames ...net.IP) (*x509.Certificate, *rsa.PrivateKey, error) {
	// TODO: Provide params via terminal input
	paramCompany := "steamcmd-cli"
	paramCountry := "US"
	paramProvince := ""
	paramLocality := "San Francisco"
	paramStreetAddress := "Somestreet"
	paramPostalCode := "94016"
	paramUntil := time.Date(2099, 1, 1, 1, 1, 1, 0, time.UTC)
	paramKeyLength := 4096

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	if err != nil {
		return nil, nil, err
	}

	certificate := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:  []string{paramCompany},
			Country:       []string{paramCountry},
			Province:      []string{paramProvince},
			Locality:      []string{paramLocality},
			StreetAddress: []string{paramStreetAddress},
			PostalCode:    []string{paramPostalCode},
		},
		NotBefore:   time.Now(),
		NotAfter:    paramUntil,
		IsCA:        true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    keyUsage,
	}

	if keyUsage&x509.KeyUsageCertSign > 0 {
		// Certificate is an (intermediate) authority
		certificate.BasicConstraintsValid = true
	}

	if len(hostnames) > 0 {
		certificate.IPAddresses = append(certificate.IPAddresses, hostnames...)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, paramKeyLength)

	if err != nil {
		return nil, nil, err
	}

	return certificate, privateKey, nil
}

type Certificate struct {
	X509          x509.Certificate
	PEM           *bytes.Buffer
	PrivateKey    *rsa.PrivateKey
	PrivateKeyPEM *bytes.Buffer
}

func NewCertificateAuthority() (*Certificate, error) {
	certificate, privateKey, err := createCertificate(x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign)

	if err != nil {
		return nil, err
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, certificate, certificate, &privateKey.PublicKey, privateKey)

	if err != nil {
		return nil, err
	}

	caPem := new(bytes.Buffer)
	pem.Encode(caPem, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	privateKeyPem := new(bytes.Buffer)
	pem.Encode(privateKeyPem, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return &Certificate{
		X509:          *certificate,
		PEM:           caPem,
		PrivateKey:    privateKey,
		PrivateKeyPEM: privateKeyPem,
	}, nil
}

type IssuedCertificate struct {
	CA          *Certificate
	Certificate *Certificate
}

func IssueCertificate(ca *Certificate, hostnames ...net.IP) (*IssuedCertificate, error) {
	// ipAddresses = []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback}
	certificate, privateKey, err := createCertificate(x509.KeyUsageDigitalSignature, hostnames...)

	if err != nil {
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, certificate, &ca.X509, &privateKey.PublicKey, ca.PrivateKey)

	if err != nil {
		return nil, err
	}

	certPem := new(bytes.Buffer)
	pem.Encode(certPem, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	privateKeyPem := new(bytes.Buffer)
	pem.Encode(privateKeyPem, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return &IssuedCertificate{
		CA: ca,
		Certificate: &Certificate{
			X509:          *certificate,
			PEM:           certPem,
			PrivateKey:    privateKey,
			PrivateKeyPEM: privateKeyPem,
		},
	}, nil
}
