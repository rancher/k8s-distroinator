package cluster

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/blang/semver"
	"github.com/rancher/rke/cloudprovider"
	"github.com/rancher/rke/docker"
	"github.com/rancher/rke/k8s"
	"github.com/rancher/rke/log"
	"github.com/rancher/rke/metadata"
	"github.com/rancher/rke/services"
	"github.com/rancher/rke/templates"
	v3 "github.com/rancher/rke/types"
	"github.com/rancher/rke/util"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	apiserverv1alpha1 "k8s.io/apiserver/pkg/apis/apiserver/v1alpha1"
	auditv1 "k8s.io/apiserver/pkg/apis/audit/v1"
)

const (
	DefaultServiceClusterIPRange = "10.43.0.0/16"
	DefaultNodePortRange         = "30000-32767"
	DefaultClusterCIDR           = "10.42.0.0/16"
	DefaultClusterDNSService     = "10.43.0.10"
	DefaultClusterDomain         = "cluster.local"
	DefaultClusterName           = "local"
	DefaultClusterSSHKeyPath     = "~/.ssh/id_rsa"

	DefaultSSHPort        = "22"
	DefaultDockerSockPath = "/var/run/docker.sock"

	DefaultAuthStrategy      = "x509"
	DefaultAuthorizationMode = "rbac"

	DefaultAuthnWebhookFile  = templates.AuthnWebhook
	DefaultAuthnCacheTimeout = "5s"

	DefaultNetworkPlugin        = "canal"
	DefaultNetworkCloudProvider = "none"

	DefaultIngressController             = "nginx"
	DefaultEtcdBackupCreationPeriod      = "12h"
	DefaultEtcdBackupRetentionPeriod     = "72h"
	DefaultEtcdSnapshot                  = true
	DefaultMonitoringProvider            = "metrics-server"
	DefaultEtcdBackupConfigIntervalHours = 12
	DefaultEtcdBackupConfigRetention     = 6
	DefaultEtcdBackupConfigTimeout       = docker.WaitTimeout

	DefaultDNSProvider = "kube-dns"
	K8sVersionCoreDNS  = "1.14.0"

	DefaultEtcdHeartbeatIntervalName  = "heartbeat-interval"
	DefaultEtcdHeartbeatIntervalValue = "500"
	DefaultEtcdElectionTimeoutName    = "election-timeout"
	DefaultEtcdElectionTimeoutValue   = "5000"

	DefaultFlannelBackendVxLan     = "vxlan"
	DefaultFlannelBackendVxLanPort = "8472"
	DefaultFlannelBackendVxLanVNI  = "1"

	DefaultCalicoFlexVolPluginDirectory = "/usr/libexec/kubernetes/kubelet-plugins/volume/exec/nodeagent~uds"

	DefaultCanalFlexVolPluginDirectory = "/usr/libexec/kubernetes/kubelet-plugins/volume/exec/nodeagent~uds"

	DefaultAciApicRefreshTime          = "1200"
	DefaultAciOVSMemoryLimit           = "1Gi"
	DefaultAciImagePullPolicy          = "Always"
	DefaultAciServiceMonitorInterval   = "5"
	DefaultAciPBRTrackingNonSnat       = "false"
	DefaultAciInstallIstio             = "false"
	DefaultAciIstioProfile             = "demo"
	DefaultAciDropLogEnable            = "true"
	DefaultAciControllerLogLevel       = "info"
	DefaultAciHostAgentLogLevel        = "info"
	DefaultAciOpflexAgentLogLevel      = "info"
	DefaultAciUseAciCniPriorityClass   = "false"
	DefaultAciNoPriorityClass          = "false"
	DefaultAciMaxNodesSvcGraph         = "32"
	DefaultAciSnatContractScope        = "global"
	DefaultAciSnatNamespace            = "aci-containers-system"
	DefaultAciCApic                    = "false"
	DefaultAciPodSubnetChunkSize       = "32"
	DefaultAciSnatPortRangeStart       = "5000"
	DefaultAciSnatPortRangeEnd         = "65000"
	DefaultAciSnatPortsPerNode         = "3000"
	DefaultAciUseHostNetnsVolume       = "false"
	DefaultAciRunGbpContainer          = "false"
	DefaultAciRunOpflexServerContainer = "false"
	DefaultAciUseAciAnywhereCRD        = "false"
	DefaultAciEnableEndpointSlice      = "false"
	DefaultAciOpflexClientSSL          = "true"
	DefaultAciUsePrivilegedContainer   = "false"
	DefaultAciUseOpflexServerVolume    = "false"

	DefaultKubeRouterRunServiceProxy = false
	DefaultKubeRouterRunFirewall     = false

	KubeAPIArgAdmissionControlConfigFile             = "admission-control-config-file"
	DefaultKubeAPIArgAdmissionControlConfigFileValue = "/etc/kubernetes/admission.yaml"

	EventRateLimitPluginName = "EventRateLimit"

	KubeAPIArgAuditLogPath                = "audit-log-path"
	KubeAPIArgAuditLogMaxAge              = "audit-log-maxage"
	KubeAPIArgAuditLogMaxBackup           = "audit-log-maxbackup"
	KubeAPIArgAuditLogMaxSize             = "audit-log-maxsize"
	KubeAPIArgAuditLogFormat              = "audit-log-format"
	KubeAPIArgAuditPolicyFile             = "audit-policy-file"
	DefaultKubeAPIArgAuditLogPathValue    = "/var/log/kube-audit/audit-log.json"
	DefaultKubeAPIArgAuditPolicyFileValue = "/etc/kubernetes/audit-policy.yaml"

	DefaultMaxUnavailableWorker       = "10%"
	DefaultMaxUnavailableControlplane = "1"
	DefaultNodeDrainTimeout           = 120
	DefaultNodeDrainGracePeriod       = -1
	DefaultHTTPPort                   = 80
	DefaultHTTPSPort                  = 443
	DefaultNetworkMode                = "hostNetwork"

	DefaultServiceKubeproxyEnabled = true
)

