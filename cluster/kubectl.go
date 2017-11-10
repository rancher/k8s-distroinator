package cluster

import (
	"context"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/rancher/rke/docker"
	"github.com/rancher/rke/pki"
)

const (
	KubectlImage    = "melsayed/kubectl:latest"
	KubctlContainer = "kubectl"
)

type KubectlCommand struct {
	Cmd []string
	Env []string
}

func (c *Cluster) buildClusterConfigEnv() []string {
	// This needs to be updated when add more configuration
	return []string{
		pki.ConvertConfigToENV(pki.KubeAdminConfigENVName, c.Certificates[pki.KubeAdminCommonName].Config),
		pki.ConvertConfigToENV(ClusterCIDREnvName, c.ClusterCIDR),
		pki.ConvertConfigToENV(ClusterDNSServerIPEnvName, c.ClusterDNSServer),
		pki.ConvertConfigToENV(ClusterDomainEnvName, c.ClusterDomain),
	}
}

func (c *Cluster) RunKubectlCmd(kubectlCmd *KubectlCommand) error {
	h := c.ControlPlaneHosts[0]

	logrus.Debugf("[kubectl] Using host [%s] for deployment", h.Hostname)
	logrus.Debugf("[kubectl] Pulling kubectl image..")

	if err := docker.PullImage(h.DClient, h.Hostname, KubectlImage); err != nil {
		return err
	}

	clusterConfigEnv := c.buildClusterConfigEnv()
	if kubectlCmd.Env != nil {
		clusterConfigEnv = append(clusterConfigEnv, kubectlCmd.Env...)
	}

	imageCfg := &container.Config{
		Image: KubectlImage,
		Env:   clusterConfigEnv,
		Cmd:   kubectlCmd.Cmd,
	}
	logrus.Debugf("[kubectl] Creating kubectl container..")
	resp, err := h.DClient.ContainerCreate(context.Background(), imageCfg, nil, nil, KubctlContainer)
	if err != nil {
		return fmt.Errorf("Failed to create kubectl container on host [%s]: %v", h.Hostname, err)
	}
	logrus.Debugf("[kubectl] Container %s created..", resp.ID)
	if err := h.DClient.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("Failed to start kubectl container on host [%s]: %v", h.Hostname, err)
	}
	logrus.Debugf("[kubectl] running command: %s", kubectlCmd.Cmd)
	statusCh, errCh := h.DClient.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("Failed to execute kubectl container on host [%s]: %v", h.Hostname, err)
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return fmt.Errorf("kubectl command failed on host [%s]: exit status %v", h.Hostname, status.StatusCode)
		}
	}
	if err := h.DClient.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{}); err != nil {
		return fmt.Errorf("Failed to remove kubectl container on host[%s]: %v", h.Hostname, err)
	}
	return nil
}
