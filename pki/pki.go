package pki

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"net"

	"github.com/rancher/rke/hosts"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/util/cert"
)

type CertificatePKI struct {
	Certificate   *x509.Certificate
	Key           *rsa.PrivateKey
	Config        string
	Name          string
	CommonName    string
	OUName        string
	EnvName       string
	Path          string
	KeyEnvName    string
	KeyPath       string
	ConfigEnvName string
	ConfigPath    string
}

// StartCertificatesGeneration ...
func StartCertificatesGeneration(cpHosts []hosts.Host, workerHosts []hosts.Host, clusterDomain, localConfigPath string, KubernetesServiceIP net.IP) (map[string]CertificatePKI, error) {
	logrus.Infof("[certificates] Generating kubernetes certificates")
	certs, err := generateCerts(cpHosts, clusterDomain, localConfigPath, KubernetesServiceIP)
	if err != nil {
		return nil, err
	}
	return certs, nil
}

func generateCerts(cpHosts []hosts.Host, clusterDomain, localConfigPath string, KubernetesServiceIP net.IP) (map[string]CertificatePKI, error) {
	certs := make(map[string]CertificatePKI)
	// generate CA certificate and key
	logrus.Infof("[certificates] Generating CA kubernetes certificates")
	caCrt, caKey, err := generateCACertAndKey()
	if err != nil {
		return nil, err
	}
	logrus.Debugf("[certificates] CA Certificate: %s", string(cert.EncodeCertPEM(caCrt)))
	certs[CACertName] = CertificatePKI{
		Certificate: caCrt,
		Key:         caKey,
		Name:        CACertName,
		EnvName:     CACertENVName,
		KeyEnvName:  CAKeyENVName,
		Path:        CACertPath,
		KeyPath:     CAKeyPath,
	}

	// generate API certificate and key
	logrus.Infof("[certificates] Generating Kubernetes API server certificates")
	kubeAPIAltNames := GetAltNames(cpHosts, clusterDomain, KubernetesServiceIP)
	kubeAPICrt, kubeAPIKey, err := GenerateKubeAPICertAndKey(caCrt, caKey, kubeAPIAltNames)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("[certificates] Kube API Certificate: %s", string(cert.EncodeCertPEM(kubeAPICrt)))
	certs[KubeAPICertName] = CertificatePKI{
		Certificate: kubeAPICrt,
		Key:         kubeAPIKey,
		Name:        KubeAPICertName,
		EnvName:     KubeAPICertENVName,
		KeyEnvName:  KubeAPIKeyENVName,
		Path:        KubeAPICertPath,
		KeyPath:     KubeAPIKeyPath,
	}

	// generate Kube controller-manager certificate and key
	logrus.Infof("[certificates] Generating Kube Controller certificates")
	kubeControllerCrt, kubeControllerKey, err := generateClientCertAndKey(caCrt, caKey, KubeControllerCommonName, []string{})
	if err != nil {
		return nil, err
	}
	logrus.Debugf("[certificates] Kube Controller Certificate: %s", string(cert.EncodeCertPEM(kubeControllerCrt)))
	certs[KubeControllerName] = CertificatePKI{
		Certificate:   kubeControllerCrt,
		Key:           kubeControllerKey,
		Config:        getKubeConfigX509("https://127.0.0.1:6443", KubeControllerName, CACertPath, KubeControllerCertPath, KubeControllerKeyPath),
		Name:          KubeControllerName,
		CommonName:    KubeControllerCommonName,
		EnvName:       KubeControllerCertENVName,
		KeyEnvName:    KubeControllerKeyENVName,
		Path:          KubeControllerCertPath,
		KeyPath:       KubeControllerKeyPath,
		ConfigEnvName: KubeControllerConfigENVName,
		ConfigPath:    KubeControllerConfigPath,
	}

	// generate Kube scheduler certificate and key
	logrus.Infof("[certificates] Generating Kube Scheduler certificates")
	kubeSchedulerCrt, kubeSchedulerKey, err := generateClientCertAndKey(caCrt, caKey, KubeSchedulerCommonName, []string{})
	if err != nil {
		return nil, err
	}
	logrus.Debugf("[certificates] Kube Scheduler Certificate: %s", string(cert.EncodeCertPEM(kubeSchedulerCrt)))
	certs[KubeSchedulerName] = CertificatePKI{
		Certificate:   kubeSchedulerCrt,
		Key:           kubeSchedulerKey,
		Config:        getKubeConfigX509("https://127.0.0.1:6443", KubeSchedulerName, CACertPath, KubeSchedulerCertPath, KubeSchedulerKeyPath),
		Name:          KubeSchedulerName,
		CommonName:    KubeSchedulerCommonName,
		EnvName:       KubeSchedulerCertENVName,
		KeyEnvName:    KubeSchedulerKeyENVName,
		Path:          KubeSchedulerCertPath,
		KeyPath:       KubeSchedulerKeyPath,
		ConfigEnvName: KubeSchedulerConfigENVName,
		ConfigPath:    KubeSchedulerConfigPath,
	}

	// generate Kube Proxy certificate and key
	logrus.Infof("[certificates] Generating Kube Proxy certificates")
	kubeProxyCrt, kubeProxyKey, err := generateClientCertAndKey(caCrt, caKey, KubeProxyCommonName, []string{})
	if err != nil {
		return nil, err
	}
	logrus.Debugf("[certificates] Kube Proxy Certificate: %s", string(cert.EncodeCertPEM(kubeProxyCrt)))
	certs[KubeProxyName] = CertificatePKI{
		Certificate:   kubeProxyCrt,
		Key:           kubeProxyKey,
		Config:        getKubeConfigX509("https://127.0.0.1:6443", KubeProxyName, CACertPath, KubeProxyCertPath, KubeProxyKeyPath),
		Name:          KubeProxyName,
		CommonName:    KubeProxyCommonName,
		EnvName:       KubeProxyCertENVName,
		Path:          KubeProxyCertPath,
		KeyEnvName:    KubeProxyKeyENVName,
		KeyPath:       KubeProxyKeyPath,
		ConfigEnvName: KubeProxyConfigENVName,
		ConfigPath:    KubeProxyConfigPath,
	}

	// generate Kubelet certificate and key
	logrus.Infof("[certificates] Generating Node certificate")
	nodeCrt, nodeKey, err := generateClientCertAndKey(caCrt, caKey, KubeNodeCommonName, []string{KubeNodeOrganizationName})
	if err != nil {
		return nil, err
	}
	logrus.Debugf("[certificates] Node Certificate: %s", string(cert.EncodeCertPEM(kubeProxyCrt)))
	certs[KubeNodeName] = CertificatePKI{
		Certificate:   nodeCrt,
		Key:           nodeKey,
		Config:        getKubeConfigX509("https://127.0.0.1:6443", KubeNodeName, CACertPath, KubeNodeCertPath, KubeNodeKeyPath),
		Name:          KubeNodeName,
		CommonName:    KubeNodeCommonName,
		OUName:        KubeNodeOrganizationName,
		EnvName:       KubeNodeCertENVName,
		KeyEnvName:    KubeNodeKeyENVName,
		Path:          KubeNodeCertPath,
		KeyPath:       KubeNodeKeyPath,
		ConfigEnvName: KubeNodeConfigENVName,
		ConfigPath:    KubeNodeCommonName,
	}
	logrus.Infof("[certificates] Generating admin certificates and kubeconfig")
	kubeAdminCrt, kubeAdminKey, err := generateClientCertAndKey(caCrt, caKey, KubeAdminCommonName, []string{KubeAdminOrganizationName})
	if err != nil {
		return nil, err
	}
	logrus.Debugf("[certificates] Admin Certificate: %s", string(cert.EncodeCertPEM(kubeAdminCrt)))
	certs[KubeAdminCommonName] = CertificatePKI{
		Certificate: kubeAdminCrt,
		Key:         kubeAdminKey,
		Config: GetKubeConfigX509WithData(
			"https://"+cpHosts[0].IP+":6443",
			KubeAdminCommonName,
			string(cert.EncodeCertPEM(caCrt)),
			string(cert.EncodeCertPEM(kubeAdminCrt)),
			string(cert.EncodePrivateKeyPEM(kubeAdminKey))),
		CommonName:    KubeAdminCommonName,
		OUName:        KubeAdminOrganizationName,
		ConfigEnvName: KubeAdminConfigENVName,
		ConfigPath:    localConfigPath,
	}
	return certs, nil
}