var (
	DefaultNodeDrainIgnoreDaemonsets      = true
	DefaultDaemonSetMaxUnavailable        = intstr.FromInt(1)
	DefaultDeploymentUpdateStrategyParams = intstr.FromString("25%")
	DefaultDaemonSetUpdateStrategy        = v3.DaemonSetUpdateStrategy{
		Strategy:      appsv1.RollingUpdateDaemonSetStrategyType,
		RollingUpdate: &appsv1.RollingUpdateDaemonSet{MaxUnavailable: &DefaultDaemonSetMaxUnavailable},
	}
	DefaultDeploymentUpdateStrategy = v3.DeploymentStrategy{
		Strategy: appsv1.RollingUpdateDeploymentStrategyType,
		RollingUpdate: &appsv1.RollingUpdateDeployment{
			MaxUnavailable: &DefaultDeploymentUpdateStrategyParams,
			MaxSurge:       &DefaultDeploymentUpdateStrategyParams,
		},
	}
	DefaultClusterProportionalAutoscalerLinearParams = v3.LinearAutoscalerParams{CoresPerReplica: 128, NodesPerReplica: 4, Min: 1, PreventSinglePointFailure: true}
	DefaultMonitoringAddonReplicas                   = int32(1)
)

type ExternalFlags struct {
	CertificateDir   string
	ClusterFilePath  string
	DinD             bool
	ConfigDir        string
	CustomCerts      bool
	DisablePortCheck bool
	GenerateCSR      bool
	Local            bool
	UpdateOnly       bool
	UseLocalState    bool
}

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

func (c *Cluster) setClusterDefaults(ctx context.Context, flags ExternalFlags) error {
	if len(c.SSHKeyPath) == 0 {
		c.SSHKeyPath = DefaultClusterSSHKeyPath
	}
	// Default Path prefix
	if len(c.PrefixPath) == 0 {
		c.PrefixPath = "/"
	}
	// Set bastion/jump host defaults
	if len(c.BastionHost.Address) > 0 {
		if len(c.BastionHost.Port) == 0 {
			c.BastionHost.Port = DefaultSSHPort
		}
		if len(c.BastionHost.SSHKeyPath) == 0 {
			c.BastionHost.SSHKeyPath = c.SSHKeyPath
		}
		c.BastionHost.SSHAgentAuth = c.SSHAgentAuth

	}
	for i, host := range c.Nodes {
		if len(host.InternalAddress) == 0 {
			c.Nodes[i].InternalAddress = c.Nodes[i].Address
		}
		if len(host.HostnameOverride) == 0 {
			// This is a temporary modification
			c.Nodes[i].HostnameOverride = c.Nodes[i].Address
		}
		if len(host.SSHKeyPath) == 0 {
			c.Nodes[i].SSHKeyPath = c.SSHKeyPath
		}
		if len(host.Port) == 0 {
			c.Nodes[i].Port = DefaultSSHPort
		}

		c.Nodes[i].HostnameOverride = strings.ToLower(c.Nodes[i].HostnameOverride)
		// For now, you can set at the global level only.
		c.Nodes[i].SSHAgentAuth = c.SSHAgentAuth
	}

	if len(c.Authorization.Mode) == 0 {
		c.Authorization.Mode = DefaultAuthorizationMode
	}
	if c.Services.KubeAPI.PodSecurityPolicy && c.Authorization.Mode != services.RBACAuthorizationMode {
		log.Warnf(ctx, "PodSecurityPolicy can't be enabled with RBAC support disabled")
		c.Services.KubeAPI.PodSecurityPolicy = false
	}
	if len(c.Ingress.Provider) == 0 {
		c.Ingress.Provider = DefaultIngressController
	}
	if len(c.ClusterName) == 0 {
		c.ClusterName = DefaultClusterName
	}
	if len(c.Version) == 0 {
		c.Version = metadata.DefaultK8sVersion
	}
	if c.AddonJobTimeout == 0 {
		c.AddonJobTimeout = k8s.DefaultTimeout
	}
	if len(c.Monitoring.Provider) == 0 {
		c.Monitoring.Provider = DefaultMonitoringProvider
	}
	//set docker private registry URL
	for _, pr := range c.PrivateRegistries {
		if pr.URL == "" {
			pr.URL = docker.DockerRegistryURL
		}
		c.PrivateRegistriesMap[pr.URL] = pr
	}

	err := c.setClusterImageDefaults()
	if err != nil {
		return err
	}

	if c.RancherKubernetesEngineConfig.RotateCertificates != nil ||
		flags.CustomCerts {
		c.ForceDeployCerts = true
	}

	err = c.setClusterDNSDefaults()
	if err != nil {
		return err
	}
	c.setClusterServicesDefaults()
	c.setClusterNetworkDefaults()
	c.setClusterAuthnDefaults()
	c.setNodeUpgradeStrategy()
	c.setAddonsDefaults()
	return nil
}

func (c *Cluster) setNodeUpgradeStrategy() {
	if c.UpgradeStrategy == nil {
		logrus.Debugf("No input provided for maxUnavailableWorker, setting it to default value of %v percent", strings.TrimRight(DefaultMaxUnavailableWorker, "%"))
		logrus.Debugf("No input provided for maxUnavailableControlplane, setting it to default value of %v", DefaultMaxUnavailableControlplane)
		c.UpgradeStrategy = &v3.NodeUpgradeStrategy{
			MaxUnavailableWorker:       DefaultMaxUnavailableWorker,
			MaxUnavailableControlplane: DefaultMaxUnavailableControlplane,
		}
		return
	}
	setDefaultIfEmpty(&c.UpgradeStrategy.MaxUnavailableWorker, DefaultMaxUnavailableWorker)
	setDefaultIfEmpty(&c.UpgradeStrategy.MaxUnavailableControlplane, DefaultMaxUnavailableControlplane)
	if c.UpgradeStrategy.Drain != nil && *c.UpgradeStrategy.Drain {
		return
	}
	if c.UpgradeStrategy.DrainInput == nil {
		c.UpgradeStrategy.DrainInput = &v3.NodeDrainInput{
			IgnoreDaemonSets: &DefaultNodeDrainIgnoreDaemonsets,
			// default to 120 seems to work better for controlplane nodes
			Timeout: DefaultNodeDrainTimeout,
			//Period of time in seconds given to each pod to terminate gracefully.
			// If negative, the default value specified in the pod will be used
			GracePeriod: DefaultNodeDrainGracePeriod,
		}
	}
}

