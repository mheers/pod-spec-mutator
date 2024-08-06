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
export NAMESPACE_REGEX='postgresoperator'
export POD_NAME_REGEX='^(logical-backup|default-).*$'
export PATCH_SPEC_JSON='{"ActiveDeadlineSeconds":null,"Affinity":null,"AutomountServiceAccountToken":null,"Containers":[{"Args":null,"Command":null,"Env":null,"EnvFrom":null,"Image":null,"ImagePullPolicy":"Never","Lifecycle":null,"LivenessProbe":null,"Name":"postgres","Ports":null,"ReadinessProbe":null,"ResizePolicy":null,"Resources":null,"RestartPolicy":null,"SecurityContext":null,"StartupProbe":null,"Stdin":null,"StdinOnce":null,"TerminationMessagePath":null,"TerminationMessagePolicy":null,"Tty":null,"VolumeDevices":null,"VolumeMounts":null,"WorkingDir":null}],"DnsConfig":null,"DnsPolicy":null,"EnableServiceLinks":null,"EphemeralContainers":null,"HostAliases":[{"Hostnames":["dex.cluster.local","minio.cluster.local"],"Ip":"172.21.0.2"}],"HostIPC":null,"HostNetwork":null,"HostPID":null,"HostUsers":null,"Hostname":null,"ImagePullSecrets":null,"InitContainers":null,"NodeName":null,"NodeSelector":null,"Os":null,"Overhead":null,"PreemptionPolicy":null,"Priority":null,"PriorityClassName":null,"ReadinessGates":null,"ResourceClaims":null,"RestartPolicy":null,"RuntimeClassName":null,"SchedulerName":null,"SchedulingGates":null,"SecurityContext":null,"ServiceAccount":null,"ServiceAccountName":null,"SetHostnameAsFQDN":null,"ShareProcessNamespace":null,"Subdomain":null,"TerminationGracePeriodSeconds":null,"Tolerations":null,"TopologySpreadConstraints":null,"Volumes":null}'
dlv --listen :2345 --headless --api-version=2 exec /tmp/pod-spec-mutator
```

```bash
telepresence leave pod-spec-mutator-postgresoperator
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
