package types

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiserverv1alpha1 "k8s.io/apiserver/pkg/apis/apiserver/v1alpha1"
	auditv1 "k8s.io/apiserver/pkg/apis/audit/v1"
	configv1 "k8s.io/apiserver/pkg/apis/config/v1"
)

type RancherKubernetesEngineConfig struct {
	// Kubernetes nodes
	Nodes []RKEConfigNode `yaml:"nodes,omitempty" json:"nodes,omitempty"`
	// Kubernetes components
	Services RKEConfigServices `yaml:"services,omitempty" json:"services,omitempty"`
	// Network configuration used in the kubernetes cluster (flannel, calico)
	Network NetworkConfig `yaml:"network,omitempty" json:"network,omitempty"`
	// Authentication configuration used in the cluster (default: x509)
	Authentication AuthnConfig `yaml:"authentication,omitempty" json:"authentication,omitempty"`
	// YAML manifest for user provided addons to be deployed on the cluster
	Addons string `yaml:"addons,omitempty" json:"addons,omitempty"`
	// List of urls or paths for addons
	AddonsInclude []string `yaml:"addons_include,omitempty" json:"addonsInclude,omitempty"`
	// List of images used internally for proxy, cert download and kubedns
	SystemImages RKESystemImages `yaml:"system_images,omitempty" json:"systemImages,omitempty"`
	// SSH Private Key Path
	SSHKeyPath string `yaml:"ssh_key_path,omitempty" json:"sshKeyPath,omitempty" norman:"nocreate,noupdate"`
	// SSH Certificate Path
	SSHCertPath string `yaml:"ssh_cert_path,omitempty" json:"sshCertPath,omitempty" norman:"nocreate,noupdate"`
	// SSH Agent Auth enable
	SSHAgentAuth bool `yaml:"ssh_agent_auth,omitempty" json:"sshAgentAuth,omitempty"`
	// Authorization mode configuration used in the cluster
	Authorization AuthzConfig `yaml:"authorization,omitempty" json:"authorization,omitempty"`
	// Enable/disable strict docker version checking
	IgnoreDockerVersion *bool `yaml:"ignore_docker_version,omitempty" json:"ignoreDockerVersion,omitempty" norman:"default=true"`
	// Enable/disable using cri-dockerd
	EnableCRIDockerd *bool `yaml:"enable_cri_dockerd,omitempty" json:"enableCRIDockerd,omitempty" norman:"default=false"`
	// Kubernetes version to use (if kubernetes image is specified, image version takes precedence)
	Version string `yaml:"kubernetes_version,omitempty" json:"kubernetesVersion,omitempty"`
	// List of private registries and their credentials
	PrivateRegistries []PrivateRegistry `yaml:"private_registries,omitempty" json:"privateRegistries,omitempty"`
	// Ingress controller used in the cluster
	Ingress IngressConfig `yaml:"ingress,omitempty" json:"ingress,omitempty"`
	// Cluster Name used in the kube config
	ClusterName string `yaml:"cluster_name,omitempty" json:"clusterName,omitempty"`
	// Cloud Provider options
	CloudProvider CloudProvider `yaml:"cloud_provider,omitempty" json:"cloudProvider,omitempty"`
	// kubernetes directory path
	PrefixPath string `yaml:"prefix_path,omitempty" json:"prefixPath,omitempty"`
	// kubernetes directory path for windows
	WindowsPrefixPath string `yaml:"win_prefix_path,omitempty" json:"winPrefixPath,omitempty"`
	// Timeout in seconds for status check on addon deployment jobs
	AddonJobTimeout int `yaml:"addon_job_timeout,omitempty" json:"addonJobTimeout,omitempty" norman:"default=45"`
	// Bastion/Jump Host configuration
	BastionHost BastionHost `yaml:"bastion_host,omitempty" json:"bastionHost,omitempty"`
	// Monitoring Config
	Monitoring MonitoringConfig `yaml:"monitoring,omitempty" json:"monitoring,omitempty"`
	// RestoreCluster flag
	Restore RestoreConfig `yaml:"restore,omitempty" json:"restore,omitempty"`
	// Rotating Certificates Option
	RotateCertificates *RotateCertificates `yaml:"rotate_certificates,omitempty" json:"rotateCertificates,omitempty"`
	// Rotate Encryption Key Option
	RotateEncryptionKey bool `yaml:"rotate_encryption_key,omitempty" json:"rotateEncryptionKey,omitempty"`
	// DNS Config
	DNS *DNSConfig `yaml:"dns,omitempty" json:"dns,omitempty"`
	// Upgrade Strategy for the cluster
	UpgradeStrategy *NodeUpgradeStrategy `yaml:"upgrade_strategy,omitempty" json:"upgradeStrategy,omitempty"`
}

func (r *RancherKubernetesEngineConfig) ObjClusterName() string {
	return r.ClusterName
}

type NodeUpgradeStrategy struct {
	// MaxUnavailableWorker input can be a number of nodes or a percentage of nodes (example, max_unavailable_worker: 2 OR max_unavailable_worker: 20%)
	MaxUnavailableWorker string `yaml:"max_unavailable_worker,omitempty" json:"maxUnavailableWorker,omitempty" norman:"min=1,default=10%"`
	// MaxUnavailableControlplane input can be a number of nodes or a percentage of nodes
	MaxUnavailableControlplane string          `yaml:"max_unavailable_controlplane,omitempty" json:"maxUnavailableControlplane,omitempty" norman:"min=1,default=1"`
	Drain                      *bool           `yaml:"drain,omitempty" json:"drain,omitempty"`
	DrainInput                 *NodeDrainInput `yaml:"node_drain_input,omitempty" json:"nodeDrainInput,omitempty"`
}

type BastionHost struct {
	// Address of Bastion Host
	Address string `yaml:"address,omitempty" json:"address,omitempty"`
	// SSH Port of Bastion Host
	Port string `yaml:"port,omitempty" json:"port,omitempty" jsonschema:"oneof_type=string;number"`
	// ssh User to Bastion Host
	User string `yaml:"user,omitempty" json:"user,omitempty"`
	// SSH Agent Auth enable
	SSHAgentAuth bool `yaml:"ssh_agent_auth,omitempty" json:"sshAgentAuth,omitempty"`
	// SSH Private Key
	SSHKey string `yaml:"ssh_key,omitempty" json:"sshKey,omitempty" norman:"type=password"`
	// SSH Private Key Path
	SSHKeyPath string `yaml:"ssh_key_path,omitempty" json:"sshKeyPath,omitempty"`
	// SSH Certificate
	SSHCert string `yaml:"ssh_cert,omitempty" json:"sshCert,omitempty"`
	// SSH Certificate Path
	SSHCertPath string `yaml:"ssh_cert_path,omitempty" json:"sshCertPath,omitempty"`
	// Ignore proxy environment variables
	IgnoreProxyEnvVars bool `yaml:"ignore_proxy_env_vars,omitempty" json:"ignoreProxyEnvVars,omitempty"`
}

type PrivateRegistry struct {
	// URL for the registry
	URL string `yaml:"url,omitempty" json:"url,omitempty"`
	// User name for registry acces
	User string `yaml:"user,omitempty" json:"user,omitempty"`
	// Password for registry access
	Password string `yaml:"password,omitempty" json:"password,omitempty" norman:"type=password"`
	// Default registry
	IsDefault bool `yaml:"is_default,omitempty" json:"isDefault,omitempty"`
}

