/*
Copyright 2022 The Firefly Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package certs

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net"
	"path/filepath"
	"time"

	certutil "k8s.io/client-go/util/cert"
	"k8s.io/client-go/util/keyutil"
	"k8s.io/kube-openapi/pkg/util/sets"
)

const (
	// certificateBlockType is a possible value for pem.Block.Type.
	certificateBlockType = "CERTIFICATE"
	rsaKeySize           = 2048
	// Duration365d Certificate validity period
	Duration365d = time.Hour * 24 * 365
)

// NewPrivateKey returns a new private key.
var NewPrivateKey = GeneratePrivateKey

// GeneratePrivateKey Generate CA Private Key
func GeneratePrivateKey(keyType x509.PublicKeyAlgorithm) (crypto.Signer, error) {
	if keyType == x509.ECDSA {
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}

	return rsa.GenerateKey(rand.Reader, rsaKeySize)
}

// CertsConfig is a wrapper around certutil.Config extending it with PublicKeyAlgorithm.
type CertsConfig struct {
	certutil.Config
	NotAfter           *time.Time
	PublicKeyAlgorithm x509.PublicKeyAlgorithm
}

// EncodeCertPEM returns PEM-endcoded certificate data
func EncodeCertPEM(cert *x509.Certificate) []byte {
	block := pem.Block{
		Type:  certificateBlockType,
		Bytes: cert.Raw,
	}
	return pem.EncodeToMemory(&block)
}

// NewCertificateAuthority creates new certificate and private key for the certificate authority
func NewCertificateAuthority(config *CertsConfig) (*x509.Certificate, crypto.Signer, error) {
	key, err := NewPrivateKey(config.PublicKeyAlgorithm)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create private key while generating CA certificate %v", err)
	}

	cert, err := certutil.NewSelfSignedCACert(config.Config, key)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create self-signed CA certificate %v", err)
	}

	return cert, key, nil
}

// NewCACertAndKey The public and private keys of the root certificate are returned
func NewCACertAndKey(cn string) (*x509.Certificate, crypto.Signer, error) {
	certCfg := &CertsConfig{Config: certutil.Config{
		CommonName: cn,
	},
	}
	caCert, caKey, err := NewCertificateAuthority(certCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failure while generating CA certificate and key: %v", err)
	}

	return caCert, caKey, nil
}

// NewSignedCert creates a signed certificate using the given CA certificate and key
func NewSignedCert(cfg *CertsConfig, key crypto.Signer, caCert *x509.Certificate, caKey crypto.Signer, isCA bool) (*x509.Certificate, error) {
	serial, err := rand.Int(rand.Reader, new(big.Int).SetInt64(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	if len(cfg.CommonName) == 0 {
		return nil, errors.New("must specify a CommonName")
	}

	keyUsage := x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
	if isCA {
		keyUsage |= x509.KeyUsageCertSign
	}

	RemoveDuplicateAltNames(&cfg.AltNames)

	notAfter := time.Now().Add(Duration365d).UTC()
	if cfg.NotAfter != nil {
		notAfter = *cfg.NotAfter
	}

	certTmpl := x509.Certificate{
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		DNSNames:              cfg.AltNames.DNSNames,
		IPAddresses:           cfg.AltNames.IPs,
		SerialNumber:          serial,
		NotBefore:             caCert.NotBefore,
		NotAfter:              notAfter,
		KeyUsage:              keyUsage,
		ExtKeyUsage:           cfg.Usages,
		BasicConstraintsValid: true,
		IsCA:                  isCA,
	}
	certDERBytes, err := x509.CreateCertificate(rand.Reader, &certTmpl, caCert, key.Public(), caKey)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}

// RemoveDuplicateAltNames removes duplicate items in altNames.
func RemoveDuplicateAltNames(altNames *certutil.AltNames) {
	if altNames == nil {
		return
	}

	if altNames.DNSNames != nil {
		altNames.DNSNames = sets.NewString(altNames.DNSNames...).List()
	}

	ipsKeys := make(map[string]struct{})
	var ips []net.IP
	for _, one := range altNames.IPs {
		if _, ok := ipsKeys[one.String()]; !ok {
			ipsKeys[one.String()] = struct{}{}
			ips = append(ips, one)
		}
	}
	altNames.IPs = ips
}

// NewCertAndKey creates new certificate and key by passing the certificate authority certificate and key
func NewCertAndKey(caCert *x509.Certificate, caKey crypto.Signer, config *CertsConfig) (*x509.Certificate, crypto.Signer, error) {
	if len(config.Usages) == 0 {
		return nil, nil, errors.New("must specify at least one ExtKeyUsage")
	}

	key, err := NewPrivateKey(config.PublicKeyAlgorithm)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create private key %v", err)
	}

	cert, err := NewSignedCert(config, key, caCert, caKey, false)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to sign certificate. %v", err)
	}

	return cert, key, nil
}

// PathForKey returns the paths for the key given the path and basename.
func PathForKey(pkiPath, name string) string {
	return filepath.Join(pkiPath, fmt.Sprintf("%s.key", name))
}

// PathForCert returns the paths for the certificate given the path and basename.
func PathForCert(pkiPath, name string) string {
	return filepath.Join(pkiPath, fmt.Sprintf("%s.crt", name))
}

// WriteCert stores the given certificate at the given location
func WriteCert(pkiPath, name string, cert *x509.Certificate) error {
	if cert == nil {
		return errors.New("certificate cannot be nil when writing to file")
	}

	certificatePath := PathForCert(pkiPath, name)
	if err := certutil.WriteCert(certificatePath, EncodeCertPEM(cert)); err != nil {
		return fmt.Errorf("unable to write certificate to file %v", err)
	}

	return nil
}

// NewCertConfig create new CertConfig
func NewCertConfig(cn string, org []string, altNames certutil.AltNames, notAfter *time.Time) *CertsConfig {
	return &CertsConfig{
		Config: certutil.Config{
			CommonName:   cn,
			Organization: org,
			Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
			AltNames:     altNames,
		},
		NotAfter: notAfter,
	}
}

// GenCerts Create CA certificate and sign etcd karmada certificate.
func GenCerts(etcdServerCertCfg, etcdClientCertCfg, karmadaCertCfg, apiserverCertCfg, frontProxyClientCertCfg *CertsConfig) (map[string][]byte, error) {
	caCert, caKey, err := NewCACertAndKey("karmada")
	if err != nil {
		return nil, err
	}

	data := make(map[string][]byte)

	encodedCAKey, err := keyutil.MarshalPrivateKeyToPEM(caKey)
	if err != nil {
		return nil, err
	}
	encodedCACert := EncodeCertPEM(caCert)
	data["ca.key"] = encodedCAKey
	data["ca.crt"] = encodedCACert

	karmadaCert, karmadaKey, err := NewCertAndKey(caCert, caKey, karmadaCertCfg)
	if err != nil {
		return nil, err
	}
	encodedKarmadaKey, err := keyutil.MarshalPrivateKeyToPEM(karmadaKey)
	if err != nil {
		return nil, err
	}
	encodedKarmadaCert := EncodeCertPEM(karmadaCert)
	data["karmada.key"] = encodedKarmadaKey
	data["karmada.crt"] = encodedKarmadaCert

	apiserverCert, apiserverKey, err := NewCertAndKey(caCert, caKey, apiserverCertCfg)
	if err != nil {
		return nil, err
	}
	encodedApiserverKey, err := keyutil.MarshalPrivateKeyToPEM(apiserverKey)
	if err != nil {
		return nil, err
	}
	encodedApiserverCert := EncodeCertPEM(apiserverCert)
	data["apiserver.key"] = encodedApiserverKey
	data["apiserver.crt"] = encodedApiserverCert

	frontProxyCaCert, frontProxyCaKey, err := NewCACertAndKey("front-proxy-ca")
	if err != nil {
		return nil, err
	}
	encodedFrontProxyCaKey, err := keyutil.MarshalPrivateKeyToPEM(frontProxyCaKey)
	if err != nil {
		return nil, err
	}
	encodedFrontProxyCaCert := EncodeCertPEM(frontProxyCaCert)
	data["front-proxy-ca.key"] = encodedFrontProxyCaKey
	data["front-proxy-ca.crt"] = encodedFrontProxyCaCert

	frontProxyClientCert, frontProxyClientKey, err := NewCertAndKey(frontProxyCaCert, frontProxyCaKey, frontProxyClientCertCfg)
	if err != nil {
		return nil, err
	}
	encodedFrontProxyClientKey, err := keyutil.MarshalPrivateKeyToPEM(frontProxyClientKey)
	if err != nil {
		return nil, err
	}
	encodedFrontProxyClientCert := EncodeCertPEM(frontProxyClientCert)
	data["front-proxy-client.key"] = encodedFrontProxyClientKey
	data["front-proxy-client.crt"] = encodedFrontProxyClientCert

	etcdCaCert, etcdCaKey, err := NewCACertAndKey("etcd-ca")
	if err != nil {
		return nil, err
	}
	encodedEtcdCaKey, err := keyutil.MarshalPrivateKeyToPEM(etcdCaKey)
	if err != nil {
		return nil, err
	}
	encodedEtcdCaCert := EncodeCertPEM(etcdCaCert)
	data["etcd-ca.key"] = encodedEtcdCaKey
	data["etcd-ca.crt"] = encodedEtcdCaCert

	etcdServerCert, etcdServerKey, err := NewCertAndKey(etcdCaCert, etcdCaKey, etcdServerCertCfg)
	if err != nil {
		return nil, err
	}
	encodedEtcdServerKey, err := keyutil.MarshalPrivateKeyToPEM(etcdServerKey)
	if err != nil {
		return nil, err
	}
	encodedEtcdServerCert := EncodeCertPEM(etcdServerCert)
	data["etcd-server.key"] = encodedEtcdServerKey
	data["etcd-server.crt"] = encodedEtcdServerCert

	etcdClientCert, etcdClientKey, err := NewCertAndKey(etcdCaCert, etcdCaKey, etcdClientCertCfg)
	if err != nil {
		return nil, err
	}
	encodedEtcdClientKey, err := keyutil.MarshalPrivateKeyToPEM(etcdClientKey)
	if err != nil {
		return nil, err
	}
	encodedEtcdClientCert := EncodeCertPEM(etcdClientCert)
	data["etcd-client.key"] = encodedEtcdClientKey
	data["etcd-client.crt"] = encodedEtcdClientCert

	return data, nil
}
