# pod-spec-mutator

Can adjust the pod spec before it is created.

## Dev

```bash
k delete -f deploy/mutatingWebhookConfiguration.yaml
telepresence helm install
kubens injector
telepresence intercept pod-spec-mutator --port 8443
```

```bash
make telepresence-get-certs
make deploy-to-telepresence
make exec-telepresence
telepresence intercept pod-spec-mutator-postgresoperator --port 8443
export PATCH_JSON='{"hostAliases":[{"ip":"192.168.1.100","hostnames":["foo.local"]}]}'
dlv --listen :2345 --headless --api-version=2 exec /tmp/pod-spec-mutator
```

```bash
telepresence leave pod-spec-mutator
telepresence uninstall --all-agents
telepresence quit
telepresence helm uninstall
```

### Test / Use

```bash
k create namespace injector
k create -f deploy/*.yaml
go run main.go &
```

# TODO:

- [x] excempt some pods from injection (especially the injector pod itself)
- [x] limit the injection to specific namespaces
- [ ] add a configmap to configure the host alias
