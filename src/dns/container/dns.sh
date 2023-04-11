#!/bin/bash

export COREDNS_NAMESPACE="${1:-kube-system}"

echo "This script runs tests on coredns"

echo "1-) CoreDNS Logs"
echo "2-) CoreDNS Pods Health Status"
echo "3-) CoreDNS Pods TCPDump"

checkCoreDnsLogs(){
	kubectl logs -l k8s-app=kube-dns -n ${COREDNS_NAMESPACE} 2>/dev/null 
}

checkCoreDnsPods(){

	kubectl -n ${COREDNS_NAMESPACE} get pods -l k8s-app=kube-dns \
	 	-o jsonpath='{range .items[*]}{.status.containerStatuses[*].ready.true}{.metadata.name}{ "\n"}{end}'

	kubectl get deployment -n ${COREDNS_NAMESPACE} |grep coredns
}

getTcpDump(){
	CORE_DNS_PODS=$(kubectl get po -n ${COREDNS_NAMESPACE} -l k8s-app=kube-dns  --no-headers -o custom-columns=":metadata.name" )

	mkdir -p tcpdump-coredns
	echo $CORE_DNS_PODS
	for pod in "${CORE_DNS_PODS[@]}"
	do
		echo $pod

		kubectl debug -q -i ${pod}  \
			-n ${COREDNS_NAMESPACE} -c nd-core-dns \
			--image emirozbir/tcpdumper:latest -- tcpdump -U -i eth0 -w - > "tcpdump-coredns/${pod}-dump.pcap"
	done
}

getTcpDump