type RKESystemImages struct {
	// etcd image
	Etcd string `yaml:"etcd,omitempty" json:"etcd,omitempty"`
	// Alpine image
	Alpine string `yaml:"alpine,omitempty" json:"alpine,omitempty"`
	// rke-nginx-proxy image
	NginxProxy string `yaml:"nginx_proxy,omitempty" json:"nginxProxy,omitempty"`
	// rke-cert-deployer image
	CertDownloader string `yaml:"cert_downloader,omitempty" json:"certDownloader,omitempty"`
	// rke-service-sidekick image
	KubernetesServicesSidecar string `yaml:"kubernetes_services_sidecar,omitempty" json:"kubernetesServicesSidecar,omitempty"`
	// KubeDNS image
	KubeDNS string `yaml:"kubedns,omitempty" json:"kubedns,omitempty"`
	// DNSMasq image
	DNSmasq string `yaml:"dnsmasq,omitempty" json:"dnsmasq,omitempty"`
	// KubeDNS side car image
	KubeDNSSidecar string `yaml:"kubedns_sidecar,omitempty" json:"kubednsSidecar,omitempty"`
	// KubeDNS autoscaler image
	KubeDNSAutoscaler string `yaml:"kubedns_autoscaler,omitempty" json:"kubednsAutoscaler,omitempty"`
	// CoreDNS image
	CoreDNS string `yaml:"coredns,omitempty" json:"coredns,omitempty"`
	// CoreDNS autoscaler image
	CoreDNSAutoscaler string `yaml:"coredns_autoscaler,omitempty" json:"corednsAutoscaler,omitempty"`
	// Nodelocal image
	Nodelocal string `yaml:"nodelocal,omitempty" json:"nodelocal,omitempty"`
	// Kubernetes image
	Kubernetes string `yaml:"kubernetes,omitempty" json:"kubernetes,omitempty"`
	// Flannel image
	Flannel string `yaml:"flannel,omitempty" json:"flannel,omitempty"`
	// Flannel CNI image
	FlannelCNI string `yaml:"flannel_cni,omitempty" json:"flannelCni,omitempty"`
	// Calico Node image
	CalicoNode string `yaml:"calico_node,omitempty" json:"calicoNode,omitempty"`
	// Calico CNI image
	CalicoCNI string `yaml:"calico_cni,omitempty" json:"calicoCni,omitempty"`
	// Calico Controllers image
	CalicoControllers string `yaml:"calico_controllers,omitempty" json:"calicoControllers,omitempty"`
	// Calicoctl image
	CalicoCtl string `yaml:"calico_ctl,omitempty" json:"calicoCtl,omitempty"`
	//CalicoFlexVol image
	CalicoFlexVol string `yaml:"calico_flexvol,omitempty" json:"calicoFlexVol,omitempty"`
	// Canal Node Image
	CanalNode string `yaml:"canal_node,omitempty" json:"canalNode,omitempty"`
	// Canal CNI image
	CanalCNI string `yaml:"canal_cni,omitempty" json:"canalCni,omitempty"`
	// Canal Controllers Image needed for Calico/Canal v3.14.0+
	CanalControllers string `yaml:"canal_controllers,omitempty" json:"canalControllers,omitempty"`
	//CanalFlannel image
	CanalFlannel string `yaml:"canal_flannel,omitempty" json:"canalFlannel,omitempty"`
	//CanalFlexVol image
	CanalFlexVol string `yaml:"canal_flexvol,omitempty" json:"canalFlexVol,omitempty"`
	//Weave Node image
	WeaveNode string `yaml:"weave_node,omitempty" json:"weaveNode,omitempty"`
	// Weave CNI image
	WeaveCNI string `yaml:"weave_cni,omitempty" json:"weaveCni,omitempty"`
	// Pod infra container image
	PodInfraContainer string `yaml:"pod_infra_container,omitempty" json:"podInfraContainer,omitempty"`
	// Ingress Controller image
	Ingress string `yaml:"ingress,omitempty" json:"ingress,omitempty"`
	// Ingress Controller Backend image
	IngressBackend string `yaml:"ingress_backend,omitempty" json:"ingressBackend,omitempty"`
	// Metrics Server image
	MetricsServer string `yaml:"metrics_server,omitempty" json:"metricsServer,omitempty"`
	// Pod infra container image for Windows
	WindowsPodInfraContainer string `yaml:"windows_pod_infra_container,omitempty" json:"windowsPodInfraContainer,omitempty"`
	// Cni deployer container image for Cisco ACI
	AciCniDeployContainer string `yaml:"aci_cni_deploy_container,omitempty" json:"aciCniDeployContainer,omitempty"`
	// host container image for Cisco ACI
	AciHostContainer string `yaml:"aci_host_container,omitempty" json:"aciHostContainer,omitempty"`
	// opflex agent container image for Cisco ACI
	AciOpflexContainer string `yaml:"aci_opflex_container,omitempty" json:"aciOpflexContainer,omitempty"`
	// mcast daemon container image for Cisco ACI
	AciMcastContainer string `yaml:"aci_mcast_container,omitempty" json:"aciMcastContainer,omitempty"`
	// OpenvSwitch container image for Cisco ACI
	AciOpenvSwitchContainer string `yaml:"aci_ovs_container,omitempty" json:"aciOvsContainer,omitempty"`
	// Controller container image for Cisco ACI
	AciControllerContainer string `yaml:"aci_controller_container,omitempty" json:"aciControllerContainer,omitempty"`
	// GBP Server container image for Cisco ACI
	AciGbpServerContainer string `yaml:"aci_gbp_server_container,omitempty" json:"aciGbpServerContainer,omitempty"`
	// Opflex Server container image for Cisco ACI
	AciOpflexServerContainer string `yaml:"aci_opflex_server_container,omitempty" json:"aciOpflexServerContainer,omitempty"`
}

type RKEConfigNode struct {
	// Name of the host provisioned via docker machine
	NodeName string `yaml:"nodeName,omitempty" json:"nodeName,omitempty" norman:"type=reference[node]"`
	// IP or FQDN that is fully resolvable and used for SSH communication
	Address string `yaml:"address,omitempty" json:"address,omitempty"`
	// Port used for SSH communication
	Port string `yaml:"port,omitempty" json:"port,omitempty" jsonschema:"oneof_type=string;number"`
	// Optional - Internal address that will be used for components communication
	InternalAddress string `yaml:"internal_address,omitempty" json:"internalAddress,omitempty"`
	// Node role in kubernetes cluster (controlplane, worker, or etcd)
	Role []string `yaml:"role,omitempty" json:"role,omitempty" norman:"type=array[enum],options=etcd|worker|controlplane"`
	// Optional - Hostname of the node
	HostnameOverride string `yaml:"hostname_override,omitempty" json:"hostnameOverride,omitempty"`
	// SSH usesr that will be used by RKE
	User string `yaml:"user,omitempty" json:"user,omitempty"`
	// Optional - Docker socket on the node that will be used in tunneling
	DockerSocket string `yaml:"docker_socket,omitempty" json:"dockerSocket,omitempty"`
	// SSH Agent Auth enable
	SSHAgentAuth bool `yaml:"ssh_agent_auth,omitempty" json:"sshAgentAuth,omitempty"`
	// SSH Private Key
	SSHKey string `yaml:"ssh_key,omitempty" json:"sshKey,omitempty" norman:"type=password"`
	// SSH Private Key Path
	SSHKeyPath string `yaml:"ssh_key_path,omitempty" json:"sshKeyPath,omitempty"`
	// SSH Certificate
	SSHCert string `yaml:"ssh_cert,omitempty" json:"sshCert,omitempty"`
	// SSH Certificate Path
	SSHCertPath string `yaml:"ssh_cert_path,omitempty" json:"sshCertPath,omitempty"`
	// Node Labels
	Labels map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	// Node Taints
	Taints []RKETaint `yaml:"taints,omitempty" json:"taints,omitempty"`
}

type K8sVersionInfo struct {
	MinRKEVersion       string `yaml:"min_rke_version,omitempty" json:"minRKEVersion,omitempty"`
	MaxRKEVersion       string `yaml:"max_rke_version,omitempty" json:"maxRKEVersion,omitempty"`
	DeprecateRKEVersion string `yaml:"deprecate_rke_version,omitempty" json:"deprecateRKEVersion,omitempty"`

	MinRancherVersion       string `yaml:"min_rancher_version,omitempty" json:"minRancherVersion,omitempty"`
	MaxRancherVersion       string `yaml:"max_rancher_version,omitempty" json:"maxRancherVersion,omitempty"`
	DeprecateRancherVersion string `yaml:"deprecate_rancher_version,omitempty" json:"deprecateRancherVersion,omitempty"`
}

type RKEConfigServices struct {
	// Etcd Service
	Etcd ETCDService `yaml:"etcd,omitempty" json:"etcd,omitempty"`
	// KubeAPI Service
	KubeAPI KubeAPIService `yaml:"kube-api,omitempty" json:"kubeApi,omitempty"`
	// KubeController Service
	KubeController KubeControllerService `yaml:"kube-controller,omitempty" json:"kubeController,omitempty"`
	// Scheduler Service
	Scheduler SchedulerService `yaml:"scheduler,omitempty" json:"scheduler,omitempty"`
	// Kubelet Service
	Kubelet KubeletService `yaml:"kubelet,omitempty" json:"kubelet,omitempty"`
	// KubeProxy Service
	Kubeproxy KubeproxyService `yaml:"kubeproxy,omitempty" json:"kubeproxy,omitempty"`
}

