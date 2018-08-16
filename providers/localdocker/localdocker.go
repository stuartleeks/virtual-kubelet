package localdocker

import (
	"fmt"
	"io"
	"time"

	"github.com/virtual-kubelet/virtual-kubelet/providers"

	"github.com/virtual-kubelet/virtual-kubelet/manager"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/remotecommand"
)

// LocalDockerProvider implements the virtual-kubelet provider and executes against the local docker context
type LocalDockerProvider struct {
	resourceManager    *manager.ResourceManager
	nodeName           string
	internalIP         string
	daemonEndpointPort int32
}

// NewLocalDockerProvider creates a new LocalDockerProvider instance
func NewLocalDockerProvider(resourcemanager *manager.ResourceManager, nodeName string, internalIP string, daemonEndpointPort int32) (*LocalDockerProvider, error) {
	if nodeName == "" {
		return nil, fmt.Errorf("nodeName is required")
	}
	provider := LocalDockerProvider{
		resourceManager:    resourcemanager,
		nodeName:           nodeName,
		internalIP:         internalIP,
		daemonEndpointPort: daemonEndpointPort,
	}
	return &provider, nil
}

// CreatePod takes a Kubernetes Pod and deploys it within the provider.
func (p *LocalDockerProvider) CreatePod(pod *v1.Pod) error {

	return fmt.Errorf("not implemented: CreatePod")
	// // Currently only handling a single container
	// if len(pod.Spec.Containers) != 1 {

	// }
	// container := pod.Spec.Containers
	// return nil
}

// UpdatePod takes a Kubernetes Pod and updates it within the provider.
func (p *LocalDockerProvider) UpdatePod(pod *v1.Pod) error {
	return fmt.Errorf("not implemented: UpdatePod")
}

// DeletePod takes a Kubernetes Pod and deletes it from the provider.
func (p *LocalDockerProvider) DeletePod(pod *v1.Pod) error {
	return fmt.Errorf("not implemented: DeletePod")
}

// GetPod retrieves a pod by name from the provider (can be cached).
func (p *LocalDockerProvider) GetPod(namespace, name string) (*v1.Pod, error) {
	return nil, fmt.Errorf("not implemented: GetPod")
}

// GetPodStatus retrievesthe status of a pod by name from the provider.
func (p *LocalDockerProvider) GetPodStatus(namespace, name string) (*v1.PodStatus, error) {
	return nil, fmt.Errorf("not implemented: GetPodStatus")
}

// GetPods retrieves a list of all pods running on the provider (can be cached).
func (p *LocalDockerProvider) GetPods() ([]*v1.Pod, error) {
	fmt.Printf("TODO: GetPods - stubbed to return empty array\n")
	var pods []*v1.Pod
	return pods, nil
}

// Capacity returns a resource list with the capacity constraints of the provider.
func (p *LocalDockerProvider) Capacity() v1.ResourceList {
	return v1.ResourceList{
		"cpu":    resource.MustParse("20"),
		"memory": resource.MustParse("100Gi"),
		"pod":    resource.MustParse("20"),
	}
}

// NodeConditions returns a list of conditions (Ready, OutOfDisk, etc), which is polled periodically to update the node status
// within Kubernetes.
func (p *LocalDockerProvider) NodeConditions() []v1.NodeCondition {
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
func (p *LocalDockerProvider) OperatingSystem() string {
	return providers.OperatingSystemLinux // just linux for now
}

// ExecInContainer executes a command in a container in the pod, copying data
// between in/out/err and the container's stdin/stdout/stderr.
func (p *LocalDockerProvider) ExecInContainer(name string, uid types.UID, container string, cmd []string, in io.Reader, out, err io.WriteCloser, tty bool, resize <-chan remotecommand.TerminalSize, timeout time.Duration) error {
	return fmt.Errorf("not implemented: ExecInContainer")
}

// GetContainerLogs retrieves the logs of a container by name from the provider.
func (p *LocalDockerProvider) GetContainerLogs(namespace, podName, containerName string, tail int) (string, error) {
	return "", fmt.Errorf("not implemented: GetContainerLogs")
}

// NodeAddresses returns a list of addresses for the node status
// within Kubernetes.
func (p *LocalDockerProvider) NodeAddresses() []v1.NodeAddress {
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
func (p *LocalDockerProvider) NodeDaemonEndpoints() *v1.NodeDaemonEndpoints {
	return &v1.NodeDaemonEndpoints{
		KubeletEndpoint: v1.DaemonEndpoint{
			Port: p.daemonEndpointPort,
		},
	}
}