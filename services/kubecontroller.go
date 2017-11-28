package services

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/rancher/rke/docker"
	"github.com/rancher/rke/hosts"
	"github.com/rancher/rke/pki"
	"github.com/rancher/types/apis/cluster.cattle.io/v1"
)

func runKubeController(host hosts.Host, kubeControllerService v1.KubeControllerService) error {
	imageCfg, hostCfg := buildKubeControllerConfig(kubeControllerService)
	return docker.DoRunContainer(host.DClient, imageCfg, hostCfg, KubeControllerContainerName, host.AdvertisedHostname, ControlRole)
}

func removeKubeController(host hosts.Host) error {
	return docker.DoRemoveContainer(host.DClient, KubeControllerContainerName, host.AdvertisedHostname)
}

func buildKubeControllerConfig(kubeControllerService v1.KubeControllerService) (*container.Config, *container.HostConfig) {
	imageCfg := &container.Config{
		Image: kubeControllerService.Image,
		Entrypoint: []string{"kube-controller-manager",
			"--address=0.0.0.0",
			"--cloud-provider=",
			"--leader-elect=true",
			"--kubeconfig=" + pki.KubeControllerConfigPath,
			"--enable-hostpath-provisioner=false",
			"--node-monitor-grace-period=40s",
			"--pod-eviction-timeout=5m0s",
			"--v=2",
			"--allocate-node-cidrs=true",
			"--cluster-cidr=" + kubeControllerService.ClusterCIDR,
			"--service-cluster-ip-range=" + kubeControllerService.ServiceClusterIPRange,
			"--service-account-private-key-file=" + pki.KubeAPIKeyPath,
			"--root-ca-file=" + pki.CACertPath,
		},
	}
	hostCfg := &container.HostConfig{
		Binds: []string{
			"/etc/kubernetes:/etc/kubernetes",
		},
		NetworkMode:   "host",
		RestartPolicy: container.RestartPolicy{Name: "always"},
	}
	for arg, value := range kubeControllerService.ExtraArgs {
		cmd := fmt.Sprintf("--%s=%s", arg, value)
		imageCfg.Cmd = append(imageCfg.Cmd, cmd)
	}
	return imageCfg, hostCfg
}