type ETCDService struct {
	// Base service properties
	BaseService `yaml:",inline" json:",inline"`
	// List of etcd urls
	ExternalURLs []string `yaml:"external_urls,omitempty" json:"externalUrls,omitempty"`
	// External CA certificate
	CACert string `yaml:"ca_cert,omitempty" json:"caCert,omitempty"`
	// External Client certificate
	Cert string `yaml:"cert,omitempty" json:"cert,omitempty"`
	// External Client key
	Key string `yaml:"key,omitempty" json:"key,omitempty"`
	// External etcd prefix
	Path string `yaml:"path,omitempty" json:"path,omitempty"`
	// UID to run etcd container as
	UID int `yaml:"uid,omitempty" json:"uid,omitempty"`
	// GID to run etcd container as
	GID int `yaml:"gid,omitempty" json:"gid,omitempty"`

	// Etcd Recurring snapshot Service, used by rke only
	Snapshot *bool `yaml:"snapshot,omitempty" json:"snapshot,omitempty" norman:"default=false"`
	// Etcd snapshot Retention period
	Retention string `yaml:"retention,omitempty" json:"retention,omitempty" norman:"default=72h"`
	// Etcd snapshot Creation period
	Creation string `yaml:"creation,omitempty" json:"creation,omitempty" norman:"default=12h"`
	// Backup backend for etcd snapshots
	BackupConfig *BackupConfig `yaml:"backup_config,omitempty" json:"backupConfig,omitempty"`
}

type KubeAPIService struct {
	// Base service properties
	BaseService `yaml:",inline" json:",inline"`
	// Virtual IP range that will be used by Kubernetes services
	ServiceClusterIPRange string `yaml:"service_cluster_ip_range,omitempty" json:"serviceClusterIpRange,omitempty"`
	// Port range for services defined with NodePort type
	ServiceNodePortRange string `yaml:"service_node_port_range,omitempty" json:"serviceNodePortRange,omitempty" norman:"default=30000-32767"`
	// Enabled/Disable PodSecurityPolicy
	PodSecurityPolicy bool `yaml:"pod_security_policy,omitempty" json:"podSecurityPolicy,omitempty"`
	// Enable/Disable AlwaysPullImages admissions plugin
	AlwaysPullImages bool `yaml:"always_pull_images,omitempty" json:"alwaysPullImages,omitempty"`
	// Secrets encryption provider config
	SecretsEncryptionConfig *SecretsEncryptionConfig `yaml:"secrets_encryption_config,omitempty" json:"secretsEncryptionConfig,omitempty"`
	// Audit Log Configuration
	AuditLog *AuditLog `yaml:"audit_log,omitempty" json:"auditLog,omitempty"`
	// AdmissionConfiguration
	AdmissionConfiguration *apiserverv1alpha1.AdmissionConfiguration `yaml:"admission_configuration,omitempty" json:"admissionConfiguration,omitempty" norman:"type=map[json]"`
	// Event Rate Limit configuration
	EventRateLimit *EventRateLimit `yaml:"event_rate_limit,omitempty" json:"eventRateLimit,omitempty"`
}

type EventRateLimit struct {
	Enabled       bool           `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Configuration *Configuration `yaml:"configuration,omitempty" json:"configuration,omitempty" norman:"type=map[json]"`
}

type AuditLog struct {
	Enabled       bool            `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Configuration *AuditLogConfig `yaml:"configuration,omitempty" json:"configuration,omitempty"`
}

type AuditLogConfig struct {
	MaxAge    int             `yaml:"max_age,omitempty" json:"maxAge,omitempty"`
	MaxBackup int             `yaml:"max_backup,omitempty" json:"maxBackup,omitempty"`
	MaxSize   int             `yaml:"max_size,omitempty" json:"maxSize,omitempty"`
	Path      string          `yaml:"path,omitempty" json:"path,omitempty"`
	Format    string          `yaml:"format,omitempty" json:"format,omitempty"`
	Policy    *auditv1.Policy `yaml:"policy,omitempty" json:"policy,omitempty" norman:"type=map[json]"`
}

type KubeControllerService struct {
	// Base service properties
	BaseService `yaml:",inline" json:",inline"`
	// CIDR Range for Pods in cluster
	ClusterCIDR string `yaml:"cluster_cidr,omitempty" json:"clusterCidr,omitempty"`
	// Virtual IP range that will be used by Kubernetes services
	ServiceClusterIPRange string `yaml:"service_cluster_ip_range,omitempty" json:"serviceClusterIpRange,omitempty"`
}

type KubeletService struct {
	// Base service properties
	BaseService `yaml:",inline" json:",inline"`
	// Domain of the cluster (default: "cluster.local")
	ClusterDomain string `yaml:"cluster_domain,omitempty" json:"clusterDomain,omitempty"`
	// The image whose network/ipc namespaces containers in each pod will use
	InfraContainerImage string `yaml:"infra_container_image,omitempty" json:"infraContainerImage,omitempty"`
	// Cluster DNS service ip
	ClusterDNSServer string `yaml:"cluster_dns_server,omitempty" json:"clusterDnsServer,omitempty"`
	// Fail if swap is enabled
	FailSwapOn bool `yaml:"fail_swap_on,omitempty" json:"failSwapOn,omitempty"`
	// Generate per node kubelet serving certificates created using kube-ca
	GenerateServingCertificate bool `yaml:"generate_serving_certificate,omitempty" json:"generateServingCertificate,omitempty"`
}

type KubeproxyService struct {
	// Base service properties
	BaseService `yaml:",inline" json:",inline"`
}

type SchedulerService struct {
	// Base service properties
	BaseService `yaml:",inline" json:",inline"`
}

type BaseService struct {
	// Docker image of the service
	Image string `yaml:"image,omitempty" json:"image,omitempty"`
	// Extra arguments that are added to the services
	ExtraArgs map[string]string `yaml:"extra_args,omitempty" json:"extraArgs,omitempty"`
	// Extra binds added to the nodes
	ExtraBinds []string `yaml:"extra_binds,omitempty" json:"extraBinds,omitempty"`
	// this is to provide extra env variable to the docker container running kubernetes service
	ExtraEnv []string `yaml:"extra_env,omitempty" json:"extraEnv,omitempty"`

	// Windows nodes only of the same as the above
	// Extra arguments that are added to the services
	WindowsExtraArgs map[string]string `yaml:"win_extra_args,omitempty" json:"winExtraArgs,omitempty"`
	// Extra binds added to the nodes
	WindowsExtraBinds []string `yaml:"win_extra_binds,omitempty" json:"winExtraBinds,omitempty"`
	// this is to provide extra env variable to the docker container running kubernetes service
	WindowsExtraEnv []string `yaml:"win_extra_env,omitempty" json:"winExtraEnv,omitempty"`
}

type NetworkConfig struct {
	// Network Plugin That will be used in kubernetes cluster
	Plugin string `yaml:"plugin,omitempty" json:"plugin,omitempty" norman:"default=canal"`
	// Plugin options to configure network properties
	Options map[string]string `yaml:"options,omitempty" json:"options,omitempty"`
	// Set MTU for CNI provider
	MTU int `yaml:"mtu,omitempty" json:"mtu,omitempty"`
	// CalicoNetworkProvider
	CalicoNetworkProvider *CalicoNetworkProvider `yaml:"calico_network_provider,omitempty" json:"calicoNetworkProvider,omitempty"`
	// CanalNetworkProvider
	CanalNetworkProvider *CanalNetworkProvider `yaml:"canal_network_provider,omitempty" json:"canalNetworkProvider,omitempty"`
	// FlannelNetworkProvider
	FlannelNetworkProvider *FlannelNetworkProvider `yaml:"flannel_network_provider,omitempty" json:"flannelNetworkProvider,omitempty"`
	// WeaveNetworkProvider
	WeaveNetworkProvider *WeaveNetworkProvider `yaml:"weave_network_provider,omitempty" json:"weaveNetworkProvider,omitempty"`
	// AciNetworkProvider
	AciNetworkProvider *AciNetworkProvider `yaml:"aci_network_provider,omitempty" json:"aciNetworkProvider,omitempty"`
	// NodeSelector key pair
	NodeSelector map[string]string `yaml:"node_selector,omitempty" json:"nodeSelector,omitempty"`
	// Network plugin daemonset upgrade strategy
	UpdateStrategy *DaemonSetUpdateStrategy `yaml:"update_strategy,omitempty" json:"updateStrategy,omitempty"`
	// Tolerations for Deployments
	Tolerations []v1.Toleration `yaml:"tolerations,omitempty" json:"tolerations,omitempty"`
}

type AuthWebhookConfig struct {
	// ConfigFile is a multiline string that represent a custom webhook config file
	ConfigFile string `yaml:"config_file,omitempty" json:"configFile,omitempty"`
	// CacheTimeout controls how long to cache authentication decisions
	CacheTimeout string `yaml:"cache_timeout,omitempty" json:"cacheTimeout,omitempty"`
}

