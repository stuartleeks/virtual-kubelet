package localdocker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/virtual-kubelet/virtual-kubelet/providers"

	"github.com/virtual-kubelet/virtual-kubelet/manager"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/remotecommand"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	dockernetwork "github.com/docker/docker/api/types/network"
	dockerclient "github.com/docker/docker/client"
)

//
// This is a simple sample provider implementation
// It isn't anticipated that it solves any particularly valuable scenario
// It is a long way from ready - no namespace or network handling as two obvious examples ;-)
//

// Provider implements the virtual-kubelet provider and executes against the local docker context
type Provider struct {
	resourceManager    *manager.ResourceManager
	nodeName           string
	internalIP         string
	daemonEndpointPort int32
	dockerClient       *dockerclient.Client
	pods               []*podInfo
}
type podInfo struct {
	pod         *v1.Pod
	containerID string
}

// NewLocalDockerProvider creates a new Provider instance
func NewLocalDockerProvider(resourcemanager *manager.ResourceManager, nodeName string, internalIP string, daemonEndpointPort int32) (*Provider, error) {
	if nodeName == "" {
		return nil, fmt.Errorf("nodeName is required")
	}
	dockerClient, err := dockerclient.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to create docker client: %v", err)
	}
	pods := make([]*podInfo, 0)
	provider := Provider{
		resourceManager:    resourcemanager,
		nodeName:           nodeName,
		internalIP:         internalIP,
		daemonEndpointPort: daemonEndpointPort,
		dockerClient:       dockerClient,
		pods:               pods,
	}
	return &provider, nil
}
func createDockerContainerName(podName string, containerName string) string { // TODO - add namespace!
	return fmt.Sprintf("VK_%s_%s", podName, containerName)
}

// CreatePod takes a Kubernetes Pod and deploys it within the provider.
func (p *Provider) CreatePod(pod *v1.Pod) error {

	// Currently only handling a single container, for simplicity
	if len(pod.Spec.Containers) != 1 {
		return fmt.Errorf("CreatePod currently only supports a single container per pod")
	}
	containerSpec := pod.Spec.Containers[0]

	config := dockercontainer.Config{
		Image: containerSpec.Image,
	}

	containerName := createDockerContainerName(pod.Name, containerSpec.Name)

	// TODO add in ImagePull

	pod.Status.Phase = v1.PodPending
	pod.Status.Message = "Creating"
	log.Printf("Creating container %s\n", containerName)
	// TODO handle exposing ports
	container, err := p.dockerClient.ContainerCreate(context.Background(), &config, &dockercontainer.HostConfig{}, &dockernetwork.NetworkingConfig{}, containerName)
	if err != nil {
		return fmt.Errorf("ContainerCreate failed: %v", err)
	}
	log.Printf("Created container %s. ID: %s\n", containerName, container.ID)

	pod.Status.Message = "Starting"
	log.Printf("Starting container %s\n", container.ID)
	err = p.dockerClient.ContainerStart(context.Background(), container.ID, dockertypes.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("ContainerStart failed: %v", err)
	}
	// TODO does ContainerStart wait for container to start before returning?
	log.Printf("Started container %s\n", container.ID)

	pod.Status.Phase = v1.PodRunning
	pod.Status.Message = "Running"

	now := metav1.NewTime(time.Now())
	pod.Status.StartTime = &now

	// containerInfo, err := p.dockerClient.ContainerInspect(context.Background(), container.ID)
	// if err != nil {

	// 	return err
	// }
	// pod.Status.HostIP = "192.168.1.181"
	// pod.Status.PodIP = containerInfo.NetworkSettings.IPAddress

	// pod.Status.HostIP = "1.2.3.4" // TODO
	// pod.Status.PodIP = "5.6.7.8"  // TODO
	pod.Status.Conditions = []v1.PodCondition{
		{
			Type:   v1.PodInitialized,
			Status: v1.ConditionTrue,
		},
		{
			Type:   v1.PodReady,
			Status: v1.ConditionTrue,
		},
		{
			Type:   v1.PodScheduled,
			Status: v1.ConditionTrue,
		},
	}

	pod.Status.ContainerStatuses = append(pod.Status.ContainerStatuses, v1.ContainerStatus{
		Name:         containerSpec.Name,
		Image:        containerSpec.Image,
		Ready:        true,
		RestartCount: 0,
		State: v1.ContainerState{
			Running: &v1.ContainerStateRunning{
				StartedAt: now,
			},
		},
	})

	podInfo := podInfo{
		pod:         pod,
		containerID: container.ID,
	}
	p.pods = append(p.pods, &podInfo)
	return nil
}

// UpdatePod takes a Kubernetes Pod and updates it within the provider.
func (p *Provider) UpdatePod(pod *v1.Pod) error {
	return fmt.Errorf("not implemented: UpdatePod")
}

// DeletePod takes a Kubernetes Pod and deletes it from the provider.
func (p *Provider) DeletePod(pod *v1.Pod) error {
	// Currently only handling a single container, for simplicity
	if len(pod.Spec.Containers) != 1 {
		return fmt.Errorf("DeletePod currently only supports a single container per pod")
	}

	for _, podInfo := range p.pods {
		if podInfo.pod.Namespace == pod.Namespace && podInfo.pod.Name == pod.Name {
			err := p.dockerClient.ContainerRemove(context.Background(), podInfo.containerID, dockertypes.ContainerRemoveOptions{Force: true})
			if err != nil {
				return fmt.Errorf("ContainerRemove failed: %v", err)
			}
			// TODO - do we need to update the pod status here?
		}
	}
	return fmt.Errorf("DeletePod: pod not found: %s:%s", pod.Namespace, pod.Name)
}