func generateClientCertAndKey(caCrt *x509.Certificate, caKey *rsa.PrivateKey, commonName string, orgs []string) (*x509.Certificate, *rsa.PrivateKey, error) {
	rootKey, err := cert.NewPrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate private key for %s certificate: %v", commonName, err)
	}
	caConfig := cert.Config{
		CommonName:   commonName,
		Organization: orgs,
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	clientCert, err := cert.NewSignedCert(caConfig, rootKey, caCrt, caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate %s certificate: %v", commonName, err)
	}

	return clientCert, rootKey, nil
}

func GenerateKubeAPICertAndKey(caCrt *x509.Certificate, caKey *rsa.PrivateKey, altNames *cert.AltNames) (*x509.Certificate, *rsa.PrivateKey, error) {
	rootKey, err := cert.NewPrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate private key for kube-apiserver certificate: %v", err)
	}
	caConfig := cert.Config{
		CommonName: KubeAPICertName,
		Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		AltNames:   *altNames,
	}
	kubeCACert, err := cert.NewSignedCert(caConfig, rootKey, caCrt, caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate kube-apiserver certificate: %v", err)
	}

	return kubeCACert, rootKey, nil
}

func generateCACertAndKey() (*x509.Certificate, *rsa.PrivateKey, error) {
	rootKey, err := cert.NewPrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate private key for CA certificate: %v", err)
	}
	caConfig := cert.Config{
		CommonName: CACertName,
	}
	kubeCACert, err := cert.NewSelfSignedCACert(caConfig, rootKey)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate CA certificate: %v", err)
	}

	return kubeCACert, rootKey, nil
}

