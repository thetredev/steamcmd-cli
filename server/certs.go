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

type CertificateParameters struct {
	Company       string
	Country       string
	Province      string
	Locality      string
	StreetAddress string
	PostalCode    string
	ValidUntil    time.Time
	KeyLength     int
	Hostnames     []string
}

func createCertificate(params *CertificateParameters, keyUsage x509.KeyUsage) (*x509.Certificate, *rsa.PrivateKey, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	if err != nil {
		return nil, nil, err
	}

	certificate := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:  []string{params.Company},
			Country:       []string{params.Country},
			Province:      []string{params.Province},
			Locality:      []string{params.Locality},
			StreetAddress: []string{params.StreetAddress},
			PostalCode:    []string{params.PostalCode},
		},
		NotBefore:   time.Now(),
		NotAfter:    params.ValidUntil,
		IsCA:        true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    keyUsage,
	}

	if keyUsage&x509.KeyUsageCertSign > 0 {
		// Certificate is an (intermediate) authority
		certificate.BasicConstraintsValid = true
	}

	for _, hostname := range params.Hostnames {
		if ip := net.ParseIP(hostname); ip != nil {
			certificate.IPAddresses = append(certificate.IPAddresses, ip)
		} else {
			certificate.DNSNames = append(certificate.DNSNames, hostname)
		}
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, params.KeyLength)

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

func NewCertificateAuthority(params *CertificateParameters) (*Certificate, error) {
	certificate, privateKey, err := createCertificate(params, x509.KeyUsageDigitalSignature|x509.KeyUsageCertSign)

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

func IssueCertificate(ca *Certificate, params *CertificateParameters) (*IssuedCertificate, error) {
	certificate, privateKey, err := createCertificate(params, x509.KeyUsageDigitalSignature)

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
