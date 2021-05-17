package cluster

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	cidr "github.com/apparentlymart/go-cidr/cidr"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/rancher/rke/docker"
	"github.com/rancher/rke/hosts"
	"github.com/rancher/rke/log"
	"github.com/rancher/rke/pki"
	"github.com/rancher/rke/templates"
	v3 "github.com/rancher/rke/types"
	"github.com/rancher/rke/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	appsv1 "k8s.io/api/apps/v1"
)

const (
	NetworkPluginResourceName = "rke-network-plugin"

	PortCheckContainer        = "rke-port-checker"
	EtcdPortListenContainer   = "rke-etcd-port-listener"
	CPPortListenContainer     = "rke-cp-port-listener"
	WorkerPortListenContainer = "rke-worker-port-listener"

	KubeAPIPort      = "6443"
	EtcdPort1        = "2379"
	EtcdPort2        = "2380"
	ScedulerPort     = "10251"
	ControllerPort   = "10252"
	KubeletPort      = "10250"
	KubeProxyPort    = "10256"
	FlannelVxLanPort = 8472

	FlannelVxLanNetworkIdentify = 1

	ProtocolTCP = "TCP"
	ProtocolUDP = "UDP"

	NoNetworkPlugin = "none"

	FlannelNetworkPlugin = "flannel"
	FlannelIface         = "flannel_iface"
	FlannelBackendType   = "flannel_backend_type"
	// FlannelBackendPort must be 4789 if using VxLan mode in the cluster with Windows nodes
	FlannelBackendPort = "flannel_backend_port"
	// FlannelBackendVxLanNetworkIdentify should be greater than or equal to 4096 if using VxLan mode in the cluster with Windows nodes
	FlannelBackendVxLanNetworkIdentify  = "flannel_backend_vni"
	KubeFlannelPriorityClassNameKeyName = "kube_flannel_priority_class_name"

	CalicoNetworkPlugin                           = "calico"
	CalicoNodeLabel                               = "calico-node"
	CalicoControllerLabel                         = "calico-kube-controllers"
	CalicoCloudProvider                           = "calico_cloud_provider"
	CalicoFlexVolPluginDirectory                  = "calico_flex_volume_plugin_dir"
	CalicoNodePriorityClassNameKeyName            = "calico_node_priority_class_name"
	CalicoKubeControllersPriorityClassNameKeyName = "calico_kube_controllers_priority_class_name"

	CanalNetworkPlugin      = "canal"
	CanalIface              = "canal_iface"
	CanalFlannelBackendType = "canal_flannel_backend_type"
	// CanalFlannelBackendPort must be 4789 if using Flannel VxLan mode in the cluster with Windows nodes
	CanalFlannelBackendPort = "canal_flannel_backend_port"
	// CanalFlannelBackendVxLanNetworkIdentify should be greater than or equal to 4096 if using Flannel VxLan mode in the cluster with Windows nodes
	CanalFlannelBackendVxLanNetworkIdentify = "canal_flannel_backend_vni"
	CanalFlexVolPluginDirectory             = "canal_flex_volume_plugin_dir"
	CanalPriorityClassNameKeyName           = "canal_priority_class_name"

	WeaveNetworkPlugin               = "weave"
	WeaveNetworkAppName              = "weave-net"
	WeaveNetPriorityClassNameKeyName = "weave_net_priority_class_name"

	AciNetworkPlugin            = "aci"
	AciOVSMemoryLimit           = "aci_ovs_memory_limit"
	AciImagePullPolicy          = "aci_image_pull_policy"
	AciPBRTrackingNonSnat       = "aci_pbr_tracking_non_snat"
	AciInstallIstio             = "aci_install_istio"
	AciIstioProfile             = "aci_istio_profile"
	AciDropLogEnable            = "aci_drop_log_enable"
	AciControllerLogLevel       = "aci_controller_log_level"
	AciHostAgentLogLevel        = "aci_host_agent_log_level"
	AciOpflexAgentLogLevel      = "aci_opflex_agent_log_level"
	AciApicRefreshTime          = "aci_apic_refresh_time"
	AciServiceMonitorInterval   = "aci_server_monitor_interval"
	AciSystemIdentifier         = "aci_system_identifier"
	AciToken                    = "aci_token"
	AciApicUserName             = "aci_apic_user_name"
	AciApicUserKey              = "aci_apic_user_key"
	AciApicUserCrt              = "aci_apic_user_crt"
	AciVmmDomain                = "aci_vmm_domain"
	AciVmmController            = "aci_vmm_controller"
	AciEncapType                = "aci_encap_type"
	AciAEP                      = "aci_aep"
	AciVRFName                  = "aci_vrf_name"
	AciVRFTenant                = "aci_vrf_tenant"
	AciL3Out                    = "aci_l3out"
	AciDynamicExternalSubnet    = "aci_dynamic_external_subnet"
	AciStaticExternalSubnet     = "aci_static_external_subnet"
	AciServiceGraphSubnet       = "aci_service_graph_subnet"
	AciKubeAPIVlan              = "aci_kubeapi_vlan"
	AciServiceVlan              = "aci_service_vlan"
	AciInfraVlan                = "aci_infra_vlan"
	AciImagePullSecret          = "aci_image_pull_secret"
	AciTenant                   = "aci_tenant"
	AciNodeSubnet               = "aci_node_subnet"
	AciMcastRangeStart          = "aci_mcast_range_start"
	AciMcastRangeEnd            = "aci_mcast_range_end"
	AciUseAciCniPriorityClass   = "aci_use_aci_cni_priority_class"
	AciNoPriorityClass          = "aci_no_priority_class"
	AciMaxNodesSvcGraph         = "aci_max_nodes_svc_graph"
	AciSnatContractScope        = "aci_snat_contract_scope"
	AciPodSubnetChunkSize       = "aci_pod_subnet_chunk_size"
	AciEnableEndpointSlice      = "aci_enable_endpoint_slice"
	AciSnatNamespace            = "aci_snat_namespace"
	AciEpRegistry               = "aci_ep_registry"
	AciOpflexMode               = "aci_opflex_mode"
	AciSnatPortRangeStart       = "aci_snat_port_range_start"
	AciSnatPortRangeEnd         = "aci_snat_port_range_end"
	AciSnatPortsPerNode         = "aci_snat_ports_per_node"
	AciOpflexClientSSL          = "aci_opflex_client_ssl"
	AciUsePrivilegedContainer   = "aci_use_privileged_container"
	AciUseHostNetnsVolume       = "aci_use_host_netns_volume"
	AciUseOpflexServerVolume    = "aci_use_opflex_server_volume"
	AciKafkaClientCrt           = "aci_kafka_client_crt"
	AciKafkaClientKey           = "aci_kafka_client_key"
	AciSubnetDomainName         = "aci_subnet_domain_name"
	AciCApic                    = "aci_capic"
	AciUseAciAnywhereCRD        = "aci_use_aci_anywhere_crd"
	AciOverlayVRFName           = "aci_overlay_vrf_name"
	AciGbpPodSubnet             = "aci_gbp_pod_subnet"
	AciRunGbpContainer          = "aci_run_gbp_container"
	AciRunOpflexServerContainer = "aci_run_opflex_server_container"
	AciOpflexServerPort         = "aci_opflex_server_port"
	AciDisableMultus            = "aci_disable_multus"
	// List of map keys to be used with network templates

	// EtcdEndpoints is the server address for Etcd, used by calico
	EtcdEndpoints = "EtcdEndpoints"
	// APIRoot is the kubernetes API address
	APIRoot = "APIRoot"
	// kubernetes client certificates and kubeconfig paths

	EtcdClientCert     = "EtcdClientCert"
	EtcdClientKey      = "EtcdClientKey"
	EtcdClientCA       = "EtcdClientCA"
	EtcdClientCertPath = "EtcdClientCertPath"
	EtcdClientKeyPath  = "EtcdClientKeyPath"
	EtcdClientCAPath   = "EtcdClientCAPath"

	ClientCertPath = "ClientCertPath"
	ClientKeyPath  = "ClientKeyPath"
	ClientCAPath   = "ClientCAPath"

	KubeCfg = "KubeCfg"

	ClusterCIDR = "ClusterCIDR"
	// Images key names

	Image              = "Image"
	CNIImage           = "CNIImage"
	NodeImage          = "NodeImage"
	ControllersImage   = "ControllersImage"
	CanalFlannelImg    = "CanalFlannelImg"
	FlexVolImg         = "FlexVolImg"
	WeaveLoopbackImage = "WeaveLoopbackImage"

	Calicoctl = "Calicoctl"

	FlannelInterface                       = "FlannelInterface"
	FlannelBackend                         = "FlannelBackend"
	KubeFlannelPriorityClassName           = "KubeFlannelPriorityClassName"
	CalicoNodePriorityClassName            = "CalicoNodePriorityClassName"
	CalicoKubeControllersPriorityClassName = "CalicoKubeControllersPriorityClassName"
	CanalInterface                         = "CanalInterface"
	CanalPriorityClassName                 = "CanalPriorityClassName"
	FlexVolPluginDir                       = "FlexVolPluginDir"
	WeavePassword                          = "WeavePassword"
	WeaveNetPriorityClassName              = "WeaveNetPriorityClassName"
	MTU                                    = "MTU"
	RBACConfig                             = "RBACConfig"
	ClusterVersion                         = "ClusterVersion"
	SystemIdentifier                       = "SystemIdentifier"
	ApicHosts                              = "ApicHosts"
	Token                                  = "Token"
	ApicUserName                           = "ApicUserName"
	ApicUserKey                            = "ApicUserKey"
	ApicUserCrt                            = "ApicUserCrt"
	ApicRefreshTime                        = "ApicRefreshTime"
	VmmDomain                              = "VmmDomain"
	VmmController                          = "VmmController"
	EncapType                              = "EncapType"
	McastRangeStart                        = "McastRangeStart"
	McastRangeEnd                          = "McastRangeEnd"
	AEP                                    = "AEP"
	VRFName                                = "VRFName"
	VRFTenant                              = "VRFTenant"
	L3Out                                  = "L3Out"
	L3OutExternalNetworks                  = "L3OutExternalNetworks"
	DynamicExternalSubnet                  = "DynamicExternalSubnet"
	StaticExternalSubnet                   = "StaticExternalSubnet"
	ServiceGraphSubnet                     = "ServiceGraphSubnet"
	KubeAPIVlan                            = "KubeAPIVlan"
	ServiceVlan                            = "ServiceVlan"
	InfraVlan                              = "InfraVlan"
	ImagePullPolicy                        = "ImagePullPolicy"
	ImagePullSecret                        = "ImagePullSecret"
	Tenant                                 = "Tenant"
	ServiceMonitorInterval                 = "ServiceMonitorInterval"
	PBRTrackingNonSnat                     = "PBRTrackingNonSnat"
	InstallIstio                           = "InstallIstio"
	IstioProfile                           = "IstioProfile"
	DropLogEnable                          = "DropLogEnable"
	ControllerLogLevel                     = "ControllerLogLevel"
	HostAgentLogLevel                      = "HostAgentLogLevel"
	OpflexAgentLogLevel                    = "OpflexAgentLogLevel"
	AciCniDeployContainer                  = "AciCniDeployContainer"
	AciHostContainer                       = "AciHostContainer"
	AciOpflexContainer                     = "AciOpflexContainer"
	AciMcastContainer                      = "AciMcastContainer"
	AciOpenvSwitchContainer                = "AciOpenvSwitchContainer"
	AciControllerContainer                 = "AciControllerContainer"
	AciGbpServerContainer                  = "AciGbpServerContainer"
	AciOpflexServerContainer               = "AciOpflexServerContainer"
	StaticServiceIPStart                   = "StaticServiceIPStart"
	StaticServiceIPEnd                     = "StaticServiceIPEnd"
	PodGateway                             = "PodGateway"
	PodIPStart                             = "PodIPStart"
	PodIPEnd                               = "PodIPEnd"
	NodeServiceIPStart                     = "NodeServiceIPStart"
	NodeServiceIPEnd                       = "NodeServiceIPEnd"
	ServiceIPStart                         = "ServiceIPStart"
	ServiceIPEnd                           = "ServiceIPEnd"
	UseAciCniPriorityClass                 = "UseAciCniPriorityClass"
	NoPriorityClass                        = "NoPriorityClass"
	MaxNodesSvcGraph                       = "MaxNodesSvcGraph"
	SnatContractScope                      = "SnatContractScope"
	PodSubnetChunkSize                     = "PodSubnetChunkSize"
	EnableEndpointSlice                    = "EnableEndpointSlice"
	SnatNamespace                          = "SnatNamespace"
	EpRegistry                             = "EpRegistry"
	OpflexMode                             = "OpflexMode"
	SnatPortRangeStart                     = "SnatPortRangeStart"
	SnatPortRangeEnd                       = "SnatPortRangeEnd"
	SnatPortsPerNode                       = "SnatPortsPerNode"
	OpflexClientSSL                        = "OpflexClientSSL"
	UsePrivilegedContainer                 = "UsePrivilegedContainer"
	UseHostNetnsVolume                     = "UseHostNetnsVolume"
	UseOpflexServerVolume                  = "UseOpflexServerVolume"
	KafkaBrokers                           = "KafkaBrokers"
	KafkaClientCrt                         = "KafkaClientCrt"
	KafkaClientKey                         = "KafkaClientKey"
	SubnetDomainName                       = "SubnetDomainName"
	CApic                                  = "CApic"
	UseAciAnywhereCRD                      = "UseAciAnywhereCRD"
	OverlayVRFName                         = "OverlayVRFName"
	GbpPodSubnet                           = "GbpPodSubnet"
	RunGbpContainer                        = "RunGbpContainer"
	RunOpflexServerContainer               = "RunOpflexServerContainer"
	OpflexServerPort                       = "OpflexServerPort"
	OVSMemoryLimit                         = "OVSMemoryLimit"
	NodeSubnet                             = "NodeSubnet"
	NodeSelector                           = "NodeSelector"
	UpdateStrategy                         = "UpdateStrategy"
	Tolerations                            = "Tolerations"
	DisableMultus                          = "DisableMultus"
)

