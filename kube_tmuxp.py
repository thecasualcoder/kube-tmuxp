import os
import sys
import errno
import yaml
from subprocess import call

kubeconfigs_dir = os.path.expanduser('~/.kube/configs')

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

def add_context(kubeconfig_filename, project_name, cluster_name, zone):
  cmd = "KUBECONFIG={0}/{1} gcloud container clusters get-credentials {2} --zone {3} --project {4}".format(kubeconfigs_dir, kubeconfig_filename, cluster_name, zone, project_name)
  execute(cmd)

def rename_context(new_context_name, project_name, cluster_name, zone):
  cmd = "KUBECONFIG={0}/{1} kubectl config rename-context gke_{2}_{3}_{4} {1}".format(kubeconfigs_dir, new_context_name, project_name, zone, cluster_name)
  execute(cmd)

def process(config_file):
  with open(config_file, 'r') as stream:
    try:
      configs = yaml.load(stream)
    except yaml.YAMLError as e:
      sys.stderr.write("Failed to load config: \n{0}\n".format(e))
      raise

  for config in configs:
    for cluster in config['clusters']:
      print("\n>>>>> Running for context: {0}\n".format(cluster['context']))
      delete_context(cluster['context'])
      add_context(cluster['context'], config['project'], cluster['name'], cluster['zone'])
      rename_context(cluster['context'], config['project'], cluster['name'], cluster['zone'])

def init():
  os.makedirs(kubeconfigs_dir, exist_ok=True)

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