func (c *Cluster) setClusterServicesDefaults() {
	// We don't accept per service images anymore.
	c.Services.KubeAPI.Image = c.SystemImages.Kubernetes
	c.Services.Scheduler.Image = c.SystemImages.Kubernetes
	c.Services.KubeController.Image = c.SystemImages.Kubernetes
	c.Services.Kubelet.Image = c.SystemImages.Kubernetes
	c.Services.Kubeproxy.Image = c.SystemImages.Kubernetes
	c.Services.Etcd.Image = c.SystemImages.Etcd

	// enable etcd snapshots by default
	if c.Services.Etcd.Snapshot == nil {
		defaultSnapshot := DefaultEtcdSnapshot
		c.Services.Etcd.Snapshot = &defaultSnapshot
	}

	// enable kubeproxy by default
	if c.Services.Kubeproxy.Enabled == nil {
		defaultServiceKubeproxyEnabled := DefaultServiceKubeproxyEnabled
		c.Services.Kubeproxy.Enabled = &defaultServiceKubeproxyEnabled
	}

	serviceConfigDefaultsMap := map[*string]string{
		&c.Services.KubeAPI.ServiceClusterIPRange:        DefaultServiceClusterIPRange,
		&c.Services.KubeAPI.ServiceNodePortRange:         DefaultNodePortRange,
		&c.Services.KubeController.ServiceClusterIPRange: DefaultServiceClusterIPRange,
		&c.Services.KubeController.ClusterCIDR:           DefaultClusterCIDR,
		&c.Services.Kubelet.ClusterDNSServer:             DefaultClusterDNSService,
		&c.Services.Kubelet.ClusterDomain:                DefaultClusterDomain,
		&c.Services.Kubelet.InfraContainerImage:          c.SystemImages.PodInfraContainer,
		&c.Services.Etcd.Creation:                        DefaultEtcdBackupCreationPeriod,
		&c.Services.Etcd.Retention:                       DefaultEtcdBackupRetentionPeriod,
	}
	for k, v := range serviceConfigDefaultsMap {
		setDefaultIfEmpty(k, v)
	}
	// Add etcd timeouts
	if c.Services.Etcd.ExtraArgs == nil {
		c.Services.Etcd.ExtraArgs = make(map[string]string)
	}
	if _, ok := c.Services.Etcd.ExtraArgs[DefaultEtcdElectionTimeoutName]; !ok {
		c.Services.Etcd.ExtraArgs[DefaultEtcdElectionTimeoutName] = DefaultEtcdElectionTimeoutValue
	}
	if _, ok := c.Services.Etcd.ExtraArgs[DefaultEtcdHeartbeatIntervalName]; !ok {
		c.Services.Etcd.ExtraArgs[DefaultEtcdHeartbeatIntervalName] = DefaultEtcdHeartbeatIntervalValue
	}

	if c.Services.Etcd.BackupConfig != nil &&
		(c.Services.Etcd.BackupConfig.Enabled == nil ||
			(c.Services.Etcd.BackupConfig.Enabled != nil && *c.Services.Etcd.BackupConfig.Enabled)) {
		if c.Services.Etcd.BackupConfig.IntervalHours == 0 {
			c.Services.Etcd.BackupConfig.IntervalHours = DefaultEtcdBackupConfigIntervalHours
		}
		if c.Services.Etcd.BackupConfig.Retention == 0 {
			c.Services.Etcd.BackupConfig.Retention = DefaultEtcdBackupConfigRetention
		}
		if c.Services.Etcd.BackupConfig.Timeout == 0 {
			c.Services.Etcd.BackupConfig.Timeout = DefaultEtcdBackupConfigTimeout
		}
	}

	if _, ok := c.Services.KubeAPI.ExtraArgs[KubeAPIArgAdmissionControlConfigFile]; !ok {
		if c.Services.KubeAPI.EventRateLimit != nil &&
			c.Services.KubeAPI.EventRateLimit.Enabled &&
			c.Services.KubeAPI.EventRateLimit.Configuration == nil {
			c.Services.KubeAPI.EventRateLimit.Configuration = newDefaultEventRateLimitConfig()
		}
	}

	enableKubeAPIAuditLog, err := checkVersionNeedsKubeAPIAuditLog(c.Version)
	if err != nil {
		logrus.Warnf("Can not determine if cluster version [%s] needs to have kube-api audit log enabled: %v", c.Version, err)
	}
	if enableKubeAPIAuditLog {
		logrus.Debugf("Enabling kube-api audit log for cluster version [%s]", c.Version)

		if c.Services.KubeAPI.AuditLog == nil {
			c.Services.KubeAPI.AuditLog = &v3.AuditLog{Enabled: true}
		}

	}
	if c.Services.KubeAPI.AuditLog != nil &&
		c.Services.KubeAPI.AuditLog.Enabled {
		if c.Services.KubeAPI.AuditLog.Configuration == nil {
			alc := newDefaultAuditLogConfig()
			c.Services.KubeAPI.AuditLog.Configuration = alc
		} else {
			if c.Services.KubeAPI.AuditLog.Configuration.Policy == nil {
				c.Services.KubeAPI.AuditLog.Configuration.Policy = newDefaultAuditPolicy()
			}
		}
	}
}

func newDefaultAuditPolicy() *auditv1.Policy {
	p := &auditv1.Policy{
		TypeMeta: v1.TypeMeta{
			Kind:       "Policy",
			APIVersion: auditv1.SchemeGroupVersion.String(),
		},
		Rules: []auditv1.PolicyRule{
			{
				Level: "Metadata",
			},
		},
		OmitStages: nil,
	}
	return p
}

func newDefaultAuditLogConfig() *v3.AuditLogConfig {
	p := newDefaultAuditPolicy()
	c := &v3.AuditLogConfig{
		MaxAge:    30,
		MaxBackup: 10,
		MaxSize:   100,
		Path:      DefaultKubeAPIArgAuditLogPathValue,
		Format:    "json",
		Policy:    p,
	}
	return c
}

func getEventRateLimitPluginFromConfig(c *v3.Configuration) (apiserverv1alpha1.AdmissionPluginConfiguration, error) {
	plugin := apiserverv1alpha1.AdmissionPluginConfiguration{
		Name: EventRateLimitPluginName,
		Configuration: &runtime.Unknown{
			ContentType: "application/json",
		},
	}

	cBytes, err := json.Marshal(c)
	if err != nil {
		return plugin, fmt.Errorf("error marshalling eventratelimit config: %v", err)
	}
	plugin.Configuration.Raw = cBytes

	return plugin, nil
}

func newDefaultEventRateLimitConfig() *v3.Configuration {
	return &v3.Configuration{
		TypeMeta: v1.TypeMeta{
			Kind:       "Configuration",
			APIVersion: "eventratelimit.admission.k8s.io/v1alpha1",
		},
		Limits: []v3.Limit{
			{
				Type:  v3.ServerLimitType,
				QPS:   5000,
				Burst: 20000,
			},
		},
	}
}