var EtcdPortList = []string{
	EtcdPort1,
	EtcdPort2,
}

var ControlPlanePortList = []string{
	KubeAPIPort,
}

var WorkerPortList = []string{
	KubeletPort,
}

var EtcdClientPortList = []string{
	EtcdPort1,
}

var CalicoNetworkLabels = []string{CalicoNodeLabel, CalicoControllerLabel}

func (c *Cluster) deployNetworkPlugin(ctx context.Context, data map[string]interface{}) error {
	log.Infof(ctx, "[network] Setting up network plugin: %s", c.Network.Plugin)
	switch c.Network.Plugin {
	case FlannelNetworkPlugin:
		return c.doFlannelDeploy(ctx, data)
	case CalicoNetworkPlugin:
		return c.doCalicoDeploy(ctx, data)
	case CanalNetworkPlugin:
		return c.doCanalDeploy(ctx, data)
	case WeaveNetworkPlugin:
		return c.doWeaveDeploy(ctx, data)
	case AciNetworkPlugin:
		return c.doAciDeploy(ctx, data)
	case NoNetworkPlugin:
		log.Infof(ctx, "[network] Not deploying a cluster network, expecting custom CNI")
		return nil
	default:
		return fmt.Errorf("[network] Unsupported network plugin: %s", c.Network.Plugin)
	}
}

