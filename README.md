# Image clone controller

This Kubernetes controller watches Deployments and Daemonsets. For each of them, the controller checks if the images are from the backup registry. If they from another registry, we push a copy of it to the backup registry, and we update de kubernetes boject with the new image. The Deployments and DaemonSets in the kube-system namespace are ignored.

## Installation

1. In order to be available to test it, a k8s cluster is required to run it locally (for example minikube).
2. In order to run it, you will need to configure the backup registry, and the user and password for it. You should modify files:
    - /k8sFiles/configmap.yaml: Set backup registry and login user
    - /k8sFiles/secret.yaml: Set login password or token
3. Create a namespace for our controller: `kubectl create namespace operators` 
4. Create all files in k8sFiles directory. `kubectl create -f /k8sFiles`
5. Everything is ready! You can see the logs in the pod created in namespace operators.

For testing, If you want to run just the go files (without running as a deployment in the cluster):

1. You should modify files:
- config/config.yaml: Set backup registry and login user
- secure-config/config.yaml: Set login password or token
2. Run go files: `go run /src/.`


## Demo

In this demo we can see how it works:

[![asciicast](https://asciinema.org/a/399726.svg)](https://asciinema.org/a/399726)

## Considerations

I have used Docker Hub as the back up registry. It should also work for other registries, but I did not test it.
