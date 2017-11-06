package cluster

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/rke/hosts"
	"github.com/rancher/rke/k8s"
	"github.com/rancher/rke/pki"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

func (c *Cluster) SaveClusterState(clusterFile string) error {
	// Reinitialize kubernetes Client
	var err error
	c.KClient, err = k8s.NewClient(pki.KubeAdminConfigPath)
	if err != nil {
		return fmt.Errorf("Failed to re-initialize Kubernetes Client: %v", err)
	}
	err = saveClusterCerts(c.KClient, c.Certificates)
	if err != nil {
		return fmt.Errorf("[certificates] Failed to Save Kubernetes certificates: %v", err)
	}
	err = saveStateToKubernetes(c.KClient, pki.KubeAdminConfigPath, []byte(clusterFile))
	if err != nil {
		return fmt.Errorf("[state] Failed to save configuration state: %v", err)
	}
	return nil
}

func (c *Cluster) GetClusterState() (*Cluster, error) {
	var err error
	var currentCluster *Cluster
	c.KClient, err = k8s.NewClient(pki.KubeAdminConfigPath)
	if err != nil {
		logrus.Warnf("Failed to initiate new Kubernetes Client: %v", err)
	} else {
		// Handle pervious kubernetes state and certificate generation
		currentCluster = getStateFromKubernetes(c.KClient, pki.KubeAdminConfigPath)
		if currentCluster != nil {
			err = currentCluster.InvertIndexHosts()
			if err != nil {
				return nil, fmt.Errorf("Failed to classify hosts from fetched cluster: %v", err)
			}
			err = hosts.ReconcileWorkers(currentCluster.WorkerHosts, c.WorkerHosts, c.KClient)
			if err != nil {
				return nil, fmt.Errorf("Failed to reconcile hosts: %v", err)
			}
		}
	}
	return currentCluster, nil
}

func saveStateToKubernetes(kClient *kubernetes.Clientset, kubeConfigPath string, clusterFile []byte) error {
	logrus.Infof("[state] Saving cluster state to Kubernetes")
	timeout := make(chan bool, 1)
	go func() {
		for {
			err := k8s.UpdateConfigMap(kClient, clusterFile, StateConfigMapName)
			if err != nil {
				time.Sleep(time.Second * 5)
				continue
			}
			logrus.Infof("[state] Successfully Saved cluster state to Kubernetes ConfigMap: %s", StateConfigMapName)
			timeout <- true
			break
		}
	}()
	select {
	case <-timeout:
		return nil
	case <-time.After(time.Second * UpdateStateTimeout):
		return fmt.Errorf("[state] Timeout waiting for kubernetes to be ready")
	}
}

func getStateFromKubernetes(kClient *kubernetes.Clientset, kubeConfigPath string) *Cluster {
	logrus.Infof("[state] Fetching cluster state from Kubernetes")
	var cfgMap *v1.ConfigMap
	var currentCluster Cluster
	var err error
	timeout := make(chan bool, 1)
	go func() {
		for {
			cfgMap, err = k8s.GetConfigMap(kClient, StateConfigMapName)
			if err != nil {
				time.Sleep(time.Second * 5)
				continue
			}
			logrus.Infof("[state] Successfully Fetched cluster state to Kubernetes ConfigMap: %s", StateConfigMapName)
			timeout <- true
			break
		}
	}()
	select {
	case <-timeout:
		clusterData := cfgMap.Data[StateConfigMapName]
		err := yaml.Unmarshal([]byte(clusterData), &currentCluster)
		if err != nil {
			return nil
		}
		return &currentCluster
	case <-time.After(time.Second * GetStateTimeout):
		logrus.Warnf("Timed out waiting for kubernetes cluster")
		return nil
	}
	return nil
}