func (c *Cluster) doFlannelDeploy(ctx context.Context, data map[string]interface{}) error {
	vni, err := atoiWithDefault(c.Network.Options[FlannelBackendVxLanNetworkIdentify], FlannelVxLanNetworkIdentify)
	if err != nil {
		return err
	}
	port, err := atoiWithDefault(c.Network.Options[FlannelBackendPort], FlannelVxLanPort)
	if err != nil {
		return err
	}

	flannelConfig := map[string]interface{}{
		ClusterCIDR:      c.ClusterCIDR,
		Image:            c.SystemImages.Flannel,
		CNIImage:         c.SystemImages.FlannelCNI,
		FlannelInterface: c.Network.Options[FlannelIface],
		FlannelBackend: map[string]interface{}{
			"Type": c.Network.Options[FlannelBackendType],
			"VNI":  vni,
			"Port": port,
		},
		RBACConfig:     c.Authorization.Mode,
		ClusterVersion: util.GetTagMajorVersion(c.Version),
		NodeSelector:   c.Network.NodeSelector,
		UpdateStrategy: &appsv1.DaemonSetUpdateStrategy{
			Type:          c.Network.UpdateStrategy.Strategy,
			RollingUpdate: c.Network.UpdateStrategy.RollingUpdate,
		},
		KubeFlannelPriorityClassName: c.Network.Options[KubeFlannelPriorityClassNameKeyName],
	}
	pluginYaml, err := c.getNetworkPluginManifest(flannelConfig, data)
	if err != nil {
		return err
	}
	return c.doAddonDeploy(ctx, pluginYaml, NetworkPluginResourceName, true)
}