func newDefaultEventRateLimitPlugin() (apiserverv1alpha1.AdmissionPluginConfiguration, error) {
	plugin := apiserverv1alpha1.AdmissionPluginConfiguration{
		Name: EventRateLimitPluginName,
		Configuration: &runtime.Unknown{
			ContentType: "application/json",
		},
	}

	c := newDefaultEventRateLimitConfig()
	cBytes, err := json.Marshal(c)
	if err != nil {
		return plugin, fmt.Errorf("error marshalling eventratelimit config: %v", err)
	}
	plugin.Configuration.Raw = cBytes

	return plugin, nil
}

func newDefaultAdmissionConfiguration() (*apiserverv1alpha1.AdmissionConfiguration, error) {
	var admissionConfiguration *apiserverv1alpha1.AdmissionConfiguration
	admissionConfiguration = &apiserverv1alpha1.AdmissionConfiguration{
		TypeMeta: v1.TypeMeta{
			Kind:       "AdmissionConfiguration",
			APIVersion: apiserverv1alpha1.SchemeGroupVersion.String(),
		},
	}
	return admissionConfiguration, nil
}

func (c *Cluster) setClusterImageDefaults() error {
	var privRegURL string

	imageDefaults, ok := metadata.K8sVersionToRKESystemImages[c.Version]
	if !ok {
		return nil
	}

	for _, privReg := range c.PrivateRegistries {
		if privReg.IsDefault {
			privRegURL = privReg.URL
			break
		}
	}
	systemImagesDefaultsMap := map[*string]string{
		&c.SystemImages.Alpine:                    d(imageDefaults.Alpine, privRegURL),
		&c.SystemImages.NginxProxy:                d(imageDefaults.NginxProxy, privRegURL),
		&c.SystemImages.CertDownloader:            d(imageDefaults.CertDownloader, privRegURL),
		&c.SystemImages.KubeDNS:                   d(imageDefaults.KubeDNS, privRegURL),
		&c.SystemImages.KubeDNSSidecar:            d(imageDefaults.KubeDNSSidecar, privRegURL),
		&c.SystemImages.DNSmasq:                   d(imageDefaults.DNSmasq, privRegURL),
		&c.SystemImages.KubeDNSAutoscaler:         d(imageDefaults.KubeDNSAutoscaler, privRegURL),
		&c.SystemImages.CoreDNS:                   d(imageDefaults.CoreDNS, privRegURL),
		&c.SystemImages.CoreDNSAutoscaler:         d(imageDefaults.CoreDNSAutoscaler, privRegURL),
		&c.SystemImages.KubernetesServicesSidecar: d(imageDefaults.KubernetesServicesSidecar, privRegURL),
		&c.SystemImages.Etcd:                      d(imageDefaults.Etcd, privRegURL),
		&c.SystemImages.Kubernetes:                d(imageDefaults.Kubernetes, privRegURL),
		&c.SystemImages.PodInfraContainer:         d(imageDefaults.PodInfraContainer, privRegURL),
		&c.SystemImages.Flannel:                   d(imageDefaults.Flannel, privRegURL),
		&c.SystemImages.FlannelCNI:                d(imageDefaults.FlannelCNI, privRegURL),
		&c.SystemImages.CalicoNode:                d(imageDefaults.CalicoNode, privRegURL),
		&c.SystemImages.CalicoCNI:                 d(imageDefaults.CalicoCNI, privRegURL),
		&c.SystemImages.CalicoCtl:                 d(imageDefaults.CalicoCtl, privRegURL),
		&c.SystemImages.CalicoControllers:         d(imageDefaults.CalicoControllers, privRegURL),
		&c.SystemImages.CalicoFlexVol:             d(imageDefaults.CalicoFlexVol, privRegURL),
		&c.SystemImages.CanalNode:                 d(imageDefaults.CanalNode, privRegURL),
		&c.SystemImages.CanalCNI:                  d(imageDefaults.CanalCNI, privRegURL),
		&c.SystemImages.CanalControllers:          d(imageDefaults.CanalControllers, privRegURL),
		&c.SystemImages.CanalFlannel:              d(imageDefaults.CanalFlannel, privRegURL),
		&c.SystemImages.CanalFlexVol:              d(imageDefaults.CanalFlexVol, privRegURL),
		&c.SystemImages.WeaveNode:                 d(imageDefaults.WeaveNode, privRegURL),
		&c.SystemImages.WeaveCNI:                  d(imageDefaults.WeaveCNI, privRegURL),
		&c.SystemImages.Ingress:                   d(imageDefaults.Ingress, privRegURL),
		&c.SystemImages.IngressBackend:            d(imageDefaults.IngressBackend, privRegURL),
		&c.SystemImages.MetricsServer:             d(imageDefaults.MetricsServer, privRegURL),
		&c.SystemImages.Nodelocal:                 d(imageDefaults.Nodelocal, privRegURL),
		&c.SystemImages.AciCniDeployContainer:     d(imageDefaults.AciCniDeployContainer, privRegURL),
		&c.SystemImages.AciHostContainer:          d(imageDefaults.AciHostContainer, privRegURL),
		&c.SystemImages.AciOpflexContainer:        d(imageDefaults.AciOpflexContainer, privRegURL),
		&c.SystemImages.AciMcastContainer:         d(imageDefaults.AciMcastContainer, privRegURL),
		&c.SystemImages.AciOpenvSwitchContainer:   d(imageDefaults.AciOpenvSwitchContainer, privRegURL),
		&c.SystemImages.AciControllerContainer:    d(imageDefaults.AciControllerContainer, privRegURL),
		&c.SystemImages.AciOpflexServerContainer:  d(imageDefaults.AciOpflexServerContainer, privRegURL),
		&c.SystemImages.AciGbpServerContainer:     d(imageDefaults.AciGbpServerContainer, privRegURL),

		// this's a stopgap, we could drop this after https://github.com/kubernetes/kubernetes/pull/75618 merged
		&c.SystemImages.WindowsPodInfraContainer: d(imageDefaults.WindowsPodInfraContainer, privRegURL),
	}

	for k, v := range systemImagesDefaultsMap {
		setDefaultIfEmpty(k, v)
	}

	return nil
}

