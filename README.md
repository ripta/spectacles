
spectacles: a Kubernetes event exporter

While attempting to write [myasnik][myasnik], I ran into problems making a
functional controller. Spectacles is a simpler app with some of the same
components (reflected, informer, lister) and setup (go modules, client-go).

# Compiling

```
git clone https://github.com/ripta/spectacles
cd spectacles
make build
```

# Running

```
export KUBECONFIG=~/.kube/config # or other path
make run
```

# Deploying

spectacles can run in-cluster, but you're on your own in build a docker image
for it. Its service account needs to be able to get events in the cluster.

There is no leader election, so you'll need to run only one copy.
