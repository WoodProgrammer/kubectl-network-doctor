package main

import (
	"context"
	"fmt"

	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func generateRandomContainerName() string {

	suffix := RandStringRunes(4)
	containerName := fmt.Sprintf("debugger-%s", suffix)
	return containerName
}

func createDebugContainer(namespaceName string, clientset *kubernetes.Clientset) string {
	debugContainerName := generateRandomContainerName()

	command := []string{"sleep", "10000000000"}
	pod, _ := clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), "coredns-64897985d-r99fw", metav1.GetOptions{})

	ec := coreV1.EphemeralContainer{
		EphemeralContainerCommon: coreV1.EphemeralContainerCommon{
			Name:    debugContainerName,
			Image:   fmt.Sprintf("%s:%s", "emirozbir/tcpdumper", "latest"),
			Command: command,
		},
	}

	pod.Spec.EphemeralContainers = append(pod.Spec.EphemeralContainers, ec)

	_, err := clientset.CoreV1().Pods(pod.Namespace).UpdateEphemeralContainers(context.TODO(), pod.Name, pod, metav1.UpdateOptions{})

	if err != nil {
		fmt.Println(err)
	}

	return debugContainerName

}
