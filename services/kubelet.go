package services

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/rancher/rke/docker"
	"github.com/rancher/rke/hosts"
	"github.com/rancher/rke/pki"
	"github.com/rancher/types/apis/management.cattle.io/v3"
)

func runKubelet(host *hosts.Host, kubeletService v3.KubeletService) error {
	imageCfg, hostCfg := buildKubeletConfig(host, kubeletService)
	return docker.DoRunContainer(host.DClient, imageCfg, hostCfg, KubeletContainerName, host.Address, WorkerRole)
}

func removeKubelet(host *hosts.Host) error {
	return docker.DoRemoveContainer(host.DClient, KubeletContainerName, host.Address)
}

func buildKubeletConfig(host *hosts.Host, kubeletService v3.KubeletService) (*container.Config, *container.HostConfig) {
	imageCfg := &container.Config{
		Image: kubeletService.Image,
		Entrypoint: []string{"kubelet",
			"--v=2",
			"--address=0.0.0.0",
			"--cluster-domain=" + kubeletService.ClusterDomain,
			"--pod-infra-container-image=" + kubeletService.InfraContainerImage,
			"--cgroup-driver=cgroupfs",
			"--cgroups-per-qos=True",
			"--enforce-node-allocatable=",
			"--hostname-override=" + host.HostnameOverride,
			"--cluster-dns=" + kubeletService.ClusterDNSServer,
			"--network-plugin=cni",
			"--cni-conf-dir=/etc/cni/net.d",
			"--cni-bin-dir=/opt/cni/bin",
			"--resolv-conf=/etc/resolv.conf",
			"--allow-privileged=true",
			"--cloud-provider=",
			"--kubeconfig=" + pki.KubeNodeConfigPath,
			"--require-kubeconfig=True",
		},
	}
	for _, role := range host.Role {
		switch role {
		case ETCDRole:
			imageCfg.Cmd = append(imageCfg.Cmd, "--node-labels=node-role.kubernetes.io/etcd=true")
		case ControlRole:
			imageCfg.Cmd = append(imageCfg.Cmd, "--node-labels=node-role.kubernetes.io/master=true")
		case WorkerRole:
			imageCfg.Cmd = append(imageCfg.Cmd, "--node-labels=node-role.kubernetes.io/worker=true")
		}
	}
	hostCfg := &container.HostConfig{
		Binds: []string{
			"/etc/kubernetes:/etc/kubernetes",
			"/etc/cni:/etc/cni:ro",
			"/opt/cni:/opt/cni:ro",
			"/etc/resolv.conf:/etc/resolv.conf",
			"/sys:/sys:ro",
			"/var/lib/docker:/var/lib/docker:rw",
			"/var/lib/kubelet:/var/lib/kubelet:shared",
			"/var/run:/var/run:rw",
			"/run:/run",
			"/dev:/host/dev"},
		NetworkMode:   "host",
		PidMode:       "host",
		Privileged:    true,
		RestartPolicy: container.RestartPolicy{Name: "always"},
		PortBindings: nat.PortMap{
			"8080/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "8080",
				},
			},
		},
	}
	for arg, value := range kubeletService.ExtraArgs {
		cmd := fmt.Sprintf("--%s=%s", arg, value)
		imageCfg.Entrypoint = append(imageCfg.Entrypoint, cmd)
	}
	return imageCfg, hostCfg
}
