#!/bin/bash
set -eou pipefail
[[ -n $DEBUG ]] && set -x

export CORE_DNS_NAMESPACE="${1:-kube-system}"
yellow=$(tput setaf 3 || true)

echo "This script runs tests on coredns"

echo "$yellow 1-) CoreDNS Logs"
echo "$yellow 2-) CoreDNS Pods Health Status"
echo "$yellow 3-) CoreDNS Pods TCPDump"

checkDNSResolution(){
	echo "$yellow This manifests use default values in hosts.txt.You can speficy extra hosts."
	kubectl -n ${CORE_DNS_NAMESPACE} apply -f src/dns/k8s/manifests  2>/dev/null 
}

checkCoreDnsLogs(){
	kubectl -n ${CORE_DNS_NAMESPACE} logs -l k8s-app=kube-dns  2>/dev/null 
}

checkCoreDnsPods(){

	kubectl -n ${CORE_DNS_NAMESPACE} get pods -l k8s-app=kube-dns \
	 	-o jsonpath='{range .items[*]}{.status.containerStatuses[*].ready.true}{.metadata.name}{ "\n"}{end}'

	kubectl -n ${CORE_DNS_NAMESPACE} get deployment  |grep coredns
}

getTcpDump(){
	CORE_DNS_PODS=$(kubectl -n ${CORE_DNS_NAMESPACE} get po  -l k8s-app=kube-dns  --no-headers -o custom-columns=":metadata.name" )

	mkdir -p tcpdump-coredns
	echo $CORE_DNS_PODS
	for pod in "${CORE_DNS_PODS[@]}"
	do
		echo $pod

		kubectl -n ${CORE_DNS_NAMESPACE} debug -q -i ${pod}  \
			 -c nd-core-dns \
			--image emirozbir/tcpdumper:latest -- tcpdump -U -i eth0 -w - > "tcpdump-coredns/${pod}-dump.pcap"
	done
}

checkTraceRoute(){
	kubectl -n ${CORE_DNS_NAMESPACE} apply -f src/traceroute/k8s/manifests  2>/dev/null
	kubectl wait --for=condition=Ready -l app=traceroute-test
	kubectl logs -f -l app=traceroute-test
}

checkCoreDnsLogs
checkCoreDnsPods
checkTraceRoute
checkDNSResolution
getTcpDump