func (c *Cluster) doCalicoDeploy(ctx context.Context, data map[string]interface{}) error {
	clientConfig := pki.GetConfigPath(pki.KubeNodeCertName)

	calicoConfig := map[string]interface{}{
		KubeCfg:          clientConfig,
		ClusterCIDR:      c.ClusterCIDR,
		CNIImage:         c.SystemImages.CalicoCNI,
		NodeImage:        c.SystemImages.CalicoNode,
		Calicoctl:        c.SystemImages.CalicoCtl,
		ControllersImage: c.SystemImages.CalicoControllers,
		CloudProvider:    c.Network.Options[CalicoCloudProvider],
		FlexVolImg:       c.SystemImages.CalicoFlexVol,
		RBACConfig:       c.Authorization.Mode,
		NodeSelector:     c.Network.NodeSelector,
		MTU:              c.Network.MTU,
		UpdateStrategy: &appsv1.DaemonSetUpdateStrategy{
			Type:          c.Network.UpdateStrategy.Strategy,
			RollingUpdate: c.Network.UpdateStrategy.RollingUpdate,
		},
		Tolerations:                            c.Network.Tolerations,
		FlexVolPluginDir:                       c.Network.Options[CalicoFlexVolPluginDirectory],
		CalicoNodePriorityClassName:            c.Network.Options[CalicoNodePriorityClassNameKeyName],
		CalicoKubeControllersPriorityClassName: c.Network.Options[CalicoKubeControllersPriorityClassNameKeyName],
	}
	pluginYaml, err := c.getNetworkPluginManifest(calicoConfig, data)
	if err != nil {
		return err
	}
	return c.doAddonDeploy(ctx, pluginYaml, NetworkPluginResourceName, true)
}

func (c *Cluster) doCanalDeploy(ctx context.Context, data map[string]interface{}) error {
	flannelVni, err := atoiWithDefault(c.Network.Options[CanalFlannelBackendVxLanNetworkIdentify], FlannelVxLanNetworkIdentify)
	if err != nil {
		return err
	}
	flannelPort, err := atoiWithDefault(c.Network.Options[CanalFlannelBackendPort], FlannelVxLanPort)
	if err != nil {
		return err
	}

	clientConfig := pki.GetConfigPath(pki.KubeNodeCertName)
	canalConfig := map[string]interface{}{
		ClientCertPath:   pki.GetCertPath(pki.KubeNodeCertName),
		APIRoot:          "https://127.0.0.1:6443",
		ClientKeyPath:    pki.GetKeyPath(pki.KubeNodeCertName),
		ClientCAPath:     pki.GetCertPath(pki.CACertName),
		KubeCfg:          clientConfig,
		ClusterCIDR:      c.ClusterCIDR,
		NodeImage:        c.SystemImages.CanalNode,
		CNIImage:         c.SystemImages.CanalCNI,
		ControllersImage: c.SystemImages.CanalControllers,
		CanalFlannelImg:  c.SystemImages.CanalFlannel,
		RBACConfig:       c.Authorization.Mode,
		CanalInterface:   c.Network.Options[CanalIface],
		FlexVolImg:       c.SystemImages.CanalFlexVol,
		FlannelBackend: map[string]interface{}{
			"Type": c.Network.Options[CanalFlannelBackendType],
			"VNI":  flannelVni,
			"Port": flannelPort,
		},
		NodeSelector: c.Network.NodeSelector,
		MTU:          c.Network.MTU,
		UpdateStrategy: &appsv1.DaemonSetUpdateStrategy{
			Type:          c.Network.UpdateStrategy.Strategy,
			RollingUpdate: c.Network.UpdateStrategy.RollingUpdate,
		},
		Tolerations:                            c.Network.Tolerations,
		FlexVolPluginDir:                       c.Network.Options[CanalFlexVolPluginDirectory],
		CanalPriorityClassName:                 c.Network.Options[CanalPriorityClassNameKeyName],
		CalicoKubeControllersPriorityClassName: c.Network.Options[CalicoKubeControllersPriorityClassNameKeyName],
	}
	pluginYaml, err := c.getNetworkPluginManifest(canalConfig, data)
	if err != nil {
		return err
	}
	return c.doAddonDeploy(ctx, pluginYaml, NetworkPluginResourceName, true)
}

