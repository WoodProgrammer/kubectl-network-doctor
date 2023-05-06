package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func readHostList(file string) string {

	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	text := string(fileContent)

	return text

}

func generateHostsFile(fileName string) map[string]string {

	configMap := make(map[string]string)

	if _, err := os.Stat(fileName); err == nil {
		fmt.Println("The file is exist:: %s", fileName)
		configMap["data"] = readHostList(fileName)

	} else {
		fmt.Println("The file %s is not exist in generating with default addresses", fileName)
		configMap["data"] = "www.youtube.com\nwww.google.com\nifconfig.co"

	}

	return configMap

}

func createConfigMap(configMapName string, namespaceName string, data map[string]string, clientset *kubernetes.Clientset) {

	cm := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespaceName,
		},
		Data: data,
	}

	err, result := clientset.CoreV1().ConfigMaps("kube-system").Create(context.TODO(), &cm, metav1.CreateOptions{})

	if err != nil {
		fmt.Println("ERROR", err)
	} else {
		fmt.Println(result)
	}
}
