package main

import (
	"errors"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
)

//Registry stores auth data and backupRegistry
type Registry struct {
	backupRegistry string
	dockerAuth     crane.Option
}

// NewRegistry creates a new registry
func NewRegistry(backupRegistry, dockerUser, dockerToken string) *Registry {
	return &Registry{
		backupRegistry: backupRegistry,
		dockerAuth:     crane.WithAuth(&authn.Basic{Username: dockerUser, Password: dockerToken})}
}

//AddImageToBackUp adds image to backup and return the name of the new backup image.
//Image will be pushed with the same tag and repository name. Only registry will be changed
//If image already exists it return an error
func (r *Registry) AddImageToBackUp(image string) (string, error) {
	_, repositoryName, imageTag := imageInfo(image)
	backupImageName := r.backupRegistry + "/" + repositoryName + ":" + imageTag

	//check if image already exists
	tags, err := crane.ListTags(r.backupRegistry+"/"+repositoryName, r.dockerAuth)
	for _, tag := range tags {
		if tag == repositoryName {
			return "", errors.New("Image already exists in backup registry as " + backupImageName)
		}
	}

	//create new tag
	tag, err := name.NewTag(backupImageName)
	if err != nil {
		return "", err
	}

	// push to remote registry
	if err := crane.Copy(image, tag.String(), r.dockerAuth); err != nil {
		return "", err
	}

	return backupImageName, nil
}

//IsImageFromBackUp checks if image names is from the registry backup
func (r *Registry) IsImageFromBackUp(imageName string) bool {
	imageRegistry, _, _ := imageInfo(imageName)
	return strings.HasSuffix(imageRegistry, r.backupRegistry)
}

//imageInfo splits image name into registry, repository and tag
func imageInfo(image string) (registry, repository, tag string) {
	ss := strings.Split(image, ":")
	imageWithoutTag := ss[0]
	tag = ss[len(ss)-1]

	ss = strings.Split(imageWithoutTag, "/")
	repository = ss[len(ss)-1]

	if len(ss) > 1 {
		registry = ss[len(ss)-2]
	}
	return
}
