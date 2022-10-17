
spectacles: a Kubernetes event exporter

While attempting to write [myasnik][myasnik], I ran into problems making a
functional controller. Spectacles is a simpler app with some of the same
components (reflected, informer, lister) and setup (go modules, client-go).

[myasnik]: https://github.com/ripta/myasnik

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
for it. Its service account needs to be able to get events in the cluster, e.g.:

```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: spectacles
rules:
  - apiGroups:
      - "" # core
      - events.k8s.io
    resources:
      - events
    verbs:
      - get
      - list
      - watch
```

There is no leader election (yet?), so you'll need to run only one copy.
