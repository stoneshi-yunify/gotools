package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

var (
	caCertPath           = flag.String("ca-cert", "ca.crt", "CA certificate path")
	caKeyPath            = flag.String("ca-key", "ca.key", "CA key path")
	apiServerCertPath    = flag.String("apiserver-cert", "apiserver.crt", "apiserver certificate path")
	newApiServerCertPath = flag.String("new-apiserver-cert", "apiserver.crt", "new apiserver certificate path")
	newApiServerKeyPath  = flag.String("new-apiserver-key", "apiserver.key", "new apiserver key path")
	extraSANIPs          = flag.String("extra-san-ips", "", "extra Subject Alternative Names (SANs) to use for the API Server serving certificate. Can only be IP addresses. Separated by comma.")
)

func main() {
	flag.Parse()

	var err error
	var netIPs []net.IP
	ips := strings.Split(*extraSANIPs, ",")
	for _, ip := range ips {
		netIP := net.ParseIP(strings.TrimSpace(ip))
		if netIP != nil {
			netIPs = append(netIPs, netIP)
		}
	}
	if netIPs == nil || len(netIPs) == 0 {
		panic("invalid argument: extra-san-ips")
	}

	caPriv, err := parsePrivateKey(*caKeyPath)
	if err != nil {
		panic(err)
	}
	caCert, err := parseCertificate(*caCertPath)
	if err != nil {
		panic(err)
	}
	apiserverCert, err := parseCertificate(*apiServerCertPath)
	if err != nil {
		panic(err)
	}
	apiserverPriv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	templateObj := *apiserverCert
	template := &templateObj
	template.IPAddresses = append(template.IPAddresses, netIPs...)

	certBytes, err := x509.CreateCertificate(rand.Reader, template, caCert, &apiserverPriv.PublicKey, caPriv)
	if err != nil {
		panic(err)
	}

	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}
	err = writePemToFile(certBlock, *newApiServerCertPath)
	if err != nil {
		panic(err)
	}

	keyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(apiserverPriv),
	}
	err = writePemToFile(keyBlock, *newApiServerKeyPath)
	if err != nil {
		panic(err)
	}
	fmt.Println("Done")
}

func writePemToFile(block *pem.Block, path string) error {
	pemOut, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	err = pem.Encode(pemOut, block)
	if err != nil {
		return err
	}
	return pemOut.Close()
}

func parseCertificate(path string) (*x509.Certificate, error) {
	pemBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemBytes)
	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("not a certificate")
	}
	return x509.ParseCertificate(block.Bytes)
}

func parsePrivateKey(path string) (crypto.PrivateKey, error) {
	pemBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemBytes)
	if block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("not a rsa private key")
	}
	der := block.Bytes
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, fmt.Errorf("found unknown private key type in PKCS#8 wrapping")
		}
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}
	return nil, fmt.Errorf("failed to parse private key")
}
