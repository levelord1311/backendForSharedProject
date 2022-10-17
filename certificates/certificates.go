package certificates

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"
)

func generateSerialNumber() (*big.Int, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}
	return serialNumber, nil
}

func getCert(serialNumber *big.Int) x509.Certificate {
	return x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Probably Volodya"},
		},
		DNSNames:  []string{"localhost"},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(3 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
}

func generateCertTemplate(privateKey *ecdsa.PrivateKey) []byte {
	serialNumber, err := generateSerialNumber()
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	template := getCert(serialNumber)

	//certificate is self-signed (&template is both template and parent args)
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	return derBytes
}

func serializeCert(derBytes []byte) {
	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if pemCert == nil {
		log.Fatal("Failed to encode certificate to PEM")
	}
	if err := os.WriteFile("certificates/cert.pem", pemCert, 0644); err != nil {
		log.Fatal(err)
	}
	log.Println("wrote cert.pem")
}

func serializeKey(privateKey *ecdsa.PrivateKey) {
	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Unable to marshal private key: %v", err)
	}
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	if pemKey == nil {
		log.Fatal("Failed to encode key to PEM")
	}
	if err := os.WriteFile("certificates/key.pem", pemKey, 0600); err != nil {
		log.Fatal(err)
	}
	log.Println("wrote key.pem")
}

func CreateCertAndKey() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	certBytes := generateCertTemplate(privateKey)
	serializeCert(certBytes)
	serializeKey(privateKey)
}
