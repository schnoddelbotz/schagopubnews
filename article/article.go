package article

import (
	"time"
)

// Article describes structure of our FireStore/FireStore documents
type Article struct {
	// Status tracks VM status
	Status string
	// VMID holds name of VM
	VMID string
	// InstanceID holds ID of GCE instance, only known after creation. It's a uint64, but datastore no supporty.
	InstanceID string
	// ManagementToken is known to VM itself, so it can request clean self-destruction
	ManagementToken string
	// CreatedAt ...
	CreatedAt time.Time

	// DockerExitCode from docker run command on VM
	DockerExitCode int
	// DockerContainerId stores container ID on VM, to enable StackDriver log filtering
	DockerContainerId string
}

const (
	// ArticleStatusCreated : initial state after creating new FireStore entry
	ArticleStatusCreated = "created"
	// ArticleStatusBooted : VM booted, tries to run docker pull on our image now
	ArticleStatusBooted = "booted"
	// ArticleStatusRunning : client has pulled image and docker command should be running now
	ArticleStatusRunning = "running"
	// ArticleStatusDone : Container has exited, VM destruction requested
	ArticleStatusDone = "done"
)