type AuthnConfig struct {
	// Authentication strategy that will be used in kubernetes cluster
	Strategy string `yaml:"strategy,omitempty" json:"strategy,omitempty" norman:"default=x509"`
	// List of additional hostnames and IPs to include in the api server PKI cert
	SANs []string `yaml:"sans,omitempty" json:"sans,omitempty"`
	// Webhook configuration options
	Webhook *AuthWebhookConfig `yaml:"webhook,omitempty" json:"webhook,omitempty"`
}

type AuthzConfig struct {
	// Authorization mode used by kubernetes
	Mode string `yaml:"mode,omitempty" json:"mode,omitempty"`
	// Authorization mode options
	Options map[string]string `yaml:"options,omitempty" json:"options,omitempty"`
}

type IngressConfig struct {
	// Ingress controller type used by kubernetes
	Provider string `yaml:"provider,omitempty" json:"provider,omitempty" norman:"default=nginx"`
	// These options are NOT for configuring Ingress's addon template.
	// They are used for its ConfigMap options specifically.
	Options map[string]string `yaml:"options,omitempty" json:"options,omitempty"`
	// NodeSelector key pair
	NodeSelector map[string]string `yaml:"node_selector,omitempty" json:"nodeSelector,omitempty"`
	// Ingress controller extra arguments
	ExtraArgs map[string]string `yaml:"extra_args,omitempty" json:"extraArgs,omitempty"`
	// DNS Policy
	DNSPolicy string `yaml:"dns_policy,omitempty" json:"dnsPolicy,omitempty"`
	// Extra Env vars
	ExtraEnvs []ExtraEnv `yaml:"extra_envs,omitempty" json:"extraEnvs,omitempty" norman:"type=array[json]"`
	// Extra volumes
	ExtraVolumes []ExtraVolume `yaml:"extra_volumes,omitempty" json:"extraVolumes,omitempty" norman:"type=array[json]"`
	// Extra volume mounts
	ExtraVolumeMounts []ExtraVolumeMount `yaml:"extra_volume_mounts,omitempty" json:"extraVolumeMounts,omitempty" norman:"type=array[json]"`
	// nginx daemonset upgrade strategy
	UpdateStrategy *DaemonSetUpdateStrategy `yaml:"update_strategy,omitempty" json:"updateStrategy,omitempty"`
	// Http port for ingress controller daemonset
	HTTPPort int `yaml:"http_port,omitempty" json:"httpPort,omitempty" jsonschema:"oneof_type=string;number"`
	// Https port for ingress controller daemonset
	HTTPSPort int `yaml:"https_port,omitempty" json:"httpsPort,omitempty" jsonschema:"oneof_type=string;number"`
	// NetworkMode selector for ingress controller pods. Default is HostNetwork
	NetworkMode string `yaml:"network_mode,omitempty" json:"networkMode,omitempty"`
	// Tolerations for Deployments
	Tolerations []v1.Toleration `yaml:"tolerations,omitempty" json:"tolerations,omitempty"`
	// Enable or disable nginx default-http-backend
	DefaultBackend *bool `yaml:"default_backend,omitempty" json:"defaultBackend,omitempty" norman:"default=true"`
	// Priority class name for Nginx-Ingress's "default-http-backend" deployment
	DefaultHTTPBackendPriorityClassName string `yaml:"default_http_backend_priority_class_name,omitempty" json:"defaultHttpBackendPriorityClassName,omitempty"`
	// Priority class name for Nginx-Ingress's "nginx-ingress-controller" daemonset
	NginxIngressControllerPriorityClassName string `yaml:"nginx_ingress_controller_priority_class_name,omitempty" json:"nginxIngressControllerPriorityClassName,omitempty"`
}

type ExtraEnv struct {
	v1.EnvVar
}

type ExtraVolume struct {
	v1.Volume
}

type ExtraVolumeMount struct {
	v1.VolumeMount
}

type RKEPlan struct {
	// List of node Plans
	Nodes []RKEConfigNodePlan `json:"nodes,omitempty"`
}

type RKEConfigNodePlan struct {
	// Node address
	Address string `json:"address,omitempty"`
	// map of named processes that should run on the node
	Processes map[string]Process `json:"processes,omitempty"`
	// List of portchecks that should be open on the node
	PortChecks []PortCheck `json:"portChecks,omitempty"`
	// List of files to deploy on the node
	Files []File `json:"files,omitempty"`
	// Node Annotations
	Annotations map[string]string `json:"annotations,omitempty"`
	// Node Labels
	Labels map[string]string `json:"labels,omitempty"`
	// Node Taints
	Taints []RKETaint `json:"taints,omitempty"`
}

type Process struct {
	// Process name, this should be the container name
	Name string `json:"name,omitempty"`
	// Process Entrypoint command
	Command []string `json:"command,omitempty"`
	// Process args
	Args []string `json:"args,omitempty"`
	// Environment variables list
	Env []string `json:"env,omitempty"`
	// Process docker image
	Image string `json:"image,omitempty"`
	//AuthConfig for image private registry
	ImageRegistryAuthConfig string `json:"imageRegistryAuthConfig,omitempty"`
	// Process docker image VolumesFrom
	VolumesFrom []string `json:"volumesFrom,omitempty"`
	// Process docker container bind mounts
	Binds []string `json:"binds,omitempty"`
	// Process docker container netwotk mode
	NetworkMode string `json:"networkMode,omitempty"`
	// Process container restart policy
	RestartPolicy string `json:"restartPolicy,omitempty"`
	// Process container pid mode
	PidMode string `json:"pidMode,omitempty"`
	// Run process in privileged container
	Privileged bool `json:"privileged,omitempty"`
	// Process healthcheck
	HealthCheck HealthCheck `json:"healthCheck,omitempty"`
	// Process docker container Labels
	Labels map[string]string `json:"labels,omitempty"`
	// Process docker publish container's port to host
	Publish []string `json:"publish,omitempty"`
	// docker will run the container with this user
	User string `json:"user,omitempty"`
}

type HealthCheck struct {
	// Healthcheck URL
	URL string `json:"url,omitempty"`
}

type PortCheck struct {
	// Portcheck address to check.
	Address string `json:"address,omitempty"`
	// Port number
	Port int `json:"port,omitempty"`
	// Port Protocol
	Protocol string `json:"protocol,omitempty"`
}

type CloudProvider struct {
	// Name of the Cloud Provider
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
	// AWSCloudProvider
	AWSCloudProvider *AWSCloudProvider `yaml:"awsCloudProvider,omitempty" json:"awsCloudProvider,omitempty"`
	// AzureCloudProvider
	AzureCloudProvider *AzureCloudProvider `yaml:"azureCloudProvider,omitempty" json:"azureCloudProvider,omitempty"`
	// OpenstackCloudProvider
	OpenstackCloudProvider *OpenstackCloudProvider `yaml:"openstackCloudProvider,omitempty" json:"openstackCloudProvider,omitempty"`
	// VsphereCloudProvider
	VsphereCloudProvider *VsphereCloudProvider `yaml:"vsphereCloudProvider,omitempty" json:"vsphereCloudProvider,omitempty"`
	// CustomCloudProvider is a multiline string that represent a custom cloud config file
	CustomCloudProvider string `yaml:"customCloudProvider,omitempty" json:"customCloudProvider,omitempty"`
}

type CalicoNetworkProvider struct {
	// Cloud provider type used with calico
	CloudProvider string `json:"cloudProvider"`
}

type FlannelNetworkProvider struct {
	// Alternate cloud interface for flannel
	Iface string `json:"iface"`
}

type CanalNetworkProvider struct {
	FlannelNetworkProvider `yaml:",inline" json:",inline"`
}

type WeaveNetworkProvider struct {
	Password string `yaml:"password,omitempty" json:"password,omitempty" norman:"type=password"`
}

