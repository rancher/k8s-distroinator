package services

import (
	"context"
	"path"

	"github.com/rancher/rke/docker"
	"github.com/rancher/rke/hosts"
	"github.com/rancher/rke/pki"
	"github.com/rancher/types/apis/management.cattle.io/v3"
)

const (
	KubeletDockerConfigPath = "/var/lib/kubelet/config.json"
)

func runKubelet(ctx context.Context, host *hosts.Host, df hosts.DialerFactory, prsMap map[string]v3.PrivateRegistry, kubeletProcess v3.Process, certMap map[string]pki.CertificatePKI, alpineImage string) error {
	imageCfg, hostCfg, healthCheckURL := GetProcessConfig(kubeletProcess)
	if err := docker.DoRunContainer(ctx, host.DClient, imageCfg, hostCfg, KubeletContainerName, host.Address, WorkerRole, prsMap); err != nil {
		return err
	}
	// we have private registries, so we write docker config for kubelet
	if len(prsMap) > 0 {
		dockerConfig, err := docker.GetKubeletDockerConfig(prsMap)
		if err != nil {
			return err
		}
		configPath := path.Join(host.PrefixPath, KubeletDockerConfigPath)
		if err := docker.WriteFileToContainer(ctx, host.DClient, host.Address, KubeletContainerName, configPath, dockerConfig); err != nil {
			return err
		}
	}
	if err := runHealthcheck(ctx, host, KubeletContainerName, df, healthCheckURL, certMap); err != nil {
		return err
	}
	return createLogLink(ctx, host, KubeletContainerName, WorkerRole, alpineImage, prsMap)
}

func removeKubelet(ctx context.Context, host *hosts.Host) error {
	return docker.DoRemoveContainer(ctx, host.DClient, KubeletContainerName, host.Address)
}
