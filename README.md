# kube-tmuxp

## Setup

* Setup `virtualenv`

  ```
  virtualenv --python=python3.6 venv
  source venv/bin/activate
  pip install -r requirements.txt
  ```

* Copy the sample config ([sample.yaml](./config/sample.yaml)) and modify it
* Generate kube configs

```
python kube-tmuxp.py ./config/my-config.yaml
```