func (c *Cluster) doWeaveDeploy(ctx context.Context, data map[string]interface{}) error {
	weaveConfig := map[string]interface{}{
		ClusterCIDR:        c.ClusterCIDR,
		WeavePassword:      c.Network.Options[WeavePassword],
		Image:              c.SystemImages.WeaveNode,
		CNIImage:           c.SystemImages.WeaveCNI,
		WeaveLoopbackImage: c.SystemImages.Alpine,
		RBACConfig:         c.Authorization.Mode,
		NodeSelector:       c.Network.NodeSelector,
		MTU:                c.Network.MTU,
		UpdateStrategy: &appsv1.DaemonSetUpdateStrategy{
			Type:          c.Network.UpdateStrategy.Strategy,
			RollingUpdate: c.Network.UpdateStrategy.RollingUpdate,
		},
		WeaveNetPriorityClassName: c.Network.Options[WeaveNetPriorityClassNameKeyName],
	}
	pluginYaml, err := c.getNetworkPluginManifest(weaveConfig, data)
	if err != nil {
		return err
	}
	return c.doAddonDeploy(ctx, pluginYaml, NetworkPluginResourceName, true)
}

func (c *Cluster) doAciDeploy(ctx context.Context, data map[string]interface{}) error {
	_, clusterCIDR, err := net.ParseCIDR(c.ClusterCIDR)
	if err != nil {
		return err
	}
	podIPStart, podIPEnd := cidr.AddressRange(clusterCIDR)
	_, staticExternalSubnet, err := net.ParseCIDR(c.Network.Options[AciStaticExternalSubnet])
	staticServiceIPStart, staticServiceIPEnd := cidr.AddressRange(staticExternalSubnet)
	_, svcGraphSubnet, err := net.ParseCIDR(c.Network.Options[AciServiceGraphSubnet])
	if err != nil {
		return err
	}
	nodeServiceIPStart, nodeServiceIPEnd := cidr.AddressRange(svcGraphSubnet)
	_, dynamicExternalSubnet, err := net.ParseCIDR(c.Network.Options[AciDynamicExternalSubnet])
	if err != nil {
		return err
	}
	serviceIPStart, serviceIPEnd := cidr.AddressRange(dynamicExternalSubnet)
	if c.Network.Options[AciTenant] == "" {
		c.Network.Options[AciTenant] = c.Network.Options[AciSystemIdentifier]
	}

	AciConfig := map[string]interface{}{
		SystemIdentifier:         c.Network.Options[AciSystemIdentifier],
		ApicHosts:                c.Network.AciNetworkProvider.ApicHosts,
		Token:                    c.Network.Options[AciToken],
		ApicUserName:             c.Network.Options[AciApicUserName],
		ApicUserKey:              c.Network.Options[AciApicUserKey],
		ApicUserCrt:              c.Network.Options[AciApicUserCrt],
		ApicRefreshTime:          c.Network.Options[AciApicRefreshTime],
		VmmDomain:                c.Network.Options[AciVmmDomain],
		VmmController:            c.Network.Options[AciVmmController],
		EncapType:                c.Network.Options[AciEncapType],
		McastRangeStart:          c.Network.Options[AciMcastRangeStart],
		McastRangeEnd:            c.Network.Options[AciMcastRangeEnd],
		NodeSubnet:               c.Network.Options[AciNodeSubnet],
		AEP:                      c.Network.Options[AciAEP],
		VRFName:                  c.Network.Options[AciVRFName],
		VRFTenant:                c.Network.Options[AciVRFTenant],
		L3Out:                    c.Network.Options[AciL3Out],
		L3OutExternalNetworks:    c.Network.AciNetworkProvider.L3OutExternalNetworks,
		DynamicExternalSubnet:    c.Network.Options[AciDynamicExternalSubnet],
		StaticExternalSubnet:     c.Network.Options[AciStaticExternalSubnet],
		ServiceGraphSubnet:       c.Network.Options[AciServiceGraphSubnet],
		KubeAPIVlan:              c.Network.Options[AciKubeAPIVlan],
		ServiceVlan:              c.Network.Options[AciServiceVlan],
		InfraVlan:                c.Network.Options[AciInfraVlan],
		ImagePullPolicy:          c.Network.Options[AciImagePullPolicy],
		ImagePullSecret:          c.Network.Options[AciImagePullSecret],
		Tenant:                   c.Network.Options[AciTenant],
		ServiceMonitorInterval:   c.Network.Options[AciServiceMonitorInterval],
		PBRTrackingNonSnat:       c.Network.Options[AciPBRTrackingNonSnat],
		InstallIstio:             c.Network.Options[AciInstallIstio],
		IstioProfile:             c.Network.Options[AciIstioProfile],
		DropLogEnable:            c.Network.Options[AciDropLogEnable],
		ControllerLogLevel:       c.Network.Options[AciControllerLogLevel],
		HostAgentLogLevel:        c.Network.Options[AciHostAgentLogLevel],
		OpflexAgentLogLevel:      c.Network.Options[AciOpflexAgentLogLevel],
		OVSMemoryLimit:           c.Network.Options[AciOVSMemoryLimit],
		ClusterCIDR:              c.ClusterCIDR,
		StaticServiceIPStart:     cidr.Inc(cidr.Inc(staticServiceIPStart)),
		StaticServiceIPEnd:       cidr.Dec(staticServiceIPEnd),
		PodGateway:               cidr.Inc(podIPStart),
		PodIPStart:               cidr.Inc(cidr.Inc(podIPStart)),
		PodIPEnd:                 cidr.Dec(podIPEnd),
		NodeServiceIPStart:       cidr.Inc(cidr.Inc(nodeServiceIPStart)),
		NodeServiceIPEnd:         cidr.Dec(nodeServiceIPEnd),
		ServiceIPStart:           cidr.Inc(cidr.Inc(serviceIPStart)),
		ServiceIPEnd:             cidr.Dec(serviceIPEnd),
		UseAciCniPriorityClass:   c.Network.Options[AciUseAciCniPriorityClass],
		NoPriorityClass:          c.Network.Options[AciNoPriorityClass],
		MaxNodesSvcGraph:         c.Network.Options[AciMaxNodesSvcGraph],
		SnatContractScope:        c.Network.Options[AciSnatContractScope],
		PodSubnetChunkSize:       c.Network.Options[AciPodSubnetChunkSize],
		EnableEndpointSlice:      c.Network.Options[AciEnableEndpointSlice],
		SnatNamespace:            c.Network.Options[AciSnatNamespace],
		EpRegistry:               c.Network.Options[AciEpRegistry],
		OpflexMode:               c.Network.Options[AciOpflexMode],
		SnatPortRangeStart:       c.Network.Options[AciSnatPortRangeStart],
		SnatPortRangeEnd:         c.Network.Options[AciSnatPortRangeEnd],
		SnatPortsPerNode:         c.Network.Options[AciSnatPortsPerNode],
		OpflexClientSSL:          c.Network.Options[AciOpflexClientSSL],
		UsePrivilegedContainer:   c.Network.Options[AciUsePrivilegedContainer],
		UseHostNetnsVolume:       c.Network.Options[AciUseHostNetnsVolume],
		UseOpflexServerVolume:    c.Network.Options[AciUseOpflexServerVolume],
		KafkaBrokers:             c.Network.AciNetworkProvider.KafkaBrokers,
		KafkaClientCrt:           c.Network.Options[AciKafkaClientCrt],
		KafkaClientKey:           c.Network.Options[AciKafkaClientKey],
		SubnetDomainName:         c.Network.Options[AciSubnetDomainName],
		CApic:                    c.Network.Options[AciCApic],
		UseAciAnywhereCRD:        c.Network.Options[AciUseAciAnywhereCRD],
		OverlayVRFName:           c.Network.Options[AciOverlayVRFName],
		GbpPodSubnet:             c.Network.Options[AciGbpPodSubnet],
		RunGbpContainer:          c.Network.Options[AciRunGbpContainer],
		RunOpflexServerContainer: c.Network.Options[AciRunOpflexServerContainer],
		OpflexServerPort:         c.Network.Options[AciOpflexServerPort],
		DisableMultus:            c.Network.Options[AciDisableMultus],
		AciCniDeployContainer:    c.SystemImages.AciCniDeployContainer,
		AciHostContainer:         c.SystemImages.AciHostContainer,
		AciOpflexContainer:       c.SystemImages.AciOpflexContainer,
		AciMcastContainer:        c.SystemImages.AciMcastContainer,
		AciOpenvSwitchContainer:  c.SystemImages.AciOpenvSwitchContainer,
		AciControllerContainer:   c.SystemImages.AciControllerContainer,
		AciGbpServerContainer:    c.SystemImages.AciGbpServerContainer,
		AciOpflexServerContainer: c.SystemImages.AciOpflexServerContainer,
		MTU:                      c.Network.MTU,
	}

	pluginYaml, err := c.getNetworkPluginManifest(AciConfig, data)
	if err != nil {
		return err
	}
	return c.doAddonDeploy(ctx, pluginYaml, NetworkPluginResourceName, true)
}