type AciNetworkProvider struct {
	SystemIdentifier         string   `yaml:"system_id,omitempty,omitempty" json:"systemId,omitempty"`
	ApicHosts                []string `yaml:"apic_hosts,omitempty" json:"apicHosts,omitempty"`
	Token                    string   `yaml:"token,omitempty" json:"token,omitempty"`
	ApicUserName             string   `yaml:"apic_user_name,omitempty" json:"apicUserName,omitempty"`
	ApicUserKey              string   `yaml:"apic_user_key,omitempty" json:"apicUserKey,omitempty"`
	ApicUserCrt              string   `yaml:"apic_user_crt,omitempty" json:"apicUserCrt,omitempty"`
	ApicRefreshTime          string   `yaml:"apic_refresh_time,omitempty" json:"apicRefreshTime,omitempty" norman:"default=1200"`
	VmmDomain                string   `yaml:"vmm_domain,omitempty" json:"vmmDomain,omitempty"`
	VmmController            string   `yaml:"vmm_controller,omitempty" json:"vmmController,omitempty"`
	EncapType                string   `yaml:"encap_type,omitempty" json:"encapType,omitempty"`
	NodeSubnet               string   `yaml:"node_subnet,omitempty" json:"nodeSubnet,omitempty"`
	McastRangeStart          string   `yaml:"mcast_range_start,omitempty" json:"mcastRangeStart,omitempty"`
	McastRangeEnd            string   `yaml:"mcast_range_end,omitempty" json:"mcastRangeEnd,omitempty"`
	AEP                      string   `yaml:"aep,omitempty" json:"aep,omitempty"`
	VRFName                  string   `yaml:"vrf_name,omitempty" json:"vrfName,omitempty"`
	VRFTenant                string   `yaml:"vrf_tenant,omitempty" json:"vrfTenant,omitempty"`
	L3Out                    string   `yaml:"l3out,omitempty" json:"l3out,omitempty"`
	L3OutExternalNetworks    []string `yaml:"l3out_external_networks,omitempty" json:"l3outExternalNetworks,omitempty"`
	DynamicExternalSubnet    string   `yaml:"extern_dynamic,omitempty" json:"externDynamic,omitempty"`
	StaticExternalSubnet     string   `yaml:"extern_static,omitempty" json:"externStatic,omitempty"`
	ServiceGraphSubnet       string   `yaml:"node_svc_subnet,omitempty" json:"nodeSvcSubnet,omitempty"`
	KubeAPIVlan              string   `yaml:"kube_api_vlan,omitempty" json:"kubeApiVlan,omitempty"`
	ServiceVlan              string   `yaml:"service_vlan,omitempty" json:"serviceVlan,omitempty"`
	InfraVlan                string   `yaml:"infra_vlan,omitempty" json:"infraVlan,omitempty"`
	Tenant                   string   `yaml:"tenant,omitempty" json:"tenant,omitempty"`
	OVSMemoryLimit           string   `yaml:"ovs_memory_limit,omitempty" json:"ovsMemoryLimit,omitempty"`
	ImagePullPolicy          string   `yaml:"image_pull_policy,omitempty" json:"imagePullPolicy,omitempty"`
	ImagePullSecret          string   `yaml:"image_pull_secret,omitempty" json:"imagePullSecret,omitempty"`
	ServiceMonitorInterval   string   `yaml:"service_monitor_interval,omitempty" json:"serviceMonitorInterval,omitempty"`
	PBRTrackingNonSnat       string   `yaml:"pbr_tracking_non_snat,omitempty" json:"pbrTrackingNonSnat,omitempty"`
	InstallIstio             string   `yaml:"install_istio,omitempty" json:"installIstio,omitempty"`
	IstioProfile             string   `yaml:"istio_profile,omitempty" json:"istioProfile,omitempty"`
	DropLogEnable            string   `yaml:"drop_log_enable,omitempty" json:"dropLogEnable,omitempty"`
	ControllerLogLevel       string   `yaml:"controller_log_level,omitempty" json:"controllerLogLevel,omitempty"`
	HostAgentLogLevel        string   `yaml:"host_agent_log_level,omitempty" json:"hostAgentLogLevel,omitempty"`
	OpflexAgentLogLevel      string   `yaml:"opflex_log_level,omitempty" json:"opflexLogLevel,omitempty"`
	UseAciCniPriorityClass   string   `yaml:"use_aci_cni_priority_class,omitempty" json:"useAciCniPriorityClass,omitempty"`
	NoPriorityClass          string   `yaml:"no_priority_class,omitempty" json:"noPriorityClass,omitempty"`
	MaxNodesSvcGraph         string   `yaml:"max_nodes_svc_graph,omitempty" json:"maxNodesSvcGraph,omitempty"`
	SnatContractScope        string   `yaml:"snat_contract_scope,omitempty" json:"snatContractScope,omitempty"`
	PodSubnetChunkSize       string   `yaml:"pod_subnet_chunk_size,omitempty" json:"podSubnetChunkSize,omitempty"`
	EnableEndpointSlice      string   `yaml:"enable_endpoint_slice,omitempty" json:"enableEndpointSlice,omitempty"`
	SnatNamespace            string   `yaml:"snat_namespace,omitempty" json:"snatNamespace,omitempty"`
	EpRegistry               string   `yaml:"ep_registry,omitempty" json:"epRegistry,omitempty"`
	OpflexMode               string   `yaml:"opflex_mode,omitempty" json:"opflexMode,omitempty"`
	SnatPortRangeStart       string   `yaml:"snat_port_range_start,omitempty" json:"snatPortRangeStart,omitempty" jsonschema:"oneof_type=string;number"`
	SnatPortRangeEnd         string   `yaml:"snat_port_range_end,omitempty" json:"snatPortRangeEnd,omitempty" jsonschema:"oneof_type=string;number"`
	SnatPortsPerNode         string   `yaml:"snat_ports_per_node,omitempty" json:"snatPortsPerNode,omitempty" jsonschema:"oneof_type=string;number"`
	OpflexClientSSL          string   `yaml:"opflex_client_ssl,omitempty" json:"opflexClientSsl,omitempty"`
	UsePrivilegedContainer   string   `yaml:"use_privileged_container,omitempty" json:"usePrivilegedContainer,omitempty"`
	UseHostNetnsVolume       string   `yaml:"use_host_netns_volume,omitempty" json:"useHostNetnsVolume,omitempty"`
	UseOpflexServerVolume    string   `yaml:"use_opflex_server_volume,omitempty" json:"useOpflexServerVolume,omitempty"`
	SubnetDomainName         string   `yaml:"subnet_domain_name,omitempty" json:"subnetDomainName,omitempty"`
	KafkaBrokers             []string `yaml:"kafka_brokers,omitempty" json:"kafkaBrokers,omitempty"`
	KafkaClientCrt           string   `yaml:"kafka_client_crt,omitempty" json:"kafkaClientCrt,omitempty"`
	KafkaClientKey           string   `yaml:"kafka_client_key,omitempty" json:"kafkaClientKey,omitempty"`
	CApic                    string   `yaml:"capic,omitempty" json:"capic,omitempty"`
	UseAciAnywhereCRD        string   `yaml:"use_aci_anywhere_crd,omitempty" json:"useAciAnywhereCrd,omitempty"`
	OverlayVRFName           string   `yaml:"overlay_vrf_name,omitempty" json:"overlayVrfName,omitempty"`
	GbpPodSubnet             string   `yaml:"gbp_pod_subnet,omitempty" json:"gbpPodSubnet,omitempty"`
	RunGbpContainer          string   `yaml:"run_gbp_container,omitempty" json:"runGbpContainer,omitempty"`
	RunOpflexServerContainer string   `yaml:"run_opflex_server_container,omitempty" json:"runOpflexServerContainer,omitempty"`
	OpflexServerPort         string   `yaml:"opflex_server_port,omitempty" json:"opflexServerPort,omitempty" jsonschema:"oneof_type=string;number"`
}

type KubernetesServicesOptions struct {
	// Additional options passed to Etcd
	Etcd map[string]string `json:"etcd"`
	// Additional options passed to KubeAPI
	KubeAPI map[string]string `json:"kubeapi"`
	// Additional options passed to Kubelet
	Kubelet map[string]string `json:"kubelet"`
	// Additional options passed to Kubeproxy
	Kubeproxy map[string]string `json:"kubeproxy"`
	// Additional options passed to KubeController
	KubeController map[string]string `json:"kubeController"`
	// Additional options passed to Scheduler
	Scheduler map[string]string `json:"scheduler"`
}

// VsphereCloudProvider options
type VsphereCloudProvider struct {
	Global        GlobalVsphereOpts              `json:"global,omitempty" yaml:"global,omitempty" ini:"Global,omitempty"`
	VirtualCenter map[string]VirtualCenterConfig `json:"virtualCenter,omitempty" yaml:"virtual_center,omitempty" ini:"VirtualCenter,omitempty"`
	Network       NetworkVshpereOpts             `json:"network,omitempty" yaml:"network,omitempty" ini:"Network,omitempty"`
	Disk          DiskVsphereOpts                `json:"disk,omitempty" yaml:"disk,omitempty" ini:"Disk,omitempty"`
	Workspace     WorkspaceVsphereOpts           `json:"workspace,omitempty" yaml:"workspace,omitempty" ini:"Workspace,omitempty"`
}

