package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// emirozbir/traceroute-test:0.0.1
// []string{"./main.sh"}

func gatherLogs(deploymentLabel string, namespaceName string, clientset *kubernetes.Clientset) {

	podLogOpts := corev1.PodLogOptions{}

	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"app": deploymentLabel}}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}

	pods, _ := clientset.CoreV1().Pods(namespaceName).List(context.TODO(), listOptions)
	for _, pod := range pods.Items {

		req := clientset.CoreV1().Pods(namespaceName).GetLogs(pod.Name, &podLogOpts)
		podLogs, err := req.Stream(context.TODO())
		if err != nil {
			fmt.Println(err)
			log.Fatal("error in opening stream")
		}

		defer podLogs.Close()

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, podLogs)
		if err != nil {
			log.Fatal("")
		}
		str := buf.String()
		fmt.Println(str)
	}
}

func createDeployment(deploymentName string, imageName string, command []string, namespaceName string, clientset *kubernetes.Clientset) {

	deploymentsClient := clientset.AppsV1().Deployments("kube-system")

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName + "-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploymentName + "-test",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deploymentName + "-test",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:    deploymentName + "-test",
							Image:   imageName,
							Command: command,
						},
					},
				},
			},
		},
	}

	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

func int32Ptr(i int32) *int32 { return &i }
