# kube-tmuxp

Easier way to work with multiple [Kubernetes](https://kubernetes.io/) clusters.

## Introduction

When working with multiple Kubernetes clusters its painful to switch context using [`kubectl`](https://github.com/kubernetes/kubernetes/tree/master/cmd/kubectl) or [`kubectx`](https://github.com/ahmetb/kubectx). There are also possibilities of making unintentional changes.

`kube-tmuxp` solves this by using one preconfigured `tmux` session per Kubernetes cluster. Each `tmux` session contains only one Kubernetes context thus preventing accidental context switching inside a session. Contexts can be switched by switching `tmux` sessions. For example: `[tmux prefix] + S`.

Given a config similar to [config.sample.yaml](./config.sample.yaml), `kube-tmuxp` generates:

* kube config (Kubernetes context) for each Kubernetes cluster under `~/.kube/configs`
* `tmuxp` config for each Kubernetes cluster under `~/.tmuxp`

The generated `tmuxp` configs are used to start preconfigured `tmux` sessions.

## How to use?

### Prerequisites

* [gcloud](https://cloud.google.com/sdk/)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [tmux](https://github.com/tmux/tmux)
* [tmuxp](https://github.com/tmux-python/tmuxp)

### Generate kube configs and `tmuxp` configs

* Setup `virtualenv`

  ```
  virtualenv --python=python3.6 venv
  source venv/bin/activate
  pip install -r requirements.txt
  ```

* Copy the sample config ([config.sample.yaml](./config.sample.yaml))

  ```
  cp config.sample.yaml config.yaml
  ```

* Add your clusters to the copied config
* Generate kube configs and tmuxp configs

```
python kube-tmuxp.py config.yaml
```

### Start a session

```
tmuxp load my-context-name
```

Now you will be inside a `tmux` session preconfigured with Kubernetes context `my-context-name`.

## Limitations

* Generates kube configs for Kubernetes clusters (GKE) on Google Cloud Platform (GCP) only
