# kube-tmuxp

[![Build Status](https://travis-ci.org/thecasualcoder/kube-tmuxp.svg?branch=master)](https://travis-ci.org/thecasualcoder/kube-tmuxp)

Tool to generate tmuxp configs that help to switch between multiple [Kubernetes](https://kubernetes.io/) contexts safely.

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

```
brew tap thecasualcoder/stable
brew install kube-tmuxp
```

### Manual

```
git clone https://github.com/thecasualcoder/kube-tmuxp.git
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

## Generate kube-tmuxp config file for gcloud

```bash
$ kube-tmuxp gcloud-generate --help
Generates configs for kube-tmuxp based on gcloud account

Usage:
  kube-tmuxp gcloud-generate [flags]

Flags:
      --additional-envs strings   Additional envs to be populated
      --all-projects              Skip confirmation for projects
      --apply                    Directly create the tmuxp configs for selected projects
  -h, --help                     help for gcloud-generate
      --project-ids strings       Comma separated Project IDs to which the configurations need to be fetched

```

#### Examples:
1) Interactively select projects:
```bash
$ kube-tmuxp gcloud-generate
# this will prompt for the projectIDs selection. Type to filter and select using space
# fuzzy search will work
```
2) Specify projectIDs:
```bash
$ kube-tmuxp gcloud-generate --projectIDs project1,project2
```
3) For all projects:
```bash
$ kube-tmuxp gcloud-generate --allProjects
```
4) Use env variables in kube-tmuxp:
> kube-tmuxp provides four envs: `KUBETMUXP_CLUSTER_NAME`, `KUBETMUXP_CLUSTER_LOCATION`, `KUBETMUXP_CLUSTER_IS_REGIONAL`, `GCP_PROJECT_ID`. We can pass additional envs also.
```bash
$ kube-tmuxp gcloud-generate --additionalEnvs 'NEW_KEY=new_value,NEW_ENV=$HOME,KUBE_CONFIG=$HOME/.kube/configs/$KUBETMUXP_CLUSTER_NAME'
# each tmux session will have 7 envs (4 predefined, 3 additionalEnvs passed as argument)
```
5) Directly create the kubeconfigs and tmuxp files (instead of kube-tmuxp config files):
```bash
$ kube-tmuxp gcloud-generate --apply
```
## Start a session

```
tmuxp load my-context-name
```

Now you will be inside a `tmux` session preconfigured with Kubernetes context `my-context-name`.

## Handy bash functions

Use the `bash` functions available [here](https://github.com/arunvelsriram/dotfiles/blob/master/bash_it_custom/plugins/kube-tmuxp.plugin.bash) to switch, kill sessions easily. Special thanks to [@jskswamy](https://github.com/jskswamy) for writing these awesome functions.

## Limitations

* Currently works for Google Kubernetes Engine (GKE) only. However, it can be extended to work with any Kubernetes clusters. Feel free to submit a PR for this.
