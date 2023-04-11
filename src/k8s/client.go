package main

import (
	"context"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	cfg, err := clientcmd.BuildConfigFromFlags(
		"",
		filepath.Join(homedir.HomeDir(), ".kube", "config"),
	)
	handleError(err)

	k8s, err := kubernetes.NewForConfig(cfg)
	handleError(err)

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "nd-dns-test-pod-one"},
		Spec: v1.PodSpec{
			RestartPolicy: v1.RestartPolicyNever,
			Containers: []v1.Container{
				v1.Container{
					Name:    "dns-checker",
					Image:   "python:3.8",
					Command: []string{"python"},
					Args:    []string{"-c", "print('hello world')"},
				},
			},
		},
	}

	_, err = k8s.CoreV1().Pods("network-doctor").Create(
		context.Background(),
		pod,
		metav1.CreateOptions{},
	)

	handleError(err)
}