type GlobalVsphereOpts struct {
	User              string `json:"user,omitempty" yaml:"user,omitempty" ini:"user,omitempty"`
	Password          string `json:"password,omitempty" yaml:"password,omitempty" ini:"password,omitempty" norman:"type=password"`
	VCenterIP         string `json:"server,omitempty" yaml:"server,omitempty" ini:"server,omitempty"`
	VCenterPort       string `json:"port,omitempty" yaml:"port,omitempty" ini:"port,omitempty" jsonschema:"oneof_type=string;number"`
	InsecureFlag      bool   `json:"insecure-flag,omitempty" yaml:"insecure-flag,omitempty" ini:"insecure-flag,omitempty"`
	Datacenter        string `json:"datacenter,omitempty" yaml:"datacenter,omitempty" ini:"datacenter,omitempty"`
	Datacenters       string `json:"datacenters,omitempty" yaml:"datacenters,omitempty" ini:"datacenters,omitempty"`
	DefaultDatastore  string `json:"datastore,omitempty" yaml:"datastore,omitempty" ini:"datastore,omitempty"`
	WorkingDir        string `json:"working-dir,omitempty" yaml:"working-dir,omitempty" ini:"working-dir,omitempty"`
	RoundTripperCount int    `json:"soap-roundtrip-count,omitempty" yaml:"soap-roundtrip-count,omitempty" ini:"soap-roundtrip-count,omitempty"`
	VMUUID            string `json:"vm-uuid,omitempty" yaml:"vm-uuid,omitempty" ini:"vm-uuid,omitempty"`
	VMName            string `json:"vm-name,omitempty" yaml:"vm-name,omitempty" ini:"vm-name,omitempty"`
}

type VirtualCenterConfig struct {
	User              string `json:"user,omitempty" yaml:"user,omitempty" ini:"user,omitempty"`
	Password          string `json:"password,omitempty" yaml:"password,omitempty" ini:"password,omitempty" norman:"type=password"`
	VCenterPort       string `json:"port,omitempty" yaml:"port,omitempty" ini:"port,omitempty" jsonschema:"oneof_type=string;number"`
	Datacenters       string `json:"datacenters,omitempty" yaml:"datacenters,omitempty" ini:"datacenters,omitempty"`
	RoundTripperCount int    `json:"soap-roundtrip-count,omitempty" yaml:"soap-roundtrip-count,omitempty" ini:"soap-roundtrip-count,omitempty"`
}

type NetworkVshpereOpts struct {
	PublicNetwork string `json:"public-network,omitempty" yaml:"public-network,omitempty" ini:"public-network,omitempty"`
}

type DiskVsphereOpts struct {
	SCSIControllerType string `json:"scsicontrollertype,omitempty" yaml:"scsicontrollertype,omitempty" ini:"scsicontrollertype,omitempty"`
}

type WorkspaceVsphereOpts struct {
	VCenterIP        string `json:"server,omitempty" yaml:"server,omitempty" ini:"server,omitempty"`
	Datacenter       string `json:"datacenter,omitempty" yaml:"datacenter,omitempty" ini:"datacenter,omitempty"`
	Folder           string `json:"folder,omitempty" yaml:"folder,omitempty" ini:"folder,omitempty"`
	DefaultDatastore string `json:"default-datastore,omitempty" yaml:"default-datastore,omitempty" ini:"default-datastore,omitempty"`
	ResourcePoolPath string `json:"resourcepool-path,omitempty" yaml:"resourcepool-path,omitempty" ini:"resourcepool-path,omitempty"`
}

// OpenstackCloudProvider options
type OpenstackCloudProvider struct {
	Global       GlobalOpenstackOpts       `json:"global" yaml:"global" ini:"Global,omitempty"`
	LoadBalancer LoadBalancerOpenstackOpts `json:"loadBalancer" yaml:"load_balancer" ini:"LoadBalancer,omitempty"`
	BlockStorage BlockStorageOpenstackOpts `json:"blockStorage" yaml:"block_storage" ini:"BlockStorage,omitempty"`
	Route        RouteOpenstackOpts        `json:"route" yaml:"route" ini:"Route,omitempty"`
	Metadata     MetadataOpenstackOpts     `json:"metadata" yaml:"metadata" ini:"Metadata,omitempty"`
}

type GlobalOpenstackOpts struct {
	AuthURL    string `json:"auth-url" yaml:"auth-url" ini:"auth-url,omitempty"`
	Username   string `json:"username" yaml:"username" ini:"username,omitempty"`
	UserID     string `json:"user-id" yaml:"user-id" ini:"user-id,omitempty"`
	Password   string `json:"password" yaml:"password" ini:"password,omitempty" norman:"type=password"`
	TenantID   string `json:"tenant-id" yaml:"tenant-id" ini:"tenant-id,omitempty"`
	TenantName string `json:"tenant-name" yaml:"tenant-name" ini:"tenant-name,omitempty"`
	TrustID    string `json:"trust-id" yaml:"trust-id" ini:"trust-id,omitempty"`
	DomainID   string `json:"domain-id" yaml:"domain-id" ini:"domain-id,omitempty"`
	DomainName string `json:"domain-name" yaml:"domain-name" ini:"domain-name,omitempty"`
	Region     string `json:"region" yaml:"region" ini:"region,omitempty"`
	CAFile     string `json:"ca-file" yaml:"ca-file" ini:"ca-file,omitempty"`
}

type LoadBalancerOpenstackOpts struct {
	LBVersion            string `json:"lb-version" yaml:"lb-version" ini:"lb-version,omitempty"`                            // overrides autodetection. Only support v2.
	UseOctavia           bool   `json:"use-octavia" yaml:"use-octavia" ini:"use-octavia,omitempty"`                         // uses Octavia V2 service catalog endpoint
	SubnetID             string `json:"subnet-id" yaml:"subnet-id" ini:"subnet-id,omitempty"`                               // overrides autodetection.
	FloatingNetworkID    string `json:"floating-network-id" yaml:"floating-network-id" ini:"floating-network-id,omitempty"` // If specified, will create floating ip for loadbalancer, or do not create floating ip.
	LBMethod             string `json:"lb-method" yaml:"lb-method" ini:"lb-method,omitempty"`                               // default to ROUND_ROBIN.
	LBProvider           string `json:"lb-provider" yaml:"lb-provider" ini:"lb-provider,omitempty"`
	CreateMonitor        bool   `json:"create-monitor" yaml:"create-monitor" ini:"create-monitor,omitempty"`
	MonitorDelay         string `json:"monitor-delay" yaml:"monitor-delay" ini:"monitor-delay,omitempty"`
	MonitorTimeout       string `json:"monitor-timeout" yaml:"monitor-timeout" ini:"monitor-timeout,omitempty"`
	MonitorMaxRetries    int    `json:"monitor-max-retries" yaml:"monitor-max-retries" ini:"monitor-max-retries,omitempty"`
	ManageSecurityGroups bool   `json:"manage-security-groups" yaml:"manage-security-groups" ini:"manage-security-groups,omitempty"`
}

type BlockStorageOpenstackOpts struct {
	BSVersion       string `json:"bs-version" yaml:"bs-version" ini:"bs-version,omitempty"`                      // overrides autodetection. v1 or v2. Defaults to auto
	TrustDevicePath bool   `json:"trust-device-path" yaml:"trust-device-path" ini:"trust-device-path,omitempty"` // See Issue #33128
	IgnoreVolumeAZ  bool   `json:"ignore-volume-az" yaml:"ignore-volume-az" ini:"ignore-volume-az,omitempty"`
}

type RouteOpenstackOpts struct {
	RouterID string `json:"router-id" yaml:"router-id" ini:"router-id,omitempty"` // required
}

type MetadataOpenstackOpts struct {
	SearchOrder    string `json:"search-order" yaml:"search-order" ini:"search-order,omitempty"`
	RequestTimeout int    `json:"request-timeout" yaml:"request-timeout" ini:"request-timeout,omitempty"`
}

