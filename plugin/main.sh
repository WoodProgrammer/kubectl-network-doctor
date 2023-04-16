#!/bin/bash

[[ -n $DEBUG ]] && set -x
set -eou pipefail
export CORE_DNS_NAMESPACE="${1:-kube-system}"

cecho(){
    BLACK='\033[1;30m'
	RED='\033[1;31m'
	GREEN='\033[1;32m'
	YELLOW='\033[1;33m'
	BLUE='\033[1;34m'
	PURPLE='\033[1;35m'
	CYAN='\033[1;36m'
	WHITE='\033[1;37m'
	NC="\033[0m" # No Color

    printf "${!1}${2} ${NC}\n"
}


cecho "RED" "Kubectl Network Doctor 0.0.1 is a triage plugin that we are using to get system dump of the essential network components"
cecho "RED" "This is demo version"

checkDNSResolution(){
	cecho "BLUE" "####### Checking DNS Resolution #########"

	cecho "YELLOW" "This manifests use default values in hosts.txt.You can speficy extra hosts."
	
	kubectl create configmap dns-func-test-cm	\
	 	--from-file=hosts.txt
	
	kubectl -n ${CORE_DNS_NAMESPACE} apply -f src/dns/k8s/manifests  2>/dev/null 
	
	kubectl wait --for=condition=Ready pod -l app=dns-func-test --timeout=60s
	kubectl logs -f -l app=dns-func-test
	cecho "BLUE" "###########################"

}

checkCoreDnsLogs(){
	cecho "YELLOW" "###########################"
	cecho "YELLOW" "Gathering CoreDNS Logs......"

	kubectl  -n ${CORE_DNS_NAMESPACE} logs -l k8s-app=kube-dns 2>/dev/null 
	cecho "YELLOW" "###########################"
}

checkCoreDnsPods(){

	cecho "RED" "###########################"
	cecho "RED" "Check coredns pod status"

	kubectl -n ${CORE_DNS_NAMESPACE} get pods -l k8s-app=kube-dns \
	 	-o jsonpath='{range .items[*]}{.status.containerStatuses[*].ready.true}{.metadata.name}{ "\n"}{end}'

	kubectl -n ${CORE_DNS_NAMESPACE} get deployment  |grep coredns
	cecho "RED" "###########################"

}

getTcpDump(){
	cecho "RED" "####### Gathering TCPDump from CoreDNS #########"

	CORE_DNS_PODS=$(kubectl -n ${CORE_DNS_NAMESPACE} get po  -l k8s-app=kube-dns  --no-headers -o custom-columns=":metadata.name" )

	mkdir -p tcpdump-coredns
	cecho "RED" "$CORE_DNS_PODS" 

	for pod in "${CORE_DNS_PODS[@]}"
	do
		echo $pod
		## add random timestamp to the nd-core-dns debug pods

		kubectl -n ${CORE_DNS_NAMESPACE} debug -q -i ${pod}  \
			 -c nd-core-dns \
			--image emirozbir/tcpdumper:latest -- tcpdump -U -i eth0 -w - > "tcpdump-coredns/${pod}-dump.pcap"
	done
	cecho "RED" ""###########################""

}

checkTraceRoute(){
	cecho "PURPLE" "####### TraceRoute Outputs #########"

	kubectl create configmap traceroute-test-cm	\
	 	--from-file=hosts.txt

	kubectl -n ${CORE_DNS_NAMESPACE} apply -f src/traceroute/k8s/manifests  2>/dev/null
	kubectl wait --for=condition=Ready pod -l app=traceroute-test --timeout=60s
	kubectl logs -f -l app=traceroute-test 
	
	cecho "PURPLE" "###########################"

}

tearDownDebugStack(){
	
	cecho "RED" "It is time to teardown entire stack"
	sleep 5
	kubectl delete cm dns-func-test-cm traceroute-test-cm
	kubectl delete -f src/traceroute/k8s/manifests
	kubectl delete -f src/dns/k8s/manifests
	cecho "RED" ""###########################""

}

checkCoreDnsLogs
checkCoreDnsPods
checkDNSResolution
checkTraceRoute
getTcpDump
tearDownDebugStack