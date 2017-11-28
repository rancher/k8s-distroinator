package services

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/rancher/rke/docker"
	"github.com/rancher/rke/hosts"
	"github.com/rancher/rke/pki"
	"github.com/rancher/types/apis/cluster.cattle.io/v1"
)

func runScheduler(host hosts.Host, schedulerService v1.SchedulerService) error {
	imageCfg, hostCfg := buildSchedulerConfig(host, schedulerService)
	return docker.DoRunContainer(host.DClient, imageCfg, hostCfg, SchedulerContainerName, host.AdvertisedHostname, ControlRole)
}

func removeScheduler(host hosts.Host) error {
	return docker.DoRemoveContainer(host.DClient, SchedulerContainerName, host.AdvertisedHostname)
}

func buildSchedulerConfig(host hosts.Host, schedulerService v1.SchedulerService) (*container.Config, *container.HostConfig) {
	imageCfg := &container.Config{
		Image: schedulerService.Image,
		Entrypoint: []string{"kube-scheduler",
			"--leader-elect=true",
			"--v=2",
			"--address=0.0.0.0",
			"--kubeconfig=" + pki.KubeSchedulerConfigPath,
		},
	}
	hostCfg := &container.HostConfig{
		Binds: []string{
			"/etc/kubernetes:/etc/kubernetes",
		},
		NetworkMode:   "host",
		RestartPolicy: container.RestartPolicy{Name: "always"},
	}
	for arg, value := range schedulerService.ExtraArgs {
		cmd := fmt.Sprintf("--%s=%s", arg, value)
		imageCfg.Cmd = append(imageCfg.Cmd, cmd)
	}
	return imageCfg, hostCfg
}
