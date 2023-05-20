package main

import (
	"bytes"
	"context"
	"io"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func getPodLogs(clientset *kubernetes.Clientset) {

	podLogOpts := corev1.PodLogOptions{}

	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"app": "dns-test"}}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}

	pods, _ := clientset.CoreV1().Pods("kube-system").List(context.TODO(), listOptions)

	for _, pod := range pods.Items {

		req := clientset.CoreV1().Pods("kube-system").GetLogs(pod.Name, &podLogOpts)
		podLogs, err := req.Stream(context.TODO())
		if err != nil {
			ErrorLogger.Println(err)
			log.Fatal("error in opening stream")
		}

		defer podLogs.Close()

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, podLogs)
		if err != nil {
			log.Fatal("")
		}
		str := buf.String()
		InfoLogger.Println(str)
	}

}
