/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const deploymentType string = "Deployment"
const daemonsetType string = "DaemonSet"

// reconcile structure
type reconcileBackup struct {
	// client can be used to retrieve objects from the APIServer.
	client   client.Client
	registry *Registry
	k8sType  string
}

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &reconcileBackup{}

func (r *reconcileBackup) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	// set up a convenient log object so we don't have to type request over and over again
	log := log.FromContext(ctx)

	if request.NamespacedName.Namespace == "kube-system" {
		log.Info("Ignoring all objects from kube-system namespace")
		return reconcile.Result{}, nil
	}

	var result reconcile.Result
	var err error

	switch r.k8sType {
	case deploymentType:
		result, err = reconcileDeployment(ctx, request, r)
	case daemonsetType:
		result, err = reconcileDaemonSet(ctx, request, r)
	default:
		log.Error(nil, "Reconcile does not support type "+r.k8sType)
	}

	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find object of type "+r.k8sType)
		return result, nil
	}
	return result, err

}

//reconcileDeployment updates the deployment images. if it is necessary it will add the image to backup registry
func reconcileDeployment(ctx context.Context, request reconcile.Request, r *reconcileBackup) (reconcile.Result, error) {
	log := log.FromContext(ctx)

	obj := &appsv1.Deployment{}
	// Fetch the object from the cache
	err := r.client.Get(ctx, request.NamespacedName, obj)
	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find object of type "+r.k8sType)
		return reconcile.Result{}, nil
	}

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch object: %+v", err)
	}

	//checking images from object
	objectChanged := false

	for i, initContainer := range obj.Spec.Template.Spec.InitContainers {
		if !r.registry.IsImageFromBackUp(initContainer.Image) {
			objectChanged = true
			log.Info("Image " + initContainer.Image + " does not exists in backup")
			newImageName, err := r.registry.AddImageToBackUp(initContainer.Image)
			if err != nil {
				log.Error(err, "Could not add image from init container to backup")
				return reconcile.Result{}, nil
			}
			log.Info("Image " + initContainer.Image + " pushed to backup as " + newImageName)
			obj.Spec.Template.Spec.InitContainers[i].Image = newImageName
		}
	}

	for i, container := range obj.Spec.Template.Spec.Containers {
		if !r.registry.IsImageFromBackUp(container.Image) {
			objectChanged = true
			log.Info("Image " + container.Image + " does not exists in backup")
			newImageName, err := r.registry.AddImageToBackUp(container.Image)
			if err != nil {
				log.Error(err, "Could not add image from container to backup")
				return reconcile.Result{}, nil
			}
			log.Info("Image " + container.Image + " pushed to backup as " + newImageName)
			obj.Spec.Template.Spec.Containers[i].Image = newImageName
		}
	}

	if !objectChanged {
		//It is not necessary to do the update
		log.Info("deployment does not need any modification")
		return reconcile.Result{}, nil
	}

	// Update the Deployment
	err = r.client.Update(ctx, obj)
	if err != nil {
		log.Info("deployment updated")
		return reconcile.Result{}, fmt.Errorf("could not write Deployment: %+v", err)
	}
	return reconcile.Result{}, nil
}

//reconcileDaemonSet updates the daemonset images. if it is necessary it will add the image to backup registry
func reconcileDaemonSet(ctx context.Context, request reconcile.Request, r *reconcileBackup) (reconcile.Result, error) {
	log := log.FromContext(ctx)

	obj := &appsv1.DaemonSet{}
	// Fetch the object from the cache
	err := r.client.Get(ctx, request.NamespacedName, obj)
	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find object of type "+r.k8sType)
		return reconcile.Result{}, nil
	}

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch object: %+v", err)
	}

	//checking images from object
	objectChanged := false

	for i, initContainer := range obj.Spec.Template.Spec.InitContainers {
		if !r.registry.IsImageFromBackUp(initContainer.Image) {
			objectChanged = true
			log.Info("Image " + initContainer.Image + " does not exists in backup")
			newImageName, err := r.registry.AddImageToBackUp(initContainer.Image)
			if err != nil {
				log.Error(err, "Could not add image from init container to backup")
				return reconcile.Result{}, nil
			}
			log.Info("Image " + initContainer.Image + " pushed to backup as " + newImageName)
			obj.Spec.Template.Spec.InitContainers[i].Image = newImageName
		}
	}

	for i, container := range obj.Spec.Template.Spec.Containers {
		if !r.registry.IsImageFromBackUp(container.Image) {
			objectChanged = true
			log.Info("Image " + container.Image + " does not exists in backup")
			newImageName, err := r.registry.AddImageToBackUp(container.Image)
			if err != nil {
				log.Error(err, "Could not add image from container to backup")
				return reconcile.Result{}, nil
			}
			log.Info("Image " + container.Image + " pushed to backup as " + newImageName)
			obj.Spec.Template.Spec.Containers[i].Image = newImageName
		}
	}

	if !objectChanged {
		//It is not necessary to do the update
		log.Info("daemonset does not need any modification")
		return reconcile.Result{}, nil
	}

	// Update the Daemonset
	err = r.client.Update(ctx, obj)
	if err != nil {
		log.Info("daemonset updated")
		return reconcile.Result{}, fmt.Errorf("could not write Daemonset: %+v", err)
	}

	return reconcile.Result{}, nil
}