// AzureCloudProvider options
type AzureCloudProvider struct {
	// The cloud environment identifier. Takes values from https://github.com/Azure/go-autorest/blob/ec5f4903f77ed9927ac95b19ab8e44ada64c1356/autorest/azure/environments.go#L13
	Cloud string `json:"cloud" yaml:"cloud"`
	// The AAD Tenant ID for the Subscription that the cluster is deployed in
	TenantID string `json:"tenantId" yaml:"tenantId"`
	// The ID of the Azure Subscription that the cluster is deployed in
	SubscriptionID string `json:"subscriptionId" yaml:"subscriptionId"`
	// The name of the resource group that the cluster is deployed in
	ResourceGroup string `json:"resourceGroup" yaml:"resourceGroup"`
	// The location of the resource group that the cluster is deployed in
	Location string `json:"location" yaml:"location"`
	// The name of the VNet that the cluster is deployed in
	VnetName string `json:"vnetName" yaml:"vnetName"`
	// The name of the resource group that the Vnet is deployed in
	VnetResourceGroup string `json:"vnetResourceGroup" yaml:"vnetResourceGroup"`
	// The name of the subnet that the cluster is deployed in
	SubnetName string `json:"subnetName" yaml:"subnetName"`
	// The name of the security group attached to the cluster's subnet
	SecurityGroupName string `json:"securityGroupName" yaml:"securityGroupName"`
	// (Optional in 1.6) The name of the route table attached to the subnet that the cluster is deployed in
	RouteTableName string `json:"routeTableName" yaml:"routeTableName"`
	// (Optional) The name of the availability set that should be used as the load balancer backend
	// If this is set, the Azure cloudprovider will only add nodes from that availability set to the load
	// balancer backend pool. If this is not set, and multiple agent pools (availability sets) are used, then
	// the cloudprovider will try to add all nodes to a single backend pool which is forbidden.
	// In other words, if you use multiple agent pools (availability sets), you MUST set this field.
	PrimaryAvailabilitySetName string `json:"primaryAvailabilitySetName" yaml:"primaryAvailabilitySetName"`
	// The type of azure nodes. Candidate valudes are: vmss and standard.
	// If not set, it will be default to standard.
	VMType string `json:"vmType" yaml:"vmType"`
	// The name of the scale set that should be used as the load balancer backend.
	// If this is set, the Azure cloudprovider will only add nodes from that scale set to the load
	// balancer backend pool. If this is not set, and multiple agent pools (scale sets) are used, then
	// the cloudprovider will try to add all nodes to a single backend pool which is forbidden.
	// In other words, if you use multiple agent pools (scale sets), you MUST set this field.
	PrimaryScaleSetName string `json:"primaryScaleSetName" yaml:"primaryScaleSetName"`
	// The ClientID for an AAD application with RBAC access to talk to Azure RM APIs
	// This's used for service principal authentication: https://github.com/Azure/aks-engine/blob/master/docs/topics/service-principals.md
	AADClientID string `json:"aadClientId" yaml:"aadClientId"`
	// The ClientSecret for an AAD application with RBAC access to talk to Azure RM APIs
	// This's used for service principal authentication: https://github.com/Azure/aks-engine/blob/master/docs/topics/service-principals.md
	AADClientSecret string `json:"aadClientSecret" yaml:"aadClientSecret" norman:"type=password"`
	// The path of a client certificate for an AAD application with RBAC access to talk to Azure RM APIs
	// This's used for client certificate authentication: https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-protocols-oauth-service-to-service
	AADClientCertPath string `json:"aadClientCertPath" yaml:"aadClientCertPath"`
	// The password of the client certificate for an AAD application with RBAC access to talk to Azure RM APIs
	// This's used for client certificate authentication: https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-protocols-oauth-service-to-service
	AADClientCertPassword string `json:"aadClientCertPassword" yaml:"aadClientCertPassword" norman:"type=password"`
	// Enable exponential backoff to manage resource request retries
	CloudProviderBackoff bool `json:"cloudProviderBackoff" yaml:"cloudProviderBackoff"`
	// Backoff retry limit
	CloudProviderBackoffRetries int `json:"cloudProviderBackoffRetries" yaml:"cloudProviderBackoffRetries"`
	// Backoff exponent
	CloudProviderBackoffExponent int `json:"cloudProviderBackoffExponent" yaml:"cloudProviderBackoffExponent"`
	// Backoff duration
	CloudProviderBackoffDuration int `json:"cloudProviderBackoffDuration" yaml:"cloudProviderBackoffDuration"`
	// Backoff jitter
	CloudProviderBackoffJitter int `json:"cloudProviderBackoffJitter" yaml:"cloudProviderBackoffJitter"`
	// Enable rate limiting
	CloudProviderRateLimit bool `json:"cloudProviderRateLimit" yaml:"cloudProviderRateLimit"`
	// Rate limit QPS
	CloudProviderRateLimitQPS int `json:"cloudProviderRateLimitQPS" yaml:"cloudProviderRateLimitQPS"`
	// Rate limit Bucket Size
	CloudProviderRateLimitBucket int `json:"cloudProviderRateLimitBucket" yaml:"cloudProviderRateLimitBucket"`
	// Use instance metadata service where possible
	UseInstanceMetadata bool `json:"useInstanceMetadata" yaml:"useInstanceMetadata"`
	// Use managed service identity for the virtual machine to access Azure ARM APIs
	// This's used for managed identity authentication: https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/overview
	// For user-assigned managed identity, need to set the below UserAssignedIdentityID
	UseManagedIdentityExtension bool `json:"useManagedIdentityExtension" yaml:"useManagedIdentityExtension"`
	// The Client ID of the user assigned MSI which is assigned to the underlying VMs
	// This's used for managed identity authentication: https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/overview
	UserAssignedIdentityID string `json:"userAssignedIdentityID,omitempty" yaml:"userAssignedIdentityID,omitempty"`
	// Maximum allowed LoadBalancer Rule Count is the limit enforced by Azure Load balancer, default(0) to 148
	MaximumLoadBalancerRuleCount int `json:"maximumLoadBalancerRuleCount" yaml:"maximumLoadBalancerRuleCount"`
	// Sku of Load Balancer and Public IP: `basic` or `standard`, default(blank) to `basic`
	LoadBalancerSku string `json:"loadBalancerSku,omitempty" yaml:"loadBalancerSku,omitempty"`
	// Excludes master nodes (labeled with `node-role.kubernetes.io/master`) from the backend pool of Azure standard loadbalancer, default(nil) to `true`
	// If want adding the master nodes to ALB, this should be set to `false` and remove the `node-role.kubernetes.io/master` label from master nodes
	ExcludeMasterFromStandardLB *bool `json:"excludeMasterFromStandardLB,omitempty" yaml:"excludeMasterFromStandardLB,omitempty"`
}

// AWSCloudProvider options
type AWSCloudProvider struct {
	Global          GlobalAwsOpts              `json:"global" yaml:"global" ini:"Global,omitempty"`
	ServiceOverride map[string]ServiceOverride `json:"serviceOverride,omitempty" yaml:"service_override,omitempty" ini:"ServiceOverride,omitempty"`
}

type ServiceOverride struct {
	Service       string `json:"service" yaml:"service" ini:"Service,omitempty"`
	Region        string `json:"region" yaml:"region" ini:"Region,omitempty"`
	URL           string `json:"url" yaml:"url" ini:"URL,omitempty"`
	SigningRegion string `json:"signing-region" yaml:"signing-region" ini:"SigningRegion,omitempty"`
	SigningMethod string `json:"signing-method" yaml:"signing-method" ini:"SigningMethod,omitempty"`
	SigningName   string `json:"signing-name" yaml:"signing-name" ini:"SigningName,omitempty"`
}

