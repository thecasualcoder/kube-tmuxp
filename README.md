# kube-tmuxp

## Setup

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