// GetPod retrieves a pod by name from the provider (can be cached).
func (p *Provider) GetPod(namespace, name string) (*v1.Pod, error) {
	log.Printf("GetPod called for %s:%s\n", namespace, name)
	for _, podInfo := range p.pods {
		if podInfo.pod.Namespace == namespace && podInfo.pod.Name == name {
			return podInfo.pod, nil
		}
	}
	return nil, fmt.Errorf("GetPod: pod not found %s:%s", namespace, name)
}

// GetPodStatus retrievesthe status of a pod by name from the provider.
func (p *Provider) GetPodStatus(namespace, name string) (*v1.PodStatus, error) {
	log.Printf("GetPodStatus called for %s:%s\n", namespace, name)
	for _, podInfo := range p.pods {
		if podInfo.pod.Namespace == namespace && podInfo.pod.Name == name {
			// TODO - check that the container is running!
			return &podInfo.pod.Status, nil
		}
	}
	return nil, fmt.Errorf("GetPodStatus: pod not found %s:%s", namespace, name)
}

// GetPods retrieves a list of all pods running on the provider (can be cached).
func (p *Provider) GetPods() ([]*v1.Pod, error) {
	log.Printf("GetPods called\n")
	pods := make([]*v1.Pod, len(p.pods))
	for index, pod := range p.pods {
		pods[index] = pod.pod
	}
	return pods, nil
}

// Capacity returns a resource list with the capacity constraints of the provider.
func (p *Provider) Capacity() v1.ResourceList {
	return v1.ResourceList{
		"cpu":    resource.MustParse("20"),
		"memory": resource.MustParse("100Gi"),
		"pod":    resource.MustParse("20"),
	}
}

// NodeConditions returns a list of conditions (Ready, OutOfDisk, etc), which is polled periodically to update the node status
// within Kubernetes.
func (p *Provider) NodeConditions() []v1.NodeCondition {
	// TODO Currently always reporting healthy - consider checking daemon is running etc
	return []v1.NodeCondition{
		{
			Type:               "Ready",
			Status:             v1.ConditionTrue,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletReady",
			Message:            "Currently optimistic about local docker daemon state",
		},
		{
			Type:               "OutOfDisk",
			Status:             v1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletHasSufficientDisk",
			Message:            "kubelet has sufficient disk space available",
		},
		{
			Type:               "MemoryPressure",
			Status:             v1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletHasSufficientMemory",
			Message:            "kubelet has sufficient memory available",
		},
		{
			Type:               "DiskPressure",
			Status:             v1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "KubeletHasNoDiskPressure",
			Message:            "kubelet has no disk pressure",
		},
		{
			Type:               "NetworkUnavailable",
			Status:             v1.ConditionFalse,
			LastHeartbeatTime:  metav1.Now(),
			LastTransitionTime: metav1.Now(),
			Reason:             "RouteCreated",
			Message:            "RouteController created a route",
		},
	}
}

// OperatingSystem returns the operating system the provider is for.
func (p *Provider) OperatingSystem() string {
	return providers.OperatingSystemLinux // just linux for now
}

// ExecInContainer executes a command in a container in the pod, copying data
// between in/out/err and the container's stdin/stdout/stderr.
func (p *Provider) ExecInContainer(name string, uid types.UID, container string, cmd []string, in io.Reader, out, err io.WriteCloser, tty bool, resize <-chan remotecommand.TerminalSize, timeout time.Duration) error {
	return fmt.Errorf("not implemented: ExecInContainer")
}

// GetContainerLogs retrieves the logs of a container by name from the provider.
func (p *Provider) GetContainerLogs(namespace, podName, containerName string, tail int) (string, error) {
	log.Printf("GetContainerLogs called for %s:%s:%s\n", namespace, podName, containerName)

	for _, podInfo := range p.pods {
		if podInfo.pod.Namespace == namespace && podInfo.pod.Name == podName {
			tailString := fmt.Sprintf("%d", tail)
			readerCloser, err := p.dockerClient.ContainerLogs(context.Background(), podInfo.containerID, dockertypes.ContainerLogsOptions{Tail: tailString})

			if err != nil {
				return "", fmt.Errorf("ContainerLogs failed: %v", err)
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(readerCloser)
			s := buf.String()

			readerCloser.Close()

			return s, nil
		}
	}
	return "", fmt.Errorf("Not found: %s:%s:%s", namespace, podName, containerName)
}

// NodeAddresses returns a list of addresses for the node status
// within Kubernetes.
func (p *Provider) NodeAddresses() []v1.NodeAddress {
	// return nil
	return []v1.NodeAddress{
		{
			Type:    "InternalIP",
			Address: p.internalIP,
		},
	}
}

// NodeDaemonEndpoints returns NodeDaemonEndpoints for the node status
// within Kubernetes.
func (p *Provider) NodeDaemonEndpoints() *v1.NodeDaemonEndpoints {
	return &v1.NodeDaemonEndpoints{
		KubeletEndpoint: v1.DaemonEndpoint{
			Port: p.daemonEndpointPort,
		},
	}
}