type GlobalAwsOpts struct {
	// TODO: Is there any use for this?  We can get it from the instance metadata service
	// Maybe if we're not running on AWS, e.g. bootstrap; for now it is not very useful
	Zone string `json:"zone" yaml:"zone" ini:"Zone,omitempty"`

	// The AWS VPC flag enables the possibility to run the master components
	// on a different aws account, on a different cloud provider or on-premises.
	// If the flag is set also the KubernetesClusterTag must be provided
	VPC string `json:"vpc" yaml:"vpc" ini:"VPC,omitempty"`
	// SubnetID enables using a specific subnet to use for ELB's
	SubnetID string `json:"subnet-id" yaml:"subnet-id" ini:"SubnetID,omitempty"`
	// RouteTableID enables using a specific RouteTable
	RouteTableID string `json:"routetable-id" yaml:"routetable-id" ini:"RouteTableID,omitempty"`

	// RoleARN is the IAM role to assume when interaction with AWS APIs.
	RoleARN string `json:"role-arn" yaml:"role-arn" ini:"RoleARN,omitempty"`

	// KubernetesClusterTag is the legacy cluster id we'll use to identify our cluster resources
	KubernetesClusterTag string `json:"kubernetes-cluster-tag" yaml:"kubernetes-cluster-tag" ini:"KubernetesClusterTag,omitempty"`
	// KubernetesClusterID is the cluster id we'll use to identify our cluster resources
	KubernetesClusterID string `json:"kubernetes-cluster-id" yaml:"kubernetes-cluster-id" ini:"KubernetesClusterID,omitempty"`

	//The aws provider creates an inbound rule per load balancer on the node security
	//group. However, this can run into the AWS security group rule limit of 50 if
	//many LoadBalancers are created.
	//
	//This flag disables the automatic ingress creation. It requires that the user
	//has setup a rule that allows inbound traffic on kubelet ports from the
	//local VPC subnet (so load balancers can access it). E.g. 10.82.0.0/16 30000-32000.
	DisableSecurityGroupIngress bool `json:"disable-security-group-ingress" yaml:"disable-security-group-ingress" ini:"DisableSecurityGroupIngress,omitempty"`

	//AWS has a hard limit of 500 security groups. For large clusters creating a security group for each ELB
	//can cause the max number of security groups to be reached. If this is set instead of creating a new
	//Security group for each ELB this security group will be used instead.
	ElbSecurityGroup string `json:"elb-security-group" yaml:"elb-security-group" ini:"ElbSecurityGroup,omitempty"`

	//During the instantiation of an new AWS cloud provider, the detected region
	//is validated against a known set of regions.
	//
	//In a non-standard, AWS like environment (e.g. Eucalyptus), this check may
	//be undesirable.  Setting this to true will disable the check and provide
	//a warning that the check was skipped.  Please note that this is an
	//experimental feature and work-in-progress for the moment.  If you find
	//yourself in an non-AWS cloud and open an issue, please indicate that in the
	//issue body.
	DisableStrictZoneCheck bool `json:"disable-strict-zone-check" yaml:"disable-strict-zone-check" ini:"DisableStrictZoneCheck,omitempty"`
}

type MonitoringConfig struct {
	// Monitoring server provider
	Provider string `yaml:"provider,omitempty" json:"provider,omitempty" norman:"default=metrics-server"`
	// These options are NOT for configuring the Metrics-Server's addon template.
	// They are used to pass command args to the metric-server's deployment containers specifically.
	Options map[string]string `yaml:"options,omitempty" json:"options,omitempty"`
	// NodeSelector key pair
	NodeSelector map[string]string `yaml:"node_selector,omitempty" json:"nodeSelector,omitempty"`
	// Update strategy
	UpdateStrategy *DeploymentStrategy `yaml:"update_strategy,omitempty" json:"updateStrategy,omitempty"`
	// Number of monitoring addon pods
	Replicas *int32 `yaml:"replicas" json:"replicas,omitempty" norman:"default=1"`
	// Tolerations for Deployments
	Tolerations []v1.Toleration `yaml:"tolerations,omitempty" json:"tolerations,omitempty"`
	// Priority class name for Metrics-Server's "metrics-server" deployment
	MetricsServerPriorityClassName string `yaml:"metrics_server_priority_class_name,omitempty" json:"metricsServerPriorityClassName,omitempty"`
}

type RestoreConfig struct {
	Restore      bool   `yaml:"restore,omitempty" json:"restore,omitempty"`
	SnapshotName string `yaml:"snapshot_name,omitempty" json:"snapshotName,omitempty"`
}
type RotateCertificates struct {
	// Rotate CA Certificates
	CACertificates bool `json:"caCertificates,omitempty"`
	// Services to rotate their certs
	Services []string `json:"services,omitempty" norman:"type=enum,options=etcd|kubelet|kube-apiserver|kube-proxy|kube-scheduler|kube-controller-manager"`
}

type DNSConfig struct {
	// DNS provider
	Provider string `yaml:"provider,omitempty" json:"provider,omitempty"`
	// DNS config options
	Options map[string]string `yaml:"options,omitempty" json:"options,omitempty"`
	// Upstream nameservers
	UpstreamNameservers []string `yaml:"upstreamnameservers,omitempty" json:"upstreamnameservers,omitempty"`
	// ReverseCIDRs
	ReverseCIDRs []string `yaml:"reversecidrs,omitempty" json:"reversecidrs,omitempty"`
	// Stubdomains
	StubDomains map[string][]string `yaml:"stubdomains,omitempty" json:"stubdomains,omitempty"`
	// NodeSelector key pair
	NodeSelector map[string]string `yaml:"node_selector,omitempty" json:"nodeSelector,omitempty"`
	// Nodelocal DNS
	Nodelocal *Nodelocal `yaml:"nodelocal,omitempty" json:"nodelocal,omitempty"`
	// Update strategy
	UpdateStrategy *DeploymentStrategy `yaml:"update_strategy,omitempty" json:"updateStrategy,omitempty"`
	// Autoscaler fields to determine number of dns replicas
	LinearAutoscalerParams *LinearAutoscalerParams `yaml:"linear_autoscaler_params,omitempty" json:"linearAutoscalerParams,omitempty"`
	// Tolerations for Deployments
	Tolerations []v1.Toleration `yaml:"tolerations,omitempty" json:"tolerations,omitempty"`
}

type Nodelocal struct {
	// link-local IP for nodelocal DNS
	IPAddress string `yaml:"ip_address,omitempty" json:"ipAddress,omitempty"`
	// Nodelocal DNS daemonset upgrade strategy
	UpdateStrategy *DaemonSetUpdateStrategy `yaml:"update_strategy,omitempty" json:"updateStrategy,omitempty"`
	// NodeSelector key pair
	NodeSelector map[string]string `yaml:"node_selector,omitempty" json:"nodeSelector,omitempty"`
	// Priority class name for NodeLocal's "node-local-dns" daemonset
	NodeLocalDNSPriorityClassName string `yaml:"node_local_dns_priority_class_name,omitempty" json:"nodeLocalDnsPriorityClassName,omitempty"`
}

// LinearAutoscalerParams contains fields expected by the cluster-proportional-autoscaler https://github.com/kubernetes-incubator/cluster-proportional-autoscaler/blob/0c61e63fc81449abdd52315aa27179a17e5d1580/pkg/autoscaler/controller/linearcontroller/linear_controller.go#L50
type LinearAutoscalerParams struct {
	CoresPerReplica           float64 `yaml:"cores_per_replica,omitempty" json:"coresPerReplica,omitempty" norman:"default=128"`
	NodesPerReplica           float64 `yaml:"nodes_per_replica,omitempty" json:"nodesPerReplica,omitempty" norman:"default=4"`
	Min                       int     `yaml:"min" json:"min,omitempty" norman:"default=1"`
	Max                       int     `yaml:"max" json:"max,omitempty"`
	PreventSinglePointFailure bool    `yaml:"prevent_single_point_failure,omitempty" json:"preventSinglePointFailure,omitempty" norman:"default=true"`
}

type RKETaint struct {
	Key       string         `json:"key,omitempty" yaml:"key,omitempty"`
	Value     string         `json:"value,omitempty" yaml:"value,omitempty"`
	Effect    v1.TaintEffect `json:"effect,omitempty" yaml:"effect,omitempty"`
	TimeAdded *metav1.Time   `json:"timeAdded,omitempty" yaml:"timeAdded,omitempty"`
}

type SecretsEncryptionConfig struct {
	// Enable/disable secrets encryption provider config
	Enabled bool `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	// Custom Encryption Provider configuration object
	CustomConfig *configv1.EncryptionConfiguration `yaml:"custom_config,omitepmty" json:"customConfig,omitempty" norman:"type=map[json]"`
}

type File struct {
	Name     string `json:"name,omitempty"`
	Contents string `json:"contents,omitempty"`
}

type NodeDrainInput struct {
	// Drain node even if there are pods not managed by a ReplicationController, Job, or DaemonSet
	// Drain will not proceed without Force set to true if there are such pods
	Force bool `yaml:"force,omitempty" json:"force,omitempty"`
	// If there are DaemonSet-managed pods, drain will not proceed without IgnoreDaemonSets set to true
	// (even when set to true, kubectl won't delete pods - so setting default to true)
	IgnoreDaemonSets *bool `yaml:"ignore_daemonsets,omitempty" json:"ignoreDaemonSets,omitempty" norman:"default=true"`
	// Continue even if there are pods using emptyDir
	DeleteLocalData bool `yaml:"delete_local_data,omitempty" json:"deleteLocalData,omitempty"`
	//Period of time in seconds given to each pod to terminate gracefully.
	// If negative, the default value specified in the pod will be used
	GracePeriod int `yaml:"grace_period,omitempty" json:"gracePeriod,omitempty" norman:"default=-1"`
	// Time to wait (in seconds) before giving up for one try
	Timeout int `yaml:"timeout" json:"timeout" norman:"min=1,max=10800,default=120"`
}