func (c *Cluster) getNetworkPluginManifest(pluginConfig, data map[string]interface{}) (string, error) {
	switch c.Network.Plugin {
	case CanalNetworkPlugin, FlannelNetworkPlugin, CalicoNetworkPlugin, WeaveNetworkPlugin, AciNetworkPlugin:
		tmplt, err := templates.GetVersionedTemplates(c.Network.Plugin, data, c.Version)
		if err != nil {
			return "", err
		}
		return templates.CompileTemplateFromMap(tmplt, pluginConfig)
	default:
		return "", fmt.Errorf("[network] Unsupported network plugin: %s", c.Network.Plugin)
	}
}

func (c *Cluster) CheckClusterPorts(ctx context.Context, currentCluster *Cluster) error {
	if currentCluster != nil {
		newEtcdHost := hosts.GetToAddHosts(currentCluster.EtcdHosts, c.EtcdHosts)
		newControlPlaneHosts := hosts.GetToAddHosts(currentCluster.ControlPlaneHosts, c.ControlPlaneHosts)
		newWorkerHosts := hosts.GetToAddHosts(currentCluster.WorkerHosts, c.WorkerHosts)

		if len(newEtcdHost) == 0 &&
			len(newWorkerHosts) == 0 &&
			len(newControlPlaneHosts) == 0 {
			log.Infof(ctx, "[network] No hosts added existing cluster, skipping port check")
			return nil
		}
	}
	if err := c.deployTCPPortListeners(ctx, currentCluster); err != nil {
		return err
	}
	if err := c.runServicePortChecks(ctx); err != nil {
		return err
	}
	// Skip kubeapi check if we are using custom k8s dialer or bastion/jump host
	if c.K8sWrapTransport == nil && len(c.BastionHost.Address) == 0 {
		if err := c.checkKubeAPIPort(ctx); err != nil {
			return err
		}
	} else {
		log.Infof(ctx, "[network] Skipping kubeapi port check")
	}

	return c.removeTCPPortListeners(ctx)
}

