package docker

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

func DoRunContainer(dClient *client.Client, imageCfg *container.Config, hostCfg *container.HostConfig, containerName string, hostname string, plane string) error {
	isRunning, err := IsContainerRunning(dClient, hostname, containerName)
	if err != nil {
		return err
	}
	if isRunning {
		logrus.Infof("[%s] Container %s is already running on host [%s]", plane, containerName, hostname)
		return nil
	}
	logrus.Debugf("[%s] Pulling Image on host [%s]", plane, hostname)
	err = PullImage(dClient, hostname, imageCfg.Image)
	if err != nil {
		return err
	}
	logrus.Infof("[%s] Successfully pulled %s image on host [%s]", plane, containerName, hostname)
	resp, err := dClient.ContainerCreate(context.Background(), imageCfg, hostCfg, nil, containerName)
	if err != nil {
		return fmt.Errorf("Failed to create %s container on host [%s]: %v", containerName, hostname, err)
	}
	if err := dClient.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("Failed to start %s container on host [%s]: %v", containerName, hostname, err)
	}
	logrus.Debugf("[%s] Successfully started %s container: %s", plane, containerName, resp.ID)
	logrus.Infof("[%s] Successfully started %s container on host [%s]", plane, containerName, hostname)
	return nil
}

func DoRollingUpdateContainer(dClient *client.Client, imageCfg *container.Config, hostCfg *container.HostConfig, containerName, hostname, plane string) error {
	logrus.Debugf("[%s] Checking for deployed %s", plane, containerName)
	isRunning, err := IsContainerRunning(dClient, hostname, containerName)
	if err != nil {
		return err
	}
	if !isRunning {
		logrus.Infof("[%s] Container %s is not running on host [%s]", plane, containerName, hostname)
		return nil
	}
	logrus.Debugf("[%s] Stopping old container", plane)
	oldContainerName := "old-" + containerName
	if err := StopRenameContainer(dClient, hostname, containerName, oldContainerName); err != nil {
		return err
	}
	logrus.Infof("[%s] Successfully stopped old container %s on host [%s]", plane, containerName, hostname)
	_, err = CreateContiner(dClient, hostname, containerName, imageCfg, hostCfg)
	if err != nil {
		return fmt.Errorf("Failed to create %s container on host [%s]: %v", containerName, hostname, err)
	}
	if err := StartContainer(dClient, hostname, containerName); err != nil {
		return fmt.Errorf("Failed to start %s container on host [%s]: %v", containerName, hostname, err)
	}
	logrus.Infof("[%s] Successfully updated %s container on host [%s]", plane, containerName, hostname)
	logrus.Debugf("[%s] Removing old container", plane)
	err = RemoveContainer(dClient, hostname, oldContainerName)
	return err
}

func IsContainerRunning(dClient *client.Client, hostname string, containerName string) (bool, error) {
	logrus.Debugf("Checking if container %s is running on host [%s]", containerName, hostname)
	containers, err := dClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return false, fmt.Errorf("Can't get Docker containers for host [%s]: %v", hostname, err)

	}
	for _, container := range containers {
		if container.Names[0] == "/"+containerName {
			return true, nil
		}
	}
	return false, nil
}

func PullImage(dClient *client.Client, hostname string, containerImage string) error {
	out, err := dClient.ImagePull(context.Background(), containerImage, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("Can't pull Docker image %s for host [%s]: %v", containerImage, hostname, err)
	}
	defer out.Close()
	if logrus.GetLevel() == logrus.DebugLevel {
		io.Copy(os.Stdout, out)
	} else {
		io.Copy(ioutil.Discard, out)
	}

	return nil
}

func RemoveContainer(dClient *client.Client, hostname string, containerName string) error {
	err := dClient.ContainerRemove(context.Background(), containerName, types.ContainerRemoveOptions{})
	if err != nil {
		return fmt.Errorf("Can't remove Docker container %s for host [%s]: %v", containerName, hostname, err)
	}
	return nil
}

func StopContainer(dClient *client.Client, hostname string, containerName string) error {
	err := dClient.ContainerStop(context.Background(), containerName, nil)
	if err != nil {
		return fmt.Errorf("Can't stop Docker container %s for host [%s]: %v", containerName, hostname, err)
	}
	return nil
}

func RenameContainer(dClient *client.Client, hostname string, oldContainerName string, newContainerName string) error {
	err := dClient.ContainerRename(context.Background(), oldContainerName, newContainerName)
	if err != nil {
		return fmt.Errorf("Can't rename Docker container %s for host [%s]: %v", oldContainerName, hostname, err)
	}
	return nil
}

func StartContainer(dClient *client.Client, hostname string, containerName string) error {
	if err := dClient.ContainerStart(context.Background(), containerName, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("Failed to start %s container on host [%s]: %v", containerName, hostname, err)
	}
	return nil
}

func CreateContiner(dClient *client.Client, hostname string, containerName string, imageCfg *container.Config, hostCfg *container.HostConfig) (container.ContainerCreateCreatedBody, error) {
	created, err := dClient.ContainerCreate(context.Background(), imageCfg, hostCfg, nil, containerName)
	if err != nil {
		return container.ContainerCreateCreatedBody{}, fmt.Errorf("Failed to create %s container on host [%s]: %v", containerName, hostname, err)
	}
	return created, nil
}

func InspectContainer(dClient *client.Client, hostname string, containerName string) (types.ContainerJSON, error) {
	inspection, err := dClient.ContainerInspect(context.Background(), containerName)
	if err != nil {
		return types.ContainerJSON{}, fmt.Errorf("Failed to inspect %s container on host [%s]: %v", containerName, hostname, err)
	}
	return inspection, nil
}

func StopRenameContainer(dClient *client.Client, hostname string, oldContainerName string, newContainerName string) error {
	if err := StopContainer(dClient, hostname, oldContainerName); err != nil {
		return err
	}
	if err := WaitForContainer(dClient, oldContainerName); err != nil {
		return nil
	}
	err := RenameContainer(dClient, hostname, oldContainerName, newContainerName)
	return err
}

func WaitForContainer(dClient *client.Client, containerName string) error {
	statusCh, errCh := dClient.ContainerWait(context.Background(), containerName, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("Error wating for container [%s]: %v", containerName, err)
		}
	case <-statusCh:
	}
	return nil
}
