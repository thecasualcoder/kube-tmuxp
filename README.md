# kube-tmuxp

[![Build Status](https://travis-ci.org/arunvelsriram/kube-tmuxp.svg?branch=master)](https://travis-ci.org/arunvelsriram/kube-tmuxp)

Easier way to work with multiple [Kubernetes](https://kubernetes.io/) clusters.

## Introduction

When working with multiple Kubernetes clusters its painful to switch context using [`kubectl`](https://github.com/kubernetes/kubernetes/tree/master/cmd/kubectl) or [`kubectx`](https://github.com/ahmetb/kubectx). There are also possibilities of making unintentional changes.

`kube-tmuxp` solves this by using one preconfigured `tmux` session per Kubernetes cluster. Each `tmux` session contains only one Kubernetes context thus preventing accidental context switching inside a session. Contexts can be switched by switching `tmux` sessions. For example: `[tmux prefix] + S`.

Given a config similar to [config.sample.yaml](./config.sample.yaml), `kube-tmuxp` generates:

* kube config (Kubernetes context) for each Kubernetes cluster under `~/.kube/configs`
* `tmuxp` config for each Kubernetes cluster under `~/.tmuxp`

The generated `tmuxp` configs can be used to start preconfigured `tmux` sessions.

## Prerequisites

* [gcloud](https://cloud.google.com/sdk/)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [tmux](https://github.com/tmux/tmux)
* [tmuxp](https://github.com/tmux-python/tmuxp)

## Install

### Homebrew

To be updated

### Manual

```
git clone https://github.com/arunvelsriram/kube-tmuxp.git
cd kube-tmuxp
make build
cp ./out/kube-tmuxp /usr/local/bin/kube-tmuxp
```

## Generate kubeconfigs and `tmuxp` configs

* Copy the sample config ([config.sample.yaml](./config.sample.yaml))

  ```
  cp config.sample.yaml ~/.kube-tmuxp.yaml
  ```

* Add your projects and clusters to the copied config
* Generate kubeconfigs and tmuxp configs

```
kube-tmuxp gen
```

Default config path is `$HOME/.kube-tmuxp.yaml`. If you are using a different path, then use the `--config` flag to specify that path. Refer `kube-tmuxp --help` for more details.

## Start a session

```
tmuxp load my-context-name
```

Now you will be inside a `tmux` session preconfigured with Kubernetes context `my-context-name`.

## Limitations

* Currently works for Google Kubernetes Engine (GKE) only. However, it can be extended to work with any Kubernetes clusters. Feel free to submit a PR for this.