func (c *Cluster) checkKubeAPIPort(ctx context.Context) error {
	log.Infof(ctx, "[network] Checking KubeAPI port Control Plane hosts")
	for _, host := range c.ControlPlaneHosts {
		logrus.Debugf("[network] Checking KubeAPI port [%s] on host: %s", KubeAPIPort, host.Address)
		address := fmt.Sprintf("%s:%s", host.Address, KubeAPIPort)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return fmt.Errorf("[network] Can't access KubeAPI port [%s] on Control Plane host: %s", KubeAPIPort, host.Address)
		}
		conn.Close()
	}
	return nil
}

func (c *Cluster) deployTCPPortListeners(ctx context.Context, currentCluster *Cluster) error {
	log.Infof(ctx, "[network] Deploying port listener containers")

	// deploy ectd listeners
	if err := c.deployListenerOnPlane(ctx, EtcdPortList, c.EtcdHosts, EtcdPortListenContainer); err != nil {
		return err
	}

	// deploy controlplane listeners
	if err := c.deployListenerOnPlane(ctx, ControlPlanePortList, c.ControlPlaneHosts, CPPortListenContainer); err != nil {
		return err
	}

	// deploy worker listeners
	if err := c.deployListenerOnPlane(ctx, WorkerPortList, c.WorkerHosts, WorkerPortListenContainer); err != nil {
		return err
	}
	log.Infof(ctx, "[network] Port listener containers deployed successfully")
	return nil
}

func (c *Cluster) deployListenerOnPlane(ctx context.Context, portList []string, hostPlane []*hosts.Host, containerName string) error {
	var errgrp errgroup.Group
	hostsQueue := util.GetObjectQueue(hostPlane)
	for w := 0; w < WorkerThreads; w++ {
		errgrp.Go(func() error {
			var errList []error
			for host := range hostsQueue {
				err := c.deployListener(ctx, host.(*hosts.Host), portList, containerName)
				if err != nil {
					errList = append(errList, err)
				}
			}
			return util.ErrList(errList)
		})
	}
	return errgrp.Wait()
}

func (c *Cluster) deployListener(ctx context.Context, host *hosts.Host, portList []string, containerName string) error {
	imageCfg := &container.Config{
		Image: c.SystemImages.Alpine,
		Cmd: []string{
			"nc",
			"-kl",
			"-p",
			"1337",
			"-e",
			"echo",
		},
		ExposedPorts: nat.PortSet{
			"1337/tcp": {},
		},
	}
	hostCfg := &container.HostConfig{
		PortBindings: nat.PortMap{
			"1337/tcp": getPortBindings("0.0.0.0", portList),
		},
	}

	logrus.Debugf("[network] Starting deployListener [%s] on host [%s]", containerName, host.Address)
	if err := docker.DoRunContainer(ctx, host.DClient, imageCfg, hostCfg, containerName, host.Address, "network", c.PrivateRegistriesMap); err != nil {
		if strings.Contains(err.Error(), "bind: address already in use") {
			logrus.Debugf("[network] Service is already up on host [%s]", host.Address)
			return nil
		}
		return err
	}
	return nil
}

func (c *Cluster) removeTCPPortListeners(ctx context.Context) error {
	log.Infof(ctx, "[network] Removing port listener containers")

	if err := removeListenerFromPlane(ctx, c.EtcdHosts, EtcdPortListenContainer); err != nil {
		return err
	}
	if err := removeListenerFromPlane(ctx, c.ControlPlaneHosts, CPPortListenContainer); err != nil {
		return err
	}
	if err := removeListenerFromPlane(ctx, c.WorkerHosts, WorkerPortListenContainer); err != nil {
		return err
	}
	log.Infof(ctx, "[network] Port listener containers removed successfully")
	return nil
}

func removeListenerFromPlane(ctx context.Context, hostPlane []*hosts.Host, containerName string) error {
	var errgrp errgroup.Group

	hostsQueue := util.GetObjectQueue(hostPlane)
	for w := 0; w < WorkerThreads; w++ {
		errgrp.Go(func() error {
			var errList []error
			for host := range hostsQueue {
				runHost := host.(*hosts.Host)
				err := docker.DoRemoveContainer(ctx, runHost.DClient, containerName, runHost.Address)
				if err != nil {
					errList = append(errList, err)
				}
			}
			return util.ErrList(errList)
		})
	}
	return errgrp.Wait()
}

