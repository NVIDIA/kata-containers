# Kata Containers packaging

## Introduction

Kata Containers currently supports packages for many distributions. Tooling to
aid in creating these packages are contained within this repository.

## Build in a container

Kata build artifacts are available within a container image, created by a
[Dockerfile](kata-deploy/Dockerfile). Reference DaemonSets are provided in
[`kata-deploy`](kata-deploy), which make installation of Kata Containers in a
running Kubernetes Cluster very straightforward.

## Build static binaries

See [the static build documentation](static-build).

## Build Kata Containers Kernel

See [the kernel documentation](kernel).

## Build QEMU

See [the QEMU documentation](qemu).

## Create a Kata Containers release

See [the release documentation](release).

## Packaging scripts

See the [scripts documentation](scripts).

## Credits

Kata Containers packaging uses [packagecloud](https://packagecloud.io) for
package hosting.

## KATA Upgrade
There is a daemonset running, likely installed by helm (or Argo). The daemonset will point to the correct image tag with each new release.

To upgrade to a newer release, we should do one of the following:

1. Argo sync obtains a new release by updating the mirror repo
    * Daemonset gets updated
    * Daemonset controller issues a rolling update on each node
    * Each new pod (kata-deploy pod) is launched with the new version (and new args if at all)
    * The new pod discovers that the artifacts do not match (it cannot realize this is an upgrade or a downgrade, and that is a feature)
    * Each pod carries out an update of the artifacts silently (containerd will get restarted behind the scenes)
2. Helm upgrade can be done where helm was used to install kata-deploy
    * Vaules.yaml needs to be updated with the new upgrade (if a helm controller is used then change the suitable CRD)
    * See the same steps as above as daemonset is updated
3. Daemonset is deleted and re-installed (whether through helm purge and helm install, or a manual daemonset delete and create)
    * All old kata-deploy pods will terminate. The side effect is that new kata workloads cannot be scheduled.
    * New image of kata-deploy will be downloaded (hopefully the base layers are unchanged), but the artifacts within (qemu/kata-shim/configs) will be pulled
in anyways
    * For the duration above (which could be several minutes), the cluster cannot take new workloads. This is disruptive but existing workloads do not get aff
ected.
    * New image will get the new pods running and it will NOT be a rolling update. All the new pods will come up simultaneously and the new artifacts will be replaced in parallel.
    * Apart from the downtime of new workloads (old workloads can be deleted), the system will be upgraded silently with the control plane using its eventual consistency model.

