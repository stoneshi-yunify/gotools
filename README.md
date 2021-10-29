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