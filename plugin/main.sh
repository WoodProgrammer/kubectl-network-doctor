#!/bin/bash
export BASE_PATH=$(pwd)
[[ -n $DEBUG ]] && set -x

set -eou pipefail
set -o nounset

IFS=$'\n\t'

trap ctrl_c INT

ctrl_c(){
	tearDownDebugStack
}

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

cecho "BLUE" "KUBECTL Network Doctor Version number::0.0.1"


cat << "EOF"
        .----. 
       ===(_)==   Kubectl Network Doctor 0.0.1
      // 6  6 \\  /
      (    7   )
       \ '--' /
        \_ ._/
       __)  (__
    /"`/`\`V/`\`\
   /   \  `Y _/_ \
  / [DR]\_ |/ / /\
  |     ( \/ / / /
   \  \  \      /
    \  `-/`  _.`
     `=._`=./
EOF

cecho "RED" "This script basically creates number of pods \n and run some queries please check WoodProgrammer/kubectl-network-doctor"

cecho "BLUE" "Kubectl Network Doctor 0.0.1 is a triage plugin that we are using to get system dump of the essential network components"
cecho "YELLOW" "This is demo version"
sleep 4

checkDNSResolution(){
	cecho "BLUE" "####### Checking DNS Resolution #########"
	cecho "YELLOW" "This manifests use default values in hosts.txt.You can speficy extra hosts."
	
	kubectl create configmap dns-func-test-cm --from-file=hosts.txt -o yaml --dry-run=client |kubectl apply -n ${CORE_DNS_NAMESPACE} -f -
	
	kubectl -n ${CORE_DNS_NAMESPACE} apply -f ${BASE_PATH}/plugin/src/dns/k8s/manifests  2>/dev/null 
	
	#kubectl wait --for=condition=Ready pod -l app=dns-func-test --timeout=60s easy way :=) 
	sleep 90
	kubectl logs -f -l app=dns-func-test


	cecho "BLUE" "--------------------------"

}

checkCoreDnsLogs(){
	cecho "YELLOW" "--------------------------"
	cecho "YELLOW" "Gathering CoreDNS Logs......"

	kubectl  -n ${CORE_DNS_NAMESPACE} logs -l k8s-app=kube-dns 2>/dev/null 
	cecho "YELLOW" "--------------------------"
}

checkCoreDnsPods(){

	cecho "RED" "--------------------------"
	cecho "RED" "Check coredns pod status"

	kubectl -n ${CORE_DNS_NAMESPACE} get pods -l k8s-app=kube-dns \
	 	-o jsonpath='{range .items[*]}{.status.containerStatuses[*].ready.true}{.metadata.name}{ "\n"}{end}'

	kubectl -n ${CORE_DNS_NAMESPACE} get deployment  |grep coredns

	cecho "RED" "--------------------------"

}

callDebugTcpDump(){
	export pod=$1

	if [ -z $pod ];
	then
		cecho "RED" "ERROR please identify name of coredns pods"
	else
		cecho "RED" "WARNING:: Gathering TCPDump of the pod :: ${pod}"


		if [ "${TCP_DUMP_MODE}" = "wireshark" ];
		then
			cecho "GREEN" "WARNING:: To continue the stream analyse you can can close Wireshark"

			kubectl -n kube-system debug -q -i ${pod} --image nicolaka/netshoot \
				-- timeout ${TCP_DUMP_TIMEOUT} tcpdump -i eth0 -w - | wireshark -k -i -
		else
			kubectl -n kube-system debug -q -i ${pod} --image nicolaka/netshoot \
				-- timeout ${TCP_DUMP_TIMEOUT} tcpdump -i eth0 -w - > tcpdump-coredns/${pod}.pcap
		fi

		if [ -s "tcpdump-coredns/${pod}.pcap" ];
		then
			cecho "BLUE" "INFO:: Tcpdump file has already written properly"
		else
			cecho "RED" "ERROR:: Tcpdump file is empty"
		fi
	fi
}	

getTcpDump(){
	cecho "RED" "####### Gathering TCPDump from CoreDNS #########"

	CORE_DNS_PODS=($(kubectl -n ${CORE_DNS_NAMESPACE} get po  -l k8s-app=kube-dns  --no-headers -o custom-columns=":metadata.name" ))
	mkdir -p tcpdump-coredns
	cecho "RED" "$CORE_DNS_PODS" 

	for pod in "${CORE_DNS_PODS[@]}"
	do
		echo $pod
	
		callDebugTcpDump $pod

		echo "Result is $?"
	done
	cecho "RED" ""--------------------------""

}

checkTraceRoute(){
	cecho "PURPLE" "####### TraceRoute Outputs #########"

	kubectl create configmap traceroute-test-cm --from-file=hosts.txt -o yaml --dry-run=client | kubectl apply -n ${CORE_DNS_NAMESPACE} -f  -
	kubectl -n ${CORE_DNS_NAMESPACE} apply \
		-f ${BASE_PATH}/plugin/src/traceroute/k8s/manifests  2>/dev/null

	#kubectl wait --for=condition=Ready pod \
	#	-l app=traceroute-test --timeout=60s
	sleep 90

	kubectl logs -f -l app=traceroute-test 
	
	cecho "PURPLE" "--------------------------"

}

tearDownDebugStack(){
	
	cecho "RED" "It is time to teardown entire stack"
	sleep 5
	kubectl delete cm dns-func-test-cm traceroute-test-cm
	
	kubectl delete deployment traceroute-test-deployment
	kubectl delete deployment dns-func-test-deployment
	
	cecho "RED" ""--------------------------""

}

createHostFile(){
	FILE_STAT=$(file -s hosts.txt)
	if [ -e hosts.txt ];
	then
		cecho "BLUE" "Hosts file exists"
	else
		cecho "BLUE" "Creating hosts file "
cat > hosts.txt <<EOF
www.youtube.com
www.google.com
google.com
EOF
	fi
}

main(){

  export CORE_DNS_NAMESPACE="${1:-kube-system}"
  export TCP_DUMP_TIMEOUT="${2:-10}"
  export TCP_DUMP_MODE="${3:-wireshark}"

  cecho "RED" "CoreDNS namespace:: ${CORE_DNS_NAMESPACE}"
  cecho "RED" "TcpDump Timeout is:: ${TCP_DUMP_TIMEOUT}"
  cecho "RED" "TcpDump Mode:: ${TCP_DUMP_MODE}"

  createHostFile
  checkCoreDnsLogs
  checkCoreDnsPods
  checkDNSResolution
  checkTraceRoute
  getTcpDump
  tearDownDebugStack

}

main