func (c *Cluster) setClusterDNSDefaults() error {
	if c.DNS != nil && len(c.DNS.Provider) != 0 {
		return nil
	}
	clusterSemVer, err := util.StrToSemVer(c.Version)
	if err != nil {
		return err
	}
	logrus.Debugf("No DNS provider configured, setting default based on cluster version [%s]", clusterSemVer)
	K8sVersionCoreDNSSemVer, err := util.StrToSemVer(K8sVersionCoreDNS)
	if err != nil {
		return err
	}
	// Default DNS provider for cluster version 1.14.0 and higher is coredns
	ClusterDNSProvider := CoreDNSProvider
	// If cluster version is less than 1.14.0 (K8sVersionCoreDNSSemVer), use kube-dns
	if clusterSemVer.LessThan(*K8sVersionCoreDNSSemVer) {
		logrus.Debugf("Cluster version [%s] is less than version [%s], using DNS provider [%s]", clusterSemVer, K8sVersionCoreDNSSemVer, DefaultDNSProvider)
		ClusterDNSProvider = DefaultDNSProvider
	}
	if c.DNS == nil {
		c.DNS = &v3.DNSConfig{}
	}
	c.DNS.Provider = ClusterDNSProvider
	logrus.Debugf("DNS provider set to [%s]", ClusterDNSProvider)
	if c.DNS.Options == nil {
		// don't break if the user didn't define options
		c.DNS.Options = make(map[string]string)
	}
	return nil
}

