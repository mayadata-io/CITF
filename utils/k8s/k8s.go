/*
Copyright 2018 The OpenEBS Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package k8s

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"errors"

	"github.com/golang/glog"
	"github.com/openebs/CITF/common"
	strutil "github.com/openebs/CITF/utils/string"
	sysutil "github.com/openebs/CITF/utils/system"
	core_v1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// GetAllNamespacesCoreV1NamespaceArray returns V1NamespaceList of all the namespaces.
// :return kubernetes.client.models.v1_namespace_list.V1NamespaceList: list of namespaces.
func (k8s K8S) GetAllNamespacesCoreV1NamespaceArray() ([]core_v1.Namespace, error) {
	nsList, err := k8s.Clientset.CoreV1().Namespaces().List(meta_v1.ListOptions{})
	return nsList.Items, err
}

// GetAllNamespacesMap returns list of the names of all the namespaces.
// :return: map[string]core_v1.Namespace: map of namespaces where key is namespace name (str)
// and value is corresponding k8s.io/api/core/v1.Namespace object.
func (k8s K8S) GetAllNamespacesMap() (map[string]core_v1.Namespace, error) {
	namespacesList, err := k8s.GetAllNamespacesCoreV1NamespaceArray()
	if err != nil {
		return nil, err
	}

	namespaces := map[string]core_v1.Namespace{}
	for _, ns := range namespacesList {
		namespaces[ns.Name] = ns
	}
	return namespaces, nil
}

// GetPod returns the Pod object for given podName in the given namespace.
// :return: *kubernetes.client.models.v1_pod.V1Pod: Pointer to Pod objects.
func (k8s K8S) GetPod(namespace, podName string) (*core_v1.Pod, error) {
	podsClient := k8s.Clientset.CoreV1().Pods(namespace)
	return podsClient.Get(podName, meta_v1.GetOptions{})
}

// GetPods returns all the Pods object which has a prefix specified in its name in the given namespace.
// :return: []kubernetes.client.models.v1_pod.V1Pod: Slice of Pod objects.
func (k8s K8S) GetPods(namespace, podNamePrefix string) ([]core_v1.Pod, error) {
	// Try to get the pod for 10 times as sometime code reaches
	// when pod is not even in ContainerCreating state
	i := 0
	thePods := []core_v1.Pod{}
	for len(thePods) == 0 && i < 10 {
		time.Sleep(2 * time.Second)

		// List pods
		pods, err := k8s.Clientset.CoreV1().Pods(namespace).List(meta_v1.ListOptions{})
		if err != nil {
			fmt.Printf("Error occurred: %+v\n", err)
		}

		// Find the Pod
		if common.DebugEnabled {
			fmt.Println(strings.Repeat("*", 80))
			fmt.Printf("Current pods in %q namespace are:\n", namespace)
		}
		for _, pod := range pods.Items {
			if common.DebugEnabled {
				fmt.Println("Complete Pod name is:", pod.Name)
			}
			if strings.HasPrefix(pod.Name, podNamePrefix) {
				thePods = append(thePods, pod)
			}
		}
		if common.DebugEnabled {
			fmt.Println(strings.Repeat("*", 80))
		}
		i++
	}

	if len(thePods) == 0 {
		return thePods, errors.New("failed getting NDM-Pod in given time")
	}

	return thePods, nil
}

// ReloadPod reloads the state of the pod supplied and return the recent one
func (k8s K8S) ReloadPod(pod core_v1.Pod) (*core_v1.Pod, error) {
	return k8s.GetPod(pod.Namespace, pod.Name)
}

// GetPodPhase returns phase of the pod passed as an k8s.io/api/core/v1.PodPhase object.
//		:param k8s.io/api/core/v1.Pod pod: pod object for which you want to get phase.
//		:return: k8s.io/api/core/v1.PodPhase: phase of the pod.
func (k8s K8S) GetPodPhase(pod core_v1.Pod) core_v1.PodPhase {
	return pod.Status.Phase
}

// GetPodPhaseStr returns phase of the pod passed in string format.
//		:param k8s.io/api/core/v1.Pod pod: pod object for which you want to get phase.
//		:return: str: phase of the pod.
func (k8s K8S) GetPodPhaseStr(pod core_v1.Pod) string {
	return string(k8s.GetPodPhase(pod))
}

// GetContainerStateInPod returns the state of the container of supplied index of the supplied Pod.
//    :param containerIndex: index of the container for which you want state.
//    :param timeout: maximum time duration to get the container's state.
//                       This method does not very strictly obey this param.
//    :return: k8s.io/api/core/v1.ContainerState: state of the container.
func (k8s K8S) GetContainerStateInPod(pod *core_v1.Pod, containerIndex int, timeout time.Duration) (core_v1.ContainerState, error) {
	if pod == nil {
		return core_v1.ContainerState{}, errors.New("nil argument supplied for pod")
	}

	var err error
	startTime := time.Now()
	for reflect.DeepEqual(pod.Status.ContainerStatuses, []core_v1.ContainerStatus(nil)) && time.Since(startTime) < timeout {
		time.Sleep(time.Second)
		pod, err = k8s.ReloadPod(*pod)
		if err != nil {
			return core_v1.ContainerState{}, err
		}
	}
	if time.Since(startTime) >= timeout {
		return core_v1.ContainerState{}, fmt.Errorf("pod %q of namespace %q had no container till %v", pod.Name, pod.Namespace, timeout)
	}

	for len(pod.Status.ContainerStatuses) <= containerIndex && time.Since(startTime) < timeout {
		time.Sleep(time.Second)
		pod, err = k8s.ReloadPod(*pod)
		if err != nil {
			return core_v1.ContainerState{}, err
		}
	}
	if time.Since(startTime) >= timeout {
		return core_v1.ContainerState{}, fmt.Errorf("pod did not had %d containers till %v", containerIndex+1, timeout)
	}

	return pod.Status.ContainerStatuses[containerIndex].State, nil
}

// GetNodes returns a list of all the nodes.
//    :return: slice: list of nodes (slice of k8s.io/api/core/v1.Node array).
func (k8s K8S) GetNodes() (nodeNames []core_v1.Node, err error) {
	nodeNames = []core_v1.Node{}

	// To handle latency it tries 10 times each after 1 second of wait
	waited := 0
	for waited < 10 {
		nodeList, err := k8s.Clientset.CoreV1().Nodes().List(meta_v1.ListOptions{})
		if err != nil {
			break
		} else if len(nodeList.Items) == 0 {
			time.Sleep(time.Second)
			waited++
			continue
		}
		nodeNames = nodeList.Items
		break
	}

	return
}

// GetNodeNames returns a list of the name of all the nodes.
//    :return: slice: list of node names (slice of string array).
func (k8s K8S) GetNodeNames() (nodeNames []string, err error) {
	nodeNames = []string{}

	nodes, err := k8s.GetNodes()
	if err != nil {
		return
	}
	for _, node := range nodes {
		nodeNames = append(nodeNames, node.Name)
	}

	return
}

// TODO: Write a function to label the node
// LabelNode label the node with the given key and value.
//    :param string node_name: Name of the node.
//    :param string key: Key of the label.
//    :param string value: Value of the label.
//    :return: error: if any error occurred or nil otherwise.
// func LabelNode(nodeName, key, value string) error { return fmt.Errorf("Not Implemented") }

// GetDaemonset returns the k8s.io/api/extensions/v1beta1.DaemonSet for the name supplied.
func (k8s K8S) GetDaemonset(daemonsetName, daemonsetNamespace string) (v1beta1.DaemonSet, error) {
	daemonsetClient := k8s.Clientset.ExtensionsV1beta1().DaemonSets(daemonsetNamespace)
	ds, err := daemonsetClient.Get(daemonsetName, meta_v1.GetOptions{})
	if err != nil {
		return v1beta1.DaemonSet{}, err
	}
	return *ds, nil
}

// ApplyDSFromManifestStruct Creates a Daemonset from the manifest supplied
func (k8s K8S) ApplyDSFromManifestStruct(manifest v1beta1.DaemonSet) (v1beta1.DaemonSet, error) {
	if manifest.Namespace == "" {
		manifest.Namespace = core_v1.NamespaceDefault
	}
	daemonsetClient := k8s.Clientset.ExtensionsV1beta1().DaemonSets(manifest.Namespace)
	ds, err := daemonsetClient.Create(&manifest)
	if err != nil {
		return v1beta1.DaemonSet{}, err
	}
	return *ds, nil
}

// GetDaemonsetStructFromYamlBytes returns k8s.io/api/extensions/v1beta1.DaemonSet
// for the yaml supplied
func (k8s K8S) GetDaemonsetStructFromYamlBytes(yamlBytes []byte) (v1beta1.DaemonSet, error) {
	ds := v1beta1.DaemonSet{}

	jsonBytes, err := strutil.ConvertYAMLtoJSON(yamlBytes)
	if err != nil {
		return ds, fmt.Errorf("error while Converting yaml string into Daemonset Structure. Error: %+v", err)
	}

	err = json.Unmarshal(jsonBytes, &ds)
	if err != nil {
		return ds, fmt.Errorf("error occurred while marshaling into Daemonset struct. Error: %+v", err)
	}

	return ds, nil
}

// TODO: Write a function to apply the YAML with the help of client-go
// YAMLApply apply the yaml specified by the argument.
//    :param str yamlPath: Path of the yaml file that is to be applied.
// func YAMLApplyAPI(yamlPath string) error { return fmt.Errorf("Not Implemented") }

// YAMLApply apply the yaml specified by the argument.
//    :param str yamlPath: Path of the yaml file that is to be applied.
func (k8s K8S) YAMLApply(yamlPath string) error {
	// TODO: Try using API call first. i.e. Using client-go

	err := sysutil.RunCommand("kubectl apply -f " + yamlPath)
	if err != nil {
		glog.Errorf("error occurred while applying the %s. Error: %+v", yamlPath, err)
		return fmt.Errorf("failed applying %s", yamlPath)
	}
	return nil
}

// ExecToPodThroughAPI performs non-interactive exec to the pod with the specified command using client-go.
// :param string command: list of the str which specify the command.
// :param string containerName: name of the container in the Pod. (If the Pod has only one container, then it can be Empty String)
// :param string pod_name: Pod name
// :param string namespace: namespace of the Pod. (If it is blank string then, namespace will be default i.e. k8s.io/api/core/v1.NamespaceDefault)
// :param io.Reader stdin: Standard Input if necessary, otherwise `nil`
// :return: string: Output of the command. (STDOUT)
//          string: Errors. (STDERR)
//           error: If any error has occurred otherwise `nil`
func (k8s K8S) ExecToPodThroughAPI(command, containerName, podName, namespace string, stdin io.Reader) (string, string, error) {
	if len(namespace) == 0 {
		namespace = core_v1.NamespaceDefault
	}

	req := k8s.Clientset.Core().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")
	scheme := runtime.NewScheme()
	if err := core_v1.AddToScheme(scheme); err != nil {
		return "", "", fmt.Errorf("error adding to scheme: %v", err)
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	podExecOptions := core_v1.PodExecOptions{
		Command: strings.Fields(command),
		Stdin:   stdin != nil,
		Stdout:  true,
		Stderr:  true,
		TTY:     false,
	}
	if len(containerName) != 0 {
		podExecOptions.Container = containerName
	}

	req.VersionedParams(&podExecOptions, parameterCodec)

	if common.DebugEnabled {
		fmt.Println("Request URL: ", req.URL().String())
	}

	exec, err := remotecommand.NewSPDYExecutor(k8s.Config, "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("error while creating Executor: %v", err)
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return "", "", fmt.Errorf("error in Stream: %v", err)
	}

	return stdout.String(), stderr.String(), nil
}

// ExecToPodThroughKubectl performs non-interactive exec to the pod with the specified command using `kubectl exec`
// :param string command: list of the str which specify the command.
// :param string containerName: name of the container in the Pod. (If the Pod has only one container, then it can be Empty String)
// :param string pod_name: Pod name
// :param string namespace: namespace of the Pod. (If it is blank string then, namespace will be default i.e. k8s.io/api/core/v1.NamespaceDefault)
// :return: string: Output of the command. (STDOUT)
//           error: If any error has occurred otherwise `nil`
func (k8s K8S) ExecToPodThroughKubectl(command, containerName, podName, namespace string) (string, error) {
	kubectlCommand := "kubectl"

	// adding namespace if namespace is not blank string
	if len(namespace) != 0 {
		kubectlCommand += " -n " + namespace
	}

	// adding podName
	kubectlCommand += " exec " + podName

	// adding container name if containerName is not a blank string
	if len(containerName) != 0 {
		kubectlCommand += " -c " + containerName
	}

	// finally adding command to execute
	kubectlCommand += " -- " + command

	return sysutil.ExecCommand(kubectlCommand)
}

// ExecToPod performs non-interactive exec to the pod with the specified command.
// first through API with `stdin` param as `nil`, if it fails then it uses `kubectl exec`
// :param string command: list of the str which specify the command.
// :param string containerName: name of the container in the Pod. (If the Pod has only one container, then it can be Empty String)
// :param string pod_name: Pod name
// :param string namespace: namespace of the Pod. (If it is blank string then, namespace will be default i.e. k8s.io/api/core/v1.NamespaceDefault)
// :return: string: Output of the command. (STDOUT)
//           error: If any error has occurred otherwise `nil`
func (k8s K8S) ExecToPod(command, containerName, podName, namespace string) (string, error) {
	stdout, stderr, err := k8s.ExecToPodThroughAPI(command, containerName, podName, namespace, nil)
	if err == nil {
		return stdout, nil
	}

	// When Exec through API fails
	glog.Errorf("error while exec into Pod through API. Stderr: %q. Error: %+v", stderr, err)
	return k8s.ExecToPodThroughKubectl(command, containerName, podName, namespace)
}

// GetLog returns the log of the pod.
// :param string pod_name: Name of the pod. (required)
// :param string namespace: Namespace of the pod. (required)
// :return: string: Log of the pod specified.
//           error: If an error has occurred, otherwise `nil`
// TODO: Fix in API call (Error: GroupVersion is required when initializing a RESTClient)
func (k8s K8S) GetLog(podName, namespace string) (string, error) {
	// We can't declare a variable somewhere which can be skipped by goto
	var req *rest.Request
	var readCloser io.ReadCloser
	var err error

	buf := new(bytes.Buffer)
	req = k8s.Clientset.CoreV1().Pods(namespace).GetLogs(
		podName,
		&core_v1.PodLogOptions{},
	)

	readCloser, err = req.Stream()
	defer readCloser.Close()
	if err != nil {
		goto use_kubectl
	}

	buf.ReadFrom(readCloser)
	if common.DebugEnabled {
		fmt.Println("Log of Pod", podName, "in Namespace", namespace, "through API:")
		fmt.Println(buf.String())
	}
	return buf.String(), nil

use_kubectl:
	glog.Errorf("Error while getting log with API call. Error: %+v", err)

	return sysutil.ExecCommand("kubectl -n " + namespace + " logs " + podName)
}
