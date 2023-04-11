# kubectl network doctor

This is a kubectl plugin it basically shows the network problems like;


# Concepts and Goals
The problems that covered by this plugin;

* DNS;
    - core-dns logs
    - core-dns healthcheck
    - memory and cpu usage
    - host configuration under these files;
        * /etc/hosts
        * /etc/resolv.conf
        * nsswitch.conf

* TCP/IP
    - Pcap dump 
    - Direct analysis based on bpf rules.

* pod-to-pod communication (based on topology keys or labels)
* node-to-node communication
* node network configuration
* traceroute outputs