func (c *Cluster) setClusterNetworkDefaults() {
	setDefaultIfEmpty(&c.Network.Plugin, DefaultNetworkPlugin)

	if c.Network.Options == nil {
		// don't break if the user didn't define options
		c.Network.Options = make(map[string]string)
	}
	networkPluginConfigDefaultsMap := make(map[string]string)
	// This is still needed because RKE doesn't use c.Network.*NetworkProvider, that's a rancher type
	switch c.Network.Plugin {
	case CalicoNetworkPlugin:
		networkPluginConfigDefaultsMap = map[string]string{
			CalicoCloudProvider:          DefaultNetworkCloudProvider,
			CalicoFlexVolPluginDirectory: DefaultCalicoFlexVolPluginDirectory,
		}
	case FlannelNetworkPlugin:
		networkPluginConfigDefaultsMap = map[string]string{
			FlannelBackendType:                 DefaultFlannelBackendVxLan,
			FlannelBackendPort:                 DefaultFlannelBackendVxLanPort,
			FlannelBackendVxLanNetworkIdentify: DefaultFlannelBackendVxLanVNI,
		}
	case CanalNetworkPlugin:
		networkPluginConfigDefaultsMap = map[string]string{
			CanalFlannelBackendType:                 DefaultFlannelBackendVxLan,
			CanalFlannelBackendPort:                 DefaultFlannelBackendVxLanPort,
			CanalFlannelBackendVxLanNetworkIdentify: DefaultFlannelBackendVxLanVNI,
			CanalFlexVolPluginDirectory:             DefaultCanalFlexVolPluginDirectory,
		}
	case AciNetworkPlugin:
		networkPluginConfigDefaultsMap = map[string]string{
			AciOVSMemoryLimit:           DefaultAciOVSMemoryLimit,
			AciImagePullPolicy:          DefaultAciImagePullPolicy,
			AciPBRTrackingNonSnat:       DefaultAciPBRTrackingNonSnat,
			AciInstallIstio:             DefaultAciInstallIstio,
			AciIstioProfile:             DefaultAciIstioProfile,
			AciDropLogEnable:            DefaultAciDropLogEnable,
			AciControllerLogLevel:       DefaultAciControllerLogLevel,
			AciHostAgentLogLevel:        DefaultAciHostAgentLogLevel,
			AciOpflexAgentLogLevel:      DefaultAciOpflexAgentLogLevel,
			AciApicRefreshTime:          DefaultAciApicRefreshTime,
			AciServiceMonitorInterval:   DefaultAciServiceMonitorInterval,
			AciUseAciCniPriorityClass:   DefaultAciUseAciCniPriorityClass,
			AciNoPriorityClass:          DefaultAciNoPriorityClass,
			AciMaxNodesSvcGraph:         DefaultAciMaxNodesSvcGraph,
			AciSnatContractScope:        DefaultAciSnatContractScope,
			AciPodSubnetChunkSize:       DefaultAciPodSubnetChunkSize,
			AciEnableEndpointSlice:      DefaultAciEnableEndpointSlice,
			AciSnatNamespace:            DefaultAciSnatNamespace,
			AciSnatPortRangeStart:       DefaultAciSnatPortRangeStart,
			AciSnatPortRangeEnd:         DefaultAciSnatPortRangeEnd,
			AciSnatPortsPerNode:         DefaultAciSnatPortsPerNode,
			AciOpflexClientSSL:          DefaultAciOpflexClientSSL,
			AciUsePrivilegedContainer:   DefaultAciUsePrivilegedContainer,
			AciUseOpflexServerVolume:    DefaultAciUseOpflexServerVolume,
			AciUseHostNetnsVolume:       DefaultAciUseHostNetnsVolume,
			AciCApic:                    DefaultAciCApic,
			AciUseAciAnywhereCRD:        DefaultAciUseAciAnywhereCRD,
			AciRunGbpContainer:          DefaultAciRunGbpContainer,
			AciRunOpflexServerContainer: DefaultAciRunOpflexServerContainer,
		}
	}
	if c.Network.CalicoNetworkProvider != nil {
		setDefaultIfEmpty(&c.Network.CalicoNetworkProvider.CloudProvider, DefaultNetworkCloudProvider)
		networkPluginConfigDefaultsMap[CalicoCloudProvider] = c.Network.CalicoNetworkProvider.CloudProvider
	}
	if c.Network.FlannelNetworkProvider != nil {
		networkPluginConfigDefaultsMap[FlannelIface] = c.Network.FlannelNetworkProvider.Iface

	}
	if c.Network.CanalNetworkProvider != nil {
		networkPluginConfigDefaultsMap[CanalIface] = c.Network.CanalNetworkProvider.Iface
	}
	if c.Network.WeaveNetworkProvider != nil {
		networkPluginConfigDefaultsMap[WeavePassword] = c.Network.WeaveNetworkProvider.Password
	}
	if c.Network.AciNetworkProvider != nil {
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.OVSMemoryLimit, DefaultAciOVSMemoryLimit)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.ImagePullPolicy, DefaultAciImagePullPolicy)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.PBRTrackingNonSnat, DefaultAciPBRTrackingNonSnat)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.InstallIstio, DefaultAciInstallIstio)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.IstioProfile, DefaultAciIstioProfile)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.DropLogEnable, DefaultAciDropLogEnable)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.ControllerLogLevel, DefaultAciControllerLogLevel)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.HostAgentLogLevel, DefaultAciHostAgentLogLevel)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.OpflexAgentLogLevel, DefaultAciOpflexAgentLogLevel)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.ApicRefreshTime, DefaultAciApicRefreshTime)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.ServiceMonitorInterval, DefaultAciServiceMonitorInterval)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.NoPriorityClass, DefaultAciNoPriorityClass)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.MaxNodesSvcGraph, DefaultAciMaxNodesSvcGraph)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.SnatContractScope, DefaultAciSnatContractScope)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.PodSubnetChunkSize, DefaultAciPodSubnetChunkSize)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.EnableEndpointSlice, DefaultAciEnableEndpointSlice)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.SnatNamespace, DefaultAciSnatNamespace)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.SnatPortRangeStart, DefaultAciSnatPortRangeStart)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.SnatPortRangeEnd, DefaultAciSnatPortRangeEnd)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.SnatPortsPerNode, DefaultAciSnatPortsPerNode)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.OpflexClientSSL, DefaultAciOpflexClientSSL)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.UsePrivilegedContainer, DefaultAciUsePrivilegedContainer)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.UseOpflexServerVolume, DefaultAciUseOpflexServerVolume)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.UseHostNetnsVolume, DefaultAciUseHostNetnsVolume)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.CApic, DefaultAciCApic)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.UseAciAnywhereCRD, DefaultAciUseAciAnywhereCRD)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.RunGbpContainer, DefaultAciRunGbpContainer)
		setDefaultIfEmpty(&c.Network.AciNetworkProvider.RunOpflexServerContainer, DefaultAciRunOpflexServerContainer)
		networkPluginConfigDefaultsMap[AciOVSMemoryLimit] = c.Network.AciNetworkProvider.OVSMemoryLimit
		networkPluginConfigDefaultsMap[AciImagePullPolicy] = c.Network.AciNetworkProvider.ImagePullPolicy
		networkPluginConfigDefaultsMap[AciPBRTrackingNonSnat] = c.Network.AciNetworkProvider.PBRTrackingNonSnat
		networkPluginConfigDefaultsMap[AciInstallIstio] = c.Network.AciNetworkProvider.InstallIstio
		networkPluginConfigDefaultsMap[AciIstioProfile] = c.Network.AciNetworkProvider.IstioProfile
		networkPluginConfigDefaultsMap[AciDropLogEnable] = c.Network.AciNetworkProvider.DropLogEnable
		networkPluginConfigDefaultsMap[AciControllerLogLevel] = c.Network.AciNetworkProvider.ControllerLogLevel
		networkPluginConfigDefaultsMap[AciHostAgentLogLevel] = c.Network.AciNetworkProvider.HostAgentLogLevel
		networkPluginConfigDefaultsMap[AciOpflexAgentLogLevel] = c.Network.AciNetworkProvider.OpflexAgentLogLevel
		networkPluginConfigDefaultsMap[AciApicRefreshTime] = c.Network.AciNetworkProvider.ApicRefreshTime
		networkPluginConfigDefaultsMap[AciServiceMonitorInterval] = c.Network.AciNetworkProvider.ServiceMonitorInterval
		networkPluginConfigDefaultsMap[AciNoPriorityClass] = c.Network.AciNetworkProvider.NoPriorityClass
		networkPluginConfigDefaultsMap[AciMaxNodesSvcGraph] = c.Network.AciNetworkProvider.MaxNodesSvcGraph
		networkPluginConfigDefaultsMap[AciSnatContractScope] = c.Network.AciNetworkProvider.SnatContractScope
		networkPluginConfigDefaultsMap[AciPodSubnetChunkSize] = c.Network.AciNetworkProvider.PodSubnetChunkSize
		networkPluginConfigDefaultsMap[AciEnableEndpointSlice] = c.Network.AciNetworkProvider.EnableEndpointSlice
		networkPluginConfigDefaultsMap[AciSnatNamespace] = c.Network.AciNetworkProvider.SnatNamespace
		networkPluginConfigDefaultsMap[AciSnatPortRangeStart] = c.Network.AciNetworkProvider.SnatPortRangeStart
		networkPluginConfigDefaultsMap[AciSnatPortRangeEnd] = c.Network.AciNetworkProvider.SnatPortRangeEnd
		networkPluginConfigDefaultsMap[AciSnatPortsPerNode] = c.Network.AciNetworkProvider.SnatPortsPerNode
		networkPluginConfigDefaultsMap[AciOpflexClientSSL] = c.Network.AciNetworkProvider.OpflexClientSSL
		networkPluginConfigDefaultsMap[AciUsePrivilegedContainer] = c.Network.AciNetworkProvider.UsePrivilegedContainer
		networkPluginConfigDefaultsMap[AciUseOpflexServerVolume] = c.Network.AciNetworkProvider.UseOpflexServerVolume
		networkPluginConfigDefaultsMap[AciUseHostNetnsVolume] = c.Network.AciNetworkProvider.UseHostNetnsVolume
		networkPluginConfigDefaultsMap[AciCApic] = c.Network.AciNetworkProvider.CApic
		networkPluginConfigDefaultsMap[AciUseAciAnywhereCRD] = c.Network.AciNetworkProvider.UseAciAnywhereCRD
		networkPluginConfigDefaultsMap[AciRunGbpContainer] = c.Network.AciNetworkProvider.RunGbpContainer
		networkPluginConfigDefaultsMap[AciRunOpflexServerContainer] = c.Network.AciNetworkProvider.RunOpflexServerContainer
		networkPluginConfigDefaultsMap[AciSystemIdentifier] = c.Network.AciNetworkProvider.SystemIdentifier
		networkPluginConfigDefaultsMap[AciToken] = c.Network.AciNetworkProvider.Token
		networkPluginConfigDefaultsMap[AciApicUserName] = c.Network.AciNetworkProvider.ApicUserName
		networkPluginConfigDefaultsMap[AciApicUserKey] = c.Network.AciNetworkProvider.ApicUserKey
		networkPluginConfigDefaultsMap[AciApicUserCrt] = c.Network.AciNetworkProvider.ApicUserCrt
		networkPluginConfigDefaultsMap[AciApicRefreshTime] = c.Network.AciNetworkProvider.ApicRefreshTime
		networkPluginConfigDefaultsMap[AciVmmDomain] = c.Network.AciNetworkProvider.VmmDomain
		networkPluginConfigDefaultsMap[AciVmmController] = c.Network.AciNetworkProvider.VmmController
		networkPluginConfigDefaultsMap[AciEncapType] = c.Network.AciNetworkProvider.EncapType
		networkPluginConfigDefaultsMap[AciMcastRangeStart] = c.Network.AciNetworkProvider.McastRangeStart
		networkPluginConfigDefaultsMap[AciMcastRangeEnd] = c.Network.AciNetworkProvider.McastRangeEnd
		networkPluginConfigDefaultsMap[AciNodeSubnet] = c.Network.AciNetworkProvider.NodeSubnet
		networkPluginConfigDefaultsMap[AciAEP] = c.Network.AciNetworkProvider.AEP
		networkPluginConfigDefaultsMap[AciVRFName] = c.Network.AciNetworkProvider.VRFName
		networkPluginConfigDefaultsMap[AciVRFTenant] = c.Network.AciNetworkProvider.VRFTenant
		networkPluginConfigDefaultsMap[AciL3Out] = c.Network.AciNetworkProvider.L3Out
		networkPluginConfigDefaultsMap[AciDynamicExternalSubnet] = c.Network.AciNetworkProvider.DynamicExternalSubnet
		networkPluginConfigDefaultsMap[AciStaticExternalSubnet] = c.Network.AciNetworkProvider.StaticExternalSubnet
		networkPluginConfigDefaultsMap[AciServiceGraphSubnet] = c.Network.AciNetworkProvider.ServiceGraphSubnet
		networkPluginConfigDefaultsMap[AciKubeAPIVlan] = c.Network.AciNetworkProvider.KubeAPIVlan
		networkPluginConfigDefaultsMap[AciServiceVlan] = c.Network.AciNetworkProvider.ServiceVlan
		networkPluginConfigDefaultsMap[AciInfraVlan] = c.Network.AciNetworkProvider.InfraVlan
		networkPluginConfigDefaultsMap[AciImagePullPolicy] = c.Network.AciNetworkProvider.ImagePullPolicy
		networkPluginConfigDefaultsMap[AciImagePullSecret] = c.Network.AciNetworkProvider.ImagePullSecret
		networkPluginConfigDefaultsMap[AciTenant] = c.Network.AciNetworkProvider.Tenant
		networkPluginConfigDefaultsMap[AciKafkaClientCrt] = c.Network.AciNetworkProvider.KafkaClientCrt
		networkPluginConfigDefaultsMap[AciKafkaClientKey] = c.Network.AciNetworkProvider.KafkaClientKey
		networkPluginConfigDefaultsMap[AciSubnetDomainName] = c.Network.AciNetworkProvider.SubnetDomainName
		networkPluginConfigDefaultsMap[AciEpRegistry] = c.Network.AciNetworkProvider.EpRegistry
		networkPluginConfigDefaultsMap[AciOpflexMode] = c.Network.AciNetworkProvider.OpflexMode
		networkPluginConfigDefaultsMap[AciOverlayVRFName] = c.Network.AciNetworkProvider.OverlayVRFName
		networkPluginConfigDefaultsMap[AciGbpPodSubnet] = c.Network.AciNetworkProvider.GbpPodSubnet
		networkPluginConfigDefaultsMap[AciOpflexServerPort] = c.Network.AciNetworkProvider.OpflexServerPort
	}
	if c.Network.KubeRouterNetworkProvider != nil {
		// by default, do not run service proxy in kube-router
		if c.Network.KubeRouterNetworkProvider.RunServiceProxy == nil {
			defaultKubeRouterRunServiceProxy := DefaultKubeRouterRunServiceProxy
			c.Network.KubeRouterNetworkProvider.RunServiceProxy = &defaultKubeRouterRunServiceProxy
		}
		// disable kube-proxy if a kube-router runs a service proxy
		if *c.Network.KubeRouterNetworkProvider.RunServiceProxy {
			serviceProxyEnabled := false
			c.Services.Kubeproxy.Enabled = &serviceProxyEnabled
		}
		// by default, do not run firewall
		if c.Network.KubeRouterNetworkProvider.RunFirewall == nil {
			defaultKubeRouterRunFirewall := DefaultKubeRouterRunFirewall
			c.Network.KubeRouterNetworkProvider.RunFirewall = &defaultKubeRouterRunFirewall
		}
	}
	for k, v := range networkPluginConfigDefaultsMap {
		setDefaultIfEmptyMapValue(c.Network.Options, k, v)
	}
}

