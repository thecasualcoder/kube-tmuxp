import os
import sys
import errno
import yaml
from subprocess import call
from jinja2 import FileSystemLoader, Environment

kubeconfigs_dir = os.path.expanduser('~/.kube/configs')
tmuxp_dir = os.path.expanduser('~/.tmuxp')

def execute(cmd):
  print("Executing command: {0}".format(cmd))
  try:
    retcode = call(cmd, shell=True)
    if retcode != 0:
      raise Exception("Command failed with exit status: {0}".format(retcode))
  except OSError as e:
    sys.stderr.write("Failed to execute command: \n{0}\n".format(e))
    raise

def delete_context(kubeconfig_filename):
  fullpath = os.path.join(kubeconfigs_dir, kubeconfig_filename)
  print("Removing: {0}".format(fullpath))
  try:
    os.remove(fullpath)
  except OSError as e:
    if e.errno != errno.ENOENT:
      raise

def add_context(kubeconfig_filename, project_name, cluster_name, location, regional):
  kubeconfig = os.path.join(kubeconfigs_dir, kubeconfig_filename)
  if regional:
    cmd = 'CLOUDSDK_CONTAINER_USE_V1_API_CLIENT=false CLOUDSDK_CONTAINER_USE_V1_API=false KUBECONFIG={kubeconfig} gcloud beta container clusters get-credentials {cluster_name} --region {location} --project {project_name}'.format(kubeconfig=kubeconfig, cluster_name=cluster_name, location=location, project_name=project_name)
  else:
    cmd = "KUBECONFIG={kubeconfig} gcloud container clusters get-credentials {cluster_name} --zone {location} --project {project_name}".format(kubeconfig=kubeconfig, cluster_name=cluster_name, location=location, project_name=project_name)

  execute(cmd)

def rename_context(new_context_name, project_name, cluster_name, zone):
  cmd = "KUBECONFIG={0}/{1} kubectl config rename-context gke_{2}_{3}_{4} {1}".format(kubeconfigs_dir, new_context_name, project_name, zone, cluster_name)
  execute(cmd)

def generate_tmuxp_config(context_name, extra_envs):
    context_file = os.path.join(kubeconfigs_dir, context_name)
    tmuxp_config_file = os.path.join(tmuxp_dir, "{0}.yaml".format(context_name))
    template_env = Environment(loader=FileSystemLoader(searchpath='./templates'))
    template = template_env.get_template('tmuxp-config.yaml.j2')
    tmuxp_config = template.render(kubeconfig=context_file, session_name=context_name, extra_envs=extra_envs)
    with open(tmuxp_config_file, 'w') as f:
      f.write(tmuxp_config)
    print("\ntmuxp config generated: {0}".format(tmuxp_config_file))

def process(config_file):
  with open(config_file, 'r') as f:
    try:
      configs = yaml.load(f)
    except yaml.YAMLError as e:
      sys.stderr.write("Failed to load config: \n{0}\n".format(e))
      raise

  for config in configs:
    for cluster in config['clusters']:
      print("\n>>>>> Running for context: {0}\n".format(cluster['context']))
      if 'region' in cluster:
        regional = True
        location = cluster['region']
      else:
        regional = False
        location = cluster['zone']

      delete_context(cluster['context'])
      add_context(cluster['context'], config['project'], cluster['name'], location, regional)
      rename_context(cluster['context'], config['project'], cluster['name'], location)
      generate_tmuxp_config(cluster['context'], cluster.get('extra_envs', {}))

def init():
  os.makedirs(kubeconfigs_dir, exist_ok=True)
  os.makedirs(tmuxp_dir, exist_ok=True)

if __name__ == '__main__':
  args = sys.argv[1:]
  args_count = len(args)
  if args_count != 1:
    sys.stderr.write("Wrong number of arguments\n")
    raise
  else:
    config_file = args[0]
    init()
    process(config_file)
