# kubectl network doctor

kubectl network doctor helps you to identify reason of the Kubernetes networking problems like DNS, east-west and south-north traffic communication.

This creates a dump of the network components in Kubernetes cluster.

# INSTALLATION


# USAGE


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