func (c *Cluster) setClusterAuthnDefaults() {
	setDefaultIfEmpty(&c.Authentication.Strategy, DefaultAuthStrategy)

	for _, strategy := range strings.Split(c.Authentication.Strategy, "|") {
		strategy = strings.ToLower(strings.TrimSpace(strategy))
		c.AuthnStrategies[strategy] = true
	}

	if c.AuthnStrategies[AuthnWebhookProvider] && c.Authentication.Webhook == nil {
		c.Authentication.Webhook = &v3.AuthWebhookConfig{}
	}
	if c.Authentication.Webhook != nil {
		webhookConfigDefaultsMap := map[*string]string{
			&c.Authentication.Webhook.ConfigFile:   DefaultAuthnWebhookFile,
			&c.Authentication.Webhook.CacheTimeout: DefaultAuthnCacheTimeout,
		}
		for k, v := range webhookConfigDefaultsMap {
			setDefaultIfEmpty(k, v)
		}
	}
}

func d(image, defaultRegistryURL string) string {
	if len(defaultRegistryURL) == 0 {
		return image
	}
	return fmt.Sprintf("%s/%s", defaultRegistryURL, image)
}

func (c *Cluster) setCloudProvider() error {
	p, err := cloudprovider.InitCloudProvider(c.CloudProvider)
	if err != nil {
		return fmt.Errorf("Failed to initialize cloud provider: %v", err)
	}
	if p != nil {
		c.CloudConfigFile, err = p.GenerateCloudConfigFile()
		if err != nil {
			return fmt.Errorf("Failed to parse cloud config file: %v", err)
		}
		c.CloudProvider.Name = p.GetName()
		if c.CloudProvider.Name == "" {
			return fmt.Errorf("Name of the cloud provider is not defined for custom provider")
		}
	}
	return nil
}

func GetExternalFlags(local, updateOnly, disablePortCheck, useLocalState bool, configDir, clusterFilePath string) ExternalFlags {
	return ExternalFlags{
		Local:            local,
		UpdateOnly:       updateOnly,
		DisablePortCheck: disablePortCheck,
		ConfigDir:        configDir,
		ClusterFilePath:  clusterFilePath,
		UseLocalState:    useLocalState,
	}
}