func (c *Cluster) runServicePortChecks(ctx context.Context) error {
	var errgrp errgroup.Group
	// check etcd <-> etcd
	// one etcd host is a pass
	if len(c.EtcdHosts) > 1 {
		log.Infof(ctx, "[network] Running etcd <-> etcd port checks")
		hostsQueue := util.GetObjectQueue(c.EtcdHosts)
		for w := 0; w < WorkerThreads; w++ {
			errgrp.Go(func() error {
				var errList []error
				for host := range hostsQueue {
					err := checkPlaneTCPPortsFromHost(ctx, host.(*hosts.Host), EtcdPortList, c.EtcdHosts, c.SystemImages.Alpine, c.PrivateRegistriesMap)
					if err != nil {
						errList = append(errList, err)
					}
				}
				return util.ErrList(errList)
			})
		}
		if err := errgrp.Wait(); err != nil {
			return err
		}
	}
	// check control -> etcd connectivity
	log.Infof(ctx, "[network] Running control plane -> etcd port checks")
	hostsQueue := util.GetObjectQueue(c.ControlPlaneHosts)
	for w := 0; w < WorkerThreads; w++ {
		errgrp.Go(func() error {
			var errList []error
			for host := range hostsQueue {
				err := checkPlaneTCPPortsFromHost(ctx, host.(*hosts.Host), EtcdClientPortList, c.EtcdHosts, c.SystemImages.Alpine, c.PrivateRegistriesMap)
				if err != nil {
					errList = append(errList, err)
				}
			}
			return util.ErrList(errList)
		})
	}
	if err := errgrp.Wait(); err != nil {
		return err
	}
	// check controle plane -> Workers
	log.Infof(ctx, "[network] Running control plane -> worker port checks")
	hostsQueue = util.GetObjectQueue(c.ControlPlaneHosts)
	for w := 0; w < WorkerThreads; w++ {
		errgrp.Go(func() error {
			var errList []error
			for host := range hostsQueue {
				err := checkPlaneTCPPortsFromHost(ctx, host.(*hosts.Host), WorkerPortList, c.WorkerHosts, c.SystemImages.Alpine, c.PrivateRegistriesMap)
				if err != nil {
					errList = append(errList, err)
				}
			}
			return util.ErrList(errList)
		})
	}
	if err := errgrp.Wait(); err != nil {
		return err
	}
	// check workers -> control plane
	log.Infof(ctx, "[network] Running workers -> control plane port checks")
	hostsQueue = util.GetObjectQueue(c.WorkerHosts)
	for w := 0; w < WorkerThreads; w++ {
		errgrp.Go(func() error {
			var errList []error
			for host := range hostsQueue {
				err := checkPlaneTCPPortsFromHost(ctx, host.(*hosts.Host), ControlPlanePortList, c.ControlPlaneHosts, c.SystemImages.Alpine, c.PrivateRegistriesMap)
				if err != nil {
					errList = append(errList, err)
				}
			}
			return util.ErrList(errList)
		})
	}
	return errgrp.Wait()
}

func checkPlaneTCPPortsFromHost(ctx context.Context, host *hosts.Host, portList []string, planeHosts []*hosts.Host, image string, prsMap map[string]v3.PrivateRegistry) error {
	var hosts []string
	var portCheckLogs string
	for _, host := range planeHosts {
		hosts = append(hosts, host.InternalAddress)
	}
	imageCfg := &container.Config{
		Image: image,
		Env: []string{
			fmt.Sprintf("HOSTS=%s", strings.Join(hosts, " ")),
			fmt.Sprintf("PORTS=%s", strings.Join(portList, " ")),
		},
		Cmd: []string{
			"sh",
			"-c",
			"for host in $HOSTS; do for port in $PORTS ; do echo \"Checking host ${host} on port ${port}\" >&1 & nc -w 5 -z $host $port > /dev/null || echo \"${host}:${port}\" >&2 & done; wait; done",
		},
	}
	hostCfg := &container.HostConfig{
		NetworkMode: "host",
		LogConfig: container.LogConfig{
			Type: "json-file",
		},
	}
	for retries := 0; retries < 3; retries++ {
		logrus.Infof("[network] Checking if host [%s] can connect to host(s) [%s] on port(s) [%s], try #%d", host.Address, strings.Join(hosts, " "), strings.Join(portList, " "), retries+1)
		if err := docker.DoRemoveContainer(ctx, host.DClient, PortCheckContainer, host.Address); err != nil {
			return err
		}
		if err := docker.DoRunContainer(ctx, host.DClient, imageCfg, hostCfg, PortCheckContainer, host.Address, "network", prsMap); err != nil {
			return err
		}

		containerLog, _, logsErr := docker.GetContainerLogsStdoutStderr(ctx, host.DClient, PortCheckContainer, "all", true)
		if logsErr != nil {
			log.Warnf(ctx, "[network] Failed to get network port check logs: %v", logsErr)
		}
		logrus.Debugf("[network] containerLog [%s] on host: %s", containerLog, host.Address)

		if err := docker.RemoveContainer(ctx, host.DClient, host.Address, PortCheckContainer); err != nil {
			return err
		}
		logrus.Debugf("[network] Length of containerLog is [%d] on host: %s", len(containerLog), host.Address)
		if len(containerLog) == 0 {
			return nil
		}
		portCheckLogs = strings.Join(strings.Split(strings.TrimSpace(containerLog), "\n"), ", ")
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("[network] Host [%s] is not able to connect to the following ports: [%s]. Please check network policies and firewall rules", host.Address, portCheckLogs)
}

func getPortBindings(hostAddress string, portList []string) []nat.PortBinding {
	portBindingList := []nat.PortBinding{}
	for _, portNumber := range portList {
		rawPort := fmt.Sprintf("%s:%s:1337/tcp", hostAddress, portNumber)
		portMapping, _ := nat.ParsePortSpec(rawPort)
		portBindingList = append(portBindingList, portMapping[0].Binding)
	}
	return portBindingList
}

func atoiWithDefault(val string, defaultVal int) (int, error) {
	if val == "" {
		return defaultVal, nil
	}

	ret, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return ret, nil
}
