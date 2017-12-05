package cluster

const (
	DefaultClusterConfig = "cluster.yml"

	DefaultServiceClusterIPRange = "10.233.0.0/18"
	DefaultClusterCIDR           = "10.233.64.0/18"
	DefaultClusterDNSService     = "10.233.0.3"
	DefaultClusterDomain         = "cluster.local"
	DefaultClusterSSHKeyPath     = "~/.ssh/id_rsa"

	DefaultAuthStrategy = "x509"

	DefaultNetworkPlugin = "flannel"

	DefaultInfraContainerImage = "gcr.io/google_containers/pause-amd64:3.0"
	DefaultAplineImage         = "alpine:latest"
	DefaultNginxProxyImage     = "rancher/rke-nginx-proxy:0.1.0"
	DefaultCertDownloaderImage = "rancher/rke-cert-deployer:0.1.0"

	DefaultFlannelImage           = "quay.io/coreos/flannel:v0.9.1"
	DefaultFlannelCNIImage        = "quay.io/coreos/flannel-cni:v0.2.0"
	DefaultCalicoNodeImage        = "quay.io/calico/node:v2.6.2"
	DefaultCalicoCNIImage         = "quay.io/calico/cni:v1.11.0"
	DefaultCalicoControllersImage = "quay.io/calico/kube-controllers:v1.0.0"
	DefaultCanalNodeImage         = "quay.io/calico/node:v2.6.2"
	DefaultCanalCNIImage          = "quay.io/calico/cni:v1.11.0"
	DefaultCanalFlannelImage      = "quay.io/coreos/flannel:v0.9.1"

	DefaultKubeDNSImage           = "gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.5"
	DefaultDNSMasqImage           = "gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.5"
	DefaultKubeDNSSidecarImage    = "gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.5"
	DefaultKubeDNSAutoScalerImage = "gcr.io/google_containers/cluster-proportional-autoscaler-amd64:1.0.0"
)

func setDefaultIfEmptyMapValue(configMap map[string]string, key string, value string) {
	if _, ok := configMap[key]; !ok {
		configMap[key] = value
	}
}
func setDefaultIfEmpty(varName *string, defaultValue string) {
	if len(*varName) == 0 {
		*varName = defaultValue
	}
}