func GetAltNames(cpHosts []hosts.Host, clusterDomain string, KubernetesServiceIP net.IP) *cert.AltNames {
	ips := []net.IP{}
	dnsNames := []string{}
	for _, host := range cpHosts {
		ips = append(ips, net.ParseIP(host.IP))
		if host.IP != host.AdvertiseAddress {
			ips = append(ips, net.ParseIP(host.AdvertiseAddress))
		}
		dnsNames = append(dnsNames, host.AdvertisedHostname)
	}
	ips = append(ips, net.ParseIP("127.0.0.1"))
	ips = append(ips, KubernetesServiceIP)
	dnsNames = append(dnsNames, []string{
		"localhost",
		"kubernetes",
		"kubernetes.default",
		"kubernetes.default.svc",
		"kubernetes.default.svc." + clusterDomain,
	}...)
	return &cert.AltNames{
		IPs:      ips,
		DNSNames: dnsNames,
	}
}

func (c *CertificatePKI) ToEnv() []string {
	env := []string{
		c.CertToEnv(),
		c.KeyToEnv(),
	}
	if c.Config != "" {
		env = append(env, c.ConfigToEnv())
	}
	return env
}

func (c *CertificatePKI) CertToEnv() string {
	encodedCrt := cert.EncodeCertPEM(c.Certificate)
	return fmt.Sprintf("%s=%s", c.EnvName, string(encodedCrt))
}

func (c *CertificatePKI) KeyToEnv() string {
	encodedKey := cert.EncodePrivateKeyPEM(c.Key)
	return fmt.Sprintf("%s=%s", c.KeyEnvName, string(encodedKey))
}

func (c *CertificatePKI) ConfigToEnv() string {
	return fmt.Sprintf("%s=%s", c.ConfigEnvName, c.Config)
}
