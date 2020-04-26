load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@io_bazel_rules_docker//go:image.bzl", "go_image")
load("@io_bazel_rules_docker//container:image.bzl", "container_image")
load("@io_bazel_rules_docker//container:bundle.bzl", "container_bundle")
load("@io_bazel_rules_docker//contrib:push-all.bzl", "container_push")
load("@io_bazel_rules_k8s//k8s:object.bzl", "k8s_object")
load("@io_bazel_rules_k8s//k8s:objects.bzl", "k8s_objects")
load("@k8s_deploy//:defaults.bzl", "k8s_deploy")
load("@k8s_service//:defaults.bzl", "k8s_service")
load("@k8s_ingress//:defaults.bzl", "k8s_ingress")

package(default_visibility = ["//visibility:public"])

NAME = "rproxy"

APP_NAME = "app"

BUNDLE_NAME = "bundle"

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/leoh0/rproxy/cmd/rproxy",
    visibility = ["//visibility:private"],
    deps = [
        "@com_github_sirupsen_logrus//:go_default_library",
    ],
)

go_binary(
    name = NAME,
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_image(
    name = APP_NAME,
    base = "@go_image_base//image",
    embed = [":go_default_library"],
    goarch = "amd64",
    goos = "linux",
    pure = "on",
)

container_image(
    name = "image",
    base = ":" + APP_NAME,
    stamp = False,
    visibility = ["//visibility:public"],
)

container_bundle(
    name = BUNDLE_NAME,
    images = {
        "{STABLE_DOCKER_REPO}/rproxy:{DOCKER_TAG}": ":image",
        "{STABLE_DOCKER_REPO}/rproxy:latest": ":image",
    },
)

container_push(
    name = "push",
    bundle = ":" + BUNDLE_NAME,
    format = "Docker",
)

k8s_deploy(
    name = "deployments",
    images = {
        "docker.io/leoh0/rproxy:latest": ":image",
    },
    substitutions = {
        "%{proxy}": "$(proxy)",
    },
    template = ":manifest/deployment.yaml",
)

k8s_service(
    name = "services",

    # A template of a Kubernetes Deployment object yaml.
    template = "manifest/service.yaml",
)

k8s_ingress(
    name = "ingresses",
    substitutions = {
        "%{ingress}": "$(ingress_domain)",
    },

    # A template of a Kubernetes Deployment object yaml.
    template = "manifest/ingress.yaml",
)

k8s_objects(
    name = "k8s",
    objects = [
        ":deployments",
        ":services",
        ":ingresses",
    ],
)
