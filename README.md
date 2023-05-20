# kubectl network doctor

kubectl network doctor helps you to identify reason of the Kubernetes networking problems like DNS, east-west and south-north traffic communication.

# Build from Scratch

This plugin has very simple installation method basically you just need to move this script somewhere in the output of the `$PATH` output on your local.

```sh
git clone git@github.com:WoodProgrammer/kubectl-network-doctor.git

pushd kubectl-network-doctor/plugin
    go build -o kubectl-nd .
    mv kubectl-nd /usr/local/bin 
popd

```

# USAGE

kubectl-network-doctor basically provides three module just for now;
* coredns (status/logs)
* traceroute & dns resolution tests
* tcpdump with debug containers

## coredns mode 
It basically check status of coredns replicas and logs.After this operation done, kubectl-nd creates one deployment with the name `dns-test`.

This deployment object run a simple pod which runs and measure the dns records which are already initalizated on `hosts.txt` file.

You can specficy specific host addresses to measure the dns resolution time.The source code of dns-test is already located under `utils/src/dns` directory.


```sh
 $ kubectl nd mode dns
 ... 
 ....
 .....

 100% |██████████████████████████████████████████████████████| (50/50, 1 it/s)
{[{www.youtube.com [] 75 false} {www.google.com [] 48 false} {ifconfig.co [] 38 false}]}
```

The output provides three important values;

`{HOST_ADDR ResolutionTime ErrStatus }`

You can populate multiple values on hosts.txt shown at below.

```txt
db.internal.co
ssm.eu-central-1.amazonaws.com
svc.default.cluster.local
svc.ns-test.cluster.local
www.google.com
```

## traceroute mode 
The traceroute mode run traceroute command and it shows the trace outputs in hosts.txt file.

```sh
    $ kubectl nd mode traceroute
    ....
    ....
    ...

    Creating deployment...
    Created deployment "traceroute-deployment".
    INFO:::Waiting for the results of the logs until the TraceRoute deployment get ready approx:: 50 sec
    100% |████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████| (50/50, 1 it/s)
    tput: No value for $TERM and no -T specified
    TraceRoute debug on this address:: www.youtube.com
    traceroute to www.youtube.com (142.250.187.110), 30 hops max, 60 byte packets
    1  10.244.0.1 (10.244.0.1)  0.065 ms  0.053 ms  0.035 ms
    2  172.18.0.1 (172.18.0.1)  0.059 ms  0.021 ms  0.047 ms
    3  192.168.0.1 (192.168.0.1)  11.294 ms  15.891 ms  17.705 ms
    4  * * *
    .....
    .......

```
## tcpdump mode 
Let's check the tcpdump usage.In this mode it basically creates debug container in specified pod.

<b>Caveats</b>

If you would like to run tcpdump mode please be ensure `EphemeralContainers` flags enabled on your cluster.

```sh
    kubectl nd mode tcpdump --pod target-pod-name --namespace kube-system --file test.pcap
```

If you do not specify name of the file, kubectl-nd creates default file by the date prefix.

So you are able to produce pcap outputs from you code block.

## The Source of the Test Components

<a href="plugin/src/dns/">Dns Tester</a>
<br>
<a href="plugin/src/traceroute/">TraceRoute Tester</a>



# Concepts and Goals
The problems that covered by this plugin;

* DNS;
    - [x] dns query test optional internal external
    - [x] core-dns logs
    - [x] core-dns healthcheck
    - memory and cpu usage
    - host configuration under these files;
        * /etc/hosts
        * /etc/resolv.conf
        * nsswitch.conf

* TCP/IP
    - [x] Pcap dump 
    - [x] traceroute outputs
    - Direct analysis based on bpf rules.

* TODO
    * pod-to-pod communication (based on topology keys or labels)
    * node-to-node communication
    * node network configuration
    * TearDown entire stack.
    * Custom TCPDump filters and command options

# RoadMap

We will provide command line flags and other user friendly options to this plugin.

# Related Links;

<a href="https://www.wireshark.org/">WireShark Installation </a>
