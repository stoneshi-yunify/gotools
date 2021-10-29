## kube-apiserver SANs Adder
Add IPs to extra Subject Alternative Names (SANs) of apiserver certificate, so they can be used to serve kube-apiserver.

### Build
```
$ make build-apiserver-sans-adder
```

### Run
1. Copy `./tmp/apiserver-sans-adder` to `/etc/kubernetes/pki` on your masters
2. Add IPs to extra SANs of kube-apiserver, do this on all masters
```
$ chmod +x ./apiserver-sans-adder
$ ./apiserver-sans-adder -extra-san-ips "1.2.3.4,5.6.7.8"
```
3. Restart all kube-apiserver pods
```
$ kubectl -n kube-system get pod | awk '{print $1}' | grep "^kube-apiserver" | xargs -n 1 kubectl -n kube-system delete pod
```
4. Replace the `server` of your kubeconfig with above IP, and try it.

### Example
```
root@node1:/etc/kubernetes/pki# ./apiserver-sans-adder -h
Usage of ./apiserver-sans-adder:
  -apiserver-cert string
    	apiserver certificate path (default "apiserver.crt")
  -ca-cert string
    	CA certificate path (default "ca.crt")
  -ca-key string
    	CA key path (default "ca.key")
  -extra-san-ips string
    	extra Subject Alternative Names (SANs) to use for the API Server serving certificate. Can only be IP addresses. Separated by comma.
  -new-apiserver-cert string
    	new apiserver certificate path (default "apiserver.crt")
  -new-apiserver-key string
    	new apiserver key path (default "apiserver.key")
root@node1:/etc/kubernetes/pki# ./apiserver-sans-adder -extra-san-ips "9.9.9.9"
Done
root@node1:/etc/kubernetes/pki# openssl x509 -in apiserver.crt -noout -ext "subjectAltName"
X509v3 Subject Alternative Name:
    DNS:kubernetes, DNS:kubernetes.default, DNS:kubernetes.default.svc, DNS:kubernetes.default.svc.cluster.local, DNS:lb.kubesphere.local, DNS:node1, IP Address:10.233.0.1, IP Address:192.168.0.12, IP Address:6.7.8.9, IP Address:5.6.7.8, IP Address:7.7.7.7, IP Address:139.198.15.193, IP Address:3.3.3.3, IP Address:3.3.3.3, IP Address:1.2.3.4, IP Address:9.9.9.9
```
