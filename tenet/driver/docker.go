package driver

import (
	"path"

	goDocker "github.com/fsouza/go-dockerclient"
	"github.com/juju/errors"

	"github.com/lingo-reviews/lingo/tenet/driver/docker"
	"github.com/lingo-reviews/lingo/util"
	"github.com/lingo-reviews/tenets/go/dev/tenet/log"
)

// Docker is a tenet driver which runs tenets inside a docker container.
type Docker struct {
	*Base

	dockerClient *goDocker.Client
}

// Pull the image for this tenet from the given registry.
func (d *Docker) Pull(update bool) error {
	dClient, err := d.getDockerClient()
	if err != nil {
		return errors.Trace(err)
	}

	if update || !docker.HaveImage(dClient, d.Name) {
		return docker.PullImage(dClient, d.Name, d.Registry, d.Tag)
	}
	return nil
}

func (d *Docker) getDockerClient() (*goDocker.Client, error) {
	if d.dockerClient == nil {
		dClient, err := util.DockerClient()
		if err != nil {
			return nil, errors.Trace(err)
		}
		d.dockerClient = dClient
	}

	return d.dockerClient, nil
}

// Init the service.
func (d *Docker) Service() (Service, error) {
	log.Print("Docker Service called")
	return docker.NewService(d.Name)
}

// Docker mounts source code under /source/ so we need to prepend this to all
// file names.
func (d *Docker) EditFilename(f string) string {
	return path.Join("/source/", f)
}
