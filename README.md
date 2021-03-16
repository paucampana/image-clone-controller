# Image clone controller

This Kubernetes controller watches Deployments and Daemonsets. For each of them, the controller checks if the images are from the backup registry. If they from another registry, we push a copy of it to the backup registry, and we update de kubernetes boject with the new image. The Deployments and DaemonSets in the kube-system namespace are ignored.

## Installation

The [k8sFiles] folder contains the necessary configuration files. Adjust them to your
taste and apply in the given order:

    kubectl apply -f /k8sFiles

