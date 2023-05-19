package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

func ExecuteRemoteCommand(command string, tcpDumpFileName string, namespaceName string, targetPodName string, containerName string) (string, string, error) {
	kubeCfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	restCfg, err := kubeCfg.ClientConfig()
	if err != nil {
		return "", "", err
	}
	coreClient, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return "", "", err
	}

	var fileWriter io.Writer

	if tcpDumpFileName != "" {
		fileWriter, err = os.Create("tcpdump-file.pcap")

	} else {
		fileWriter, err = os.Create(tcpDumpFileName)
	}

	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}

	request := coreClient.CoreV1().RESTClient().
		Post().
		Namespace(namespaceName).
		Resource("pods").
		Name(targetPodName).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command:   []string{"/bin/bash", "-c", command},
			Container: containerName,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(restCfg, "POST", request.URL())
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: fileWriter,
		Stderr: errBuf,
		Tty:    false,
	})

	if err != nil {
		return "", "", fmt.Errorf("%w Failed executing command %s on %v/%v", err, command, "kube-system", containerName)
	}

	return buf.String(), errBuf.String(), nil
}
