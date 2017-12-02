package services

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/rancher/rke/docker"
	"github.com/rancher/rke/hosts"
	"github.com/rancher/rke/pki"
	"github.com/rancher/types/apis/cluster.cattle.io/v1"
)

func runKubeproxy(host *hosts.Host, kubeproxyService v1.KubeproxyService) error {
	imageCfg, hostCfg := buildKubeproxyConfig(host, kubeproxyService)
	return docker.DoRunContainer(host.DClient, imageCfg, hostCfg, KubeproxyContainerName, host.Address, WorkerRole)
}

func removeKubeproxy(host *hosts.Host) error {
	return docker.DoRemoveContainer(host.DClient, KubeproxyContainerName, host.Address)
}

func buildKubeproxyConfig(host *hosts.Host, kubeproxyService v1.KubeproxyService) (*container.Config, *container.HostConfig) {
	imageCfg := &container.Config{
		Image: kubeproxyService.Image,
		Entrypoint: []string{"kube-proxy",
			"--v=2",
			"--healthz-bind-address=0.0.0.0",
			"--kubeconfig=" + pki.KubeProxyConfigPath,
		},
	}
	hostCfg := &container.HostConfig{
		Binds: []string{
			"/etc/kubernetes:/etc/kubernetes",
		},
		NetworkMode:   "host",
		RestartPolicy: container.RestartPolicy{Name: "always"},
		Privileged:    true,
	}
	for arg, value := range kubeproxyService.ExtraArgs {
		cmd := fmt.Sprintf("--%s=%s", arg, value)
		imageCfg.Entrypoint = append(imageCfg.Entrypoint, cmd)
	}
	return imageCfg, hostCfg
}