func (c *Cluster) setAddonsDefaults() {
	c.Ingress.UpdateStrategy = setDaemonsetAddonDefaults(c.Ingress.UpdateStrategy)
	c.Network.UpdateStrategy = setDaemonsetAddonDefaults(c.Network.UpdateStrategy)
	c.DNS.UpdateStrategy = setDNSDeploymentAddonDefaults(c.DNS.UpdateStrategy, c.DNS.Provider)
	if c.DNS.LinearAutoscalerParams == nil {
		c.DNS.LinearAutoscalerParams = &DefaultClusterProportionalAutoscalerLinearParams
	}
	c.Monitoring.UpdateStrategy = setDeploymentAddonDefaults(c.Monitoring.UpdateStrategy)
	if c.Monitoring.Replicas == nil {
		c.Monitoring.Replicas = &DefaultMonitoringAddonReplicas
	}

	if c.Ingress.NetworkMode == "" {
		c.Ingress.NetworkMode = DefaultNetworkMode
	}

	if c.Ingress.HTTPPort == 0 {
		c.Ingress.HTTPPort = DefaultHTTPPort
	}

	if c.Ingress.HTTPSPort == 0 {
		c.Ingress.HTTPSPort = DefaultHTTPSPort
	}

	if c.Ingress.DefaultBackend == nil {
		defaultBackend := true
		k8sVersion := c.RancherKubernetesEngineConfig.Version
		toMatch, err := semver.Make(k8sVersion[1:])
		if err != nil {
			logrus.Warnf("Cluster version [%s] can not be parsed as semver", k8sVersion)
		}
		logrus.Debugf("Checking ingress default backend for cluster version [%s]", k8sVersion)
		// default backend will be false for k8s 1.21 and up
		disableIngressDefaultBackendRange, err := semver.ParseRange(">=1.21.0-rancher0")
		if err != nil {
			logrus.Warnf("Failed to parse semver range for checking ingress default backend")
		}
		if disableIngressDefaultBackendRange(toMatch) {
			logrus.Debugf("Cluster version [%s] needs to have ingress default backend disabled", k8sVersion)
			defaultBackend = false
		}
		c.Ingress.DefaultBackend = &defaultBackend
	}
}

func setDaemonsetAddonDefaults(updateStrategy *v3.DaemonSetUpdateStrategy) *v3.DaemonSetUpdateStrategy {
	if updateStrategy != nil && updateStrategy.Strategy != appsv1.RollingUpdateDaemonSetStrategyType {
		return updateStrategy
	}
	if updateStrategy == nil || updateStrategy.RollingUpdate == nil || updateStrategy.RollingUpdate.MaxUnavailable == nil {
		return &DefaultDaemonSetUpdateStrategy
	}
	return updateStrategy
}

func setDeploymentAddonDefaults(updateStrategy *v3.DeploymentStrategy) *v3.DeploymentStrategy {
	if updateStrategy != nil && updateStrategy.Strategy != appsv1.RollingUpdateDeploymentStrategyType {
		return updateStrategy
	}
	if updateStrategy == nil || updateStrategy.RollingUpdate == nil {
		return &DefaultDeploymentUpdateStrategy
	}
	if updateStrategy.RollingUpdate.MaxUnavailable == nil {
		updateStrategy.RollingUpdate.MaxUnavailable = &DefaultDeploymentUpdateStrategyParams
	}
	if updateStrategy.RollingUpdate.MaxSurge == nil {
		updateStrategy.RollingUpdate.MaxSurge = &DefaultDeploymentUpdateStrategyParams
	}
	return updateStrategy
}

func setDNSDeploymentAddonDefaults(updateStrategy *v3.DeploymentStrategy, dnsProvider string) *v3.DeploymentStrategy {
	var (
		coreDNSMaxUnavailable, coreDNSMaxSurge = intstr.FromInt(1), intstr.FromInt(0)
		kubeDNSMaxSurge, kubeDNSMaxUnavailable = intstr.FromString("10%"), intstr.FromInt(0)
	)
	if updateStrategy != nil && updateStrategy.Strategy != appsv1.RollingUpdateDeploymentStrategyType {
		return updateStrategy
	}
	switch dnsProvider {
	case CoreDNSProvider:
		if updateStrategy == nil || updateStrategy.RollingUpdate == nil {
			return &v3.DeploymentStrategy{
				Strategy: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &coreDNSMaxUnavailable,
					MaxSurge:       &coreDNSMaxSurge,
				},
			}
		}
		if updateStrategy.RollingUpdate.MaxUnavailable == nil {
			updateStrategy.RollingUpdate.MaxUnavailable = &coreDNSMaxUnavailable
		}
	case KubeDNSProvider:
		if updateStrategy == nil || updateStrategy.RollingUpdate == nil {
			return &v3.DeploymentStrategy{
				Strategy: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &kubeDNSMaxUnavailable,
					MaxSurge:       &kubeDNSMaxSurge,
				},
			}
		}
		if updateStrategy.RollingUpdate.MaxSurge == nil {
			updateStrategy.RollingUpdate.MaxSurge = &kubeDNSMaxSurge
		}
	}

	return updateStrategy
}

func checkVersionNeedsKubeAPIAuditLog(k8sVersion string) (bool, error) {
	toMatch, err := semver.Make(k8sVersion[1:])
	if err != nil {
		return false, fmt.Errorf("Cluster version [%s] can not be parsed as semver", k8sVersion[1:])
	}
	logrus.Debugf("Checking if cluster version [%s] needs to have kube-api audit log enabled", k8sVersion[1:])
	// kube-api audit log needs to be enabled for k8s 1.15.11 and up, k8s 1.16.8 and up, k8s 1.17.4 and up
	clusterKubeAPIAuditLogRange, err := semver.ParseRange(">=1.15.11-rancher0 <=1.15.99 || >=1.16.8-rancher0 <=1.16.99 || >=1.17.4-rancher0")
	if err != nil {
		return false, errors.New("Failed to parse semver range for checking to enable kube-api audit log")
	}
	if clusterKubeAPIAuditLogRange(toMatch) {
		logrus.Debugf("Cluster version [%s] needs to have kube-api audit log enabled", k8sVersion[1:])
		return true, nil
	}
	logrus.Debugf("Cluster version [%s] does not need to have kube-api audit log enabled", k8sVersion[1:])
	return false, nil
}
