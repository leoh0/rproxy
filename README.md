# RProxy

Simple reverse proxy example

Motivated from [ghproxy](https://github.com/kubernetes/test-infra/blob/master/ghproxy/README.md)

## RProxy Options

* `-log-level`

Log level is one of [panic fatal error warning info debug trace]. (default "debug") 

* `-port`

Port to listen on. (default 8080)

* `-upstream`

Scheme, host, and base path of reverse proxy upstream. (default "https://hooks.slack.com")

``` sh
$ bazel run //cmd/rproxy:rproxy --verbose_failures -- -h
INFO: Analyzed target //cmd/rproxy:rproxy (0 packages loaded, 0 targets configured).
INFO: Found 1 target...
Target //cmd/rproxy:rproxy up-to-date:
  bazel-bin/cmd/rproxy/darwin_amd64_stripped/rproxy
INFO: Elapsed time: 0.193s, Critical Path: 0.09s
INFO: 0 processes.
INFO: Build completed successfully, 1 total action
INFO: Build completed successfully, 1 total action
Usage of /private/var/tmp/_bazel_al/f218ceb5136271f54137115021d86db1/execroot/__main__/bazel-out/darwin-fastbuild/bin/cmd/rproxy/darwin_amd64_stripped/rproxy:
  -log-level string
    	Log level is one of [panic fatal error warning info debug trace]. (default "debug")
  -port int
    	Port to listen on. (default 8080)
  -upstream string
    	Scheme, host, and base path of reverse proxy upstream. (default "https://hooks.slack.com")
```

## How to use this repository

### Command list

``` sh
$ make help
apply-yaml                                         Apply k8s yaml
bazel-target-list                                  List bazel targets
build-image                                        Build docker image
check-yaml                                         Check k8s yaml
help                                               List commands
print-workspace-status                             Print workspace status
push-image                                         Push docker image
run-linux-bin                                      Run linux go binary
run-macos-bin                                      Run macos go binary
```

#### Run binary

``` sh
$ make run-macos-bin
bazel run //cmd/rproxy:rproxy --verbose_failures
INFO: Analyzed target //cmd/rproxy:rproxy (0 packages loaded, 0 targets configured).
INFO: Found 1 target...
Target //cmd/rproxy:rproxy up-to-date:
  bazel-bin/cmd/rproxy/darwin_amd64_stripped/rproxy
INFO: Elapsed time: 0.174s, Critical Path: 0.09s
INFO: 0 processes.
INFO: Build completed successfully, 1 total action
INFO: Build completed successfully, 1 total action
{"component":"rproxy","file":"external/com_github_sirupsen_logrus/logger.go:192","func":"github.com/sirupsen/logrus.(*Logger).Log","level":"info","msg":"Server started.","time":"2020-04-27T13:39:23+09:00"}
```

#### Build and push image

``` sh
$ make build-image
bazel run //cmd/rproxy:bundle --verbose_failures
INFO: Analyzed target //cmd/rproxy:bundle (0 packages loaded, 297 targets configured).
INFO: Found 1 target...
Target //cmd/rproxy:bundle up-to-date (nothing to build)
INFO: Elapsed time: 0.311s, Critical Path: 0.19s
INFO: 0 processes.
INFO: Build completed successfully, 1 total action
INFO: Build completed successfully, 1 total action
Loaded image ID: sha256:9d22ad4fe0b534435ff1d9f28f06ea718bad304f33bf5dff2805f62a0d966800
Loaded image ID: sha256:9d22ad4fe0b534435ff1d9f28f06ea718bad304f33bf5dff2805f62a0d966800
Tagging 9d22ad4fe0b534435ff1d9f28f06ea718bad304f33bf5dff2805f62a0d966800 as docker.io/leoh0/rproxy:v20200427-92e90b2-dirty
Tagging 9d22ad4fe0b534435ff1d9f28f06ea718bad304f33bf5dff2805f62a0d966800 as docker.io/leoh0/rproxy:latest
```

Or just build and push

``` sh
$ make push-image
bazel run //cmd/rproxy:push --verbose_failures
INFO: Analyzed target //cmd/rproxy:push (0 packages loaded, 75 targets configured).
INFO: Found 1 target...
Target //cmd/rproxy:push up-to-date:
  bazel-bin/cmd/rproxy/push
INFO: Elapsed time: 0.290s, Critical Path: 0.09s
INFO: 0 processes.
INFO: Build completed successfully, 1 total action
INFO: Build completed successfully, 1 total action
2020/04/27 13:46:02 Destination {STABLE_DOCKER_REPO}/rproxy:latest was resolved to docker.io/leoh0/rproxy:latest after stamping.
2020/04/27 13:46:02 Destination {STABLE_DOCKER_REPO}/rproxy:{DOCKER_TAG} was resolved to docker.io/leoh0/rproxy:v20200427-92e90b2-dirty after stamping.
2020/04/27 13:46:06 Successfully pushed Docker image to docker.io/leoh0/rproxy:v20200427-92e90b2-dirty
2020/04/27 13:46:07 Successfully pushed Docker image to docker.io/leoh0/rproxy:latest
```

#### Check and apply k8s yaml

``` sh
$ make check-yaml INGRESS_DOMAIN=igress.com PROXY=http://proxy.com
bazel run //cmd/rproxy:k8s --define ingress_domain=igress.com --define proxy=http://proxy.com --verbose_failures
INFO: Build option --define has changed, discarding analysis cache.
INFO: Analyzed target //cmd/rproxy:k8s (0 packages loaded, 7293 targets configured).
INFO: Found 1 target...
Target //cmd/rproxy:k8s up-to-date:
  bazel-bin/cmd/rproxy/k8s
INFO: Elapsed time: 0.393s, Critical Path: 0.14s
INFO: 0 processes.
INFO: Build completed successfully, 3 total actions
INFO: Build completed successfully, 3 total actions
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: rproxy
  name: rproxy
spec:
  replicas: 1
  selector:
    matchLabels:
      app: rproxy
  template:
    metadata:
      labels:
        app: rproxy
    spec:
      containers:
      - env:
        - name: HTTP_PROXY
          value: http://proxy.com
        - name: HTTPS_PROXY
          value: http://proxy.com
        - name: http_proxy
          value: http://proxy.com
        - name: https_proxy
          value: http://proxy.com
        image: index.docker.io/leoh0/rproxy@sha256:f44a7551c18f4456e2049c36cc267c786296f18be872c0338795dc0b70a2d803
        imagePullPolicy: Always
        name: rproxy
        ports:
        - containerPort: 8080
          name: proxy

---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: rproxy
  name: rproxy
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: rproxy

---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  labels:
    app: rproxy
  name: rproxy
spec:
  rules:
  - host: igress.com
    http:
      paths:
      - backend:
          serviceName: rproxy
          servicePort: 8080
        path: /
```

Or just build and apply yaml

``` sh
$ make apply-yaml INGRESS_DOMAIN=igress.com PROXY=http://proxy.com
bazel run //cmd/rproxy:k8s.apply --define ingress_domain=ingress.com --define proxy=http://proxy.com --verbose_failures
INFO: Build option --define has changed, discarding analysis cache.
INFO: Analyzed target //cmd/rproxy:k8s.apply (0 packages loaded, 7298 targets configured).
INFO: Found 1 target...
Target //cmd/rproxy:k8s.apply up-to-date:
  bazel-bin/cmd/rproxy/k8s.apply
INFO: Elapsed time: 0.362s, Critical Path: 0.15s
INFO: 0 processes.
INFO: Build completed successfully, 3 total actions
INFO: Build completed successfully, 3 total actions
$ /usr/local/bin/kubectl --kubeconfig= --cluster=al-cluster --context=al-context --user= --namespace=rproxy apply -f -
deployment.apps/rproxy unchanged
$ /usr/local/bin/kubectl --kubeconfig= --cluster=al-cluster --context=al-context --user= --namespace=rproxy apply -f -
service/rproxy unchanged
$ /usr/local/bin/kubectl --kubeconfig= --cluster=al-cluster --context=al-context --user= --namespace=rproxy apply -f -
ingress.networking.k8s.io/rproxy unchanged
```

#### Check all bazel targets

``` sh
$ make bazel-target-list
bazel query '...'
//cmd/rproxy:services.describe
//cmd/rproxy:rproxy
//cmd/rproxy:push
//cmd/rproxy:k8s.replace
//cmd/rproxy:services.replace
//cmd/rproxy:k8s.delete
//cmd/rproxy:services.delete
//cmd/rproxy:services.reversed
//cmd/rproxy:k8s.create
//cmd/rproxy:services.create
//cmd/rproxy:k8s.apply
//cmd/rproxy:services.apply
//cmd/rproxy:k8s
//cmd/rproxy:services
//cmd/rproxy:ingresses.replace
//cmd/rproxy:ingresses.describe
//cmd/rproxy:ingresses.delete
//cmd/rproxy:ingresses.reversed
//cmd/rproxy:ingresses.create
//cmd/rproxy:ingresses.apply
//cmd/rproxy:ingresses
//cmd/rproxy:deployments.replace
//cmd/rproxy:deployments.describe
//cmd/rproxy:deployments.delete
//cmd/rproxy:deployments.reversed
//cmd/rproxy:deployments.create
//cmd/rproxy:deployments.apply
//cmd/rproxy:deployments
//cmd/rproxy:bundle
//cmd/rproxy:image
//cmd/rproxy:app
//cmd/rproxy:app.binary
//cmd/rproxy:go_default_library
Loading: 0 packages loaded
```

### Change configuration

``` sh
$ make print-workspace-status
hack/print-workspace-status.sh
STABLE_DOCKER_REPO docker.io/leoh0
CLUSTER al-cluster
CONTEXT al-context
NAMESPACE rproxy
DOCKER_TAG v20200427-92e90b2-dirty
```

If you want to update configuration, then `export` variables
``` sh
export DOCKER_REPO_OVERRIDE=my.docker.com/image
```
