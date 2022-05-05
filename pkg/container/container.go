// Package container package container
package container

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"trade_agent/pkg/log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// RunContainer RunContainer
func RunContainer(imageName string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Get().Panic(err)
	}
	// imageName = "bfirsh/reticulate-splines"
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		log.Get().Panic(err)
	}
	defer func() {
		_ = out.Close()
	}()
	_, _ = io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{Image: imageName}, nil, nil, nil, "")
	if err != nil {
		log.Get().Panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Get().Panic(err)
	}

	fmt.Println(resp.ID)
}

// ListAllContainers ListAllContainers
func ListAllContainers() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Get().Panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Get().Panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID)
	}
}

// StopAllContainers StopAllContainers
func StopAllContainers() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Get().Panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Get().Panic(err)
	}

	for _, container := range containers {
		fmt.Print("Stopping container ", container.ID[:10], "... ")
		if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
			log.Get().Panic(err)
		}
		fmt.Println("Success")
	}
}

// ListAllImages ListAllImages
func ListAllImages() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Get().Panic(err)
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		log.Get().Panic(err)
	}

	for _, image := range images {
		fmt.Println(image.ID)
	}
}

// PullImage PullImage
func PullImage() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Get().Panic(err)
	}

	out, err := cli.ImagePull(ctx, "gitlab.tocraw.com:5050/root/sinopac_mq_srv:d10b5c4d", types.ImagePullOptions{})
	if err != nil {
		log.Get().Panic(err)
	}

	defer func() {
		_ = out.Close()
	}()

	_, _ = io.Copy(os.Stdout, out)
}

// PullImageWithAuth PullImageWithAuth
func PullImageWithAuth() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Get().Panic(err)
	}

	authConfig := types.AuthConfig{
		Username: "root",
		Password: "notriq-9Korzi-xuwduj",
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		log.Get().Panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	out, err := cli.ImagePull(ctx, "gitlab.tocraw.com:5050/root/sinopac_mq_srv:d10b5c4d", types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		log.Get().Panic(err)
	}

	defer func() {
		_ = out.Close()
	}()
	_, _ = io.Copy(os.Stdout, out)
	// body, err := ioutil.ReadAll(out)
	// if err != nil {
	// 	log.Get().Panic(err)
	// }
	// for _, v := range strings.Split(string(body), "\n") {

	// }
}

// PullStatus PullStatus
type PullStatus struct {
	Status         string `json:"status"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
	Progress string `json:"progress"`
	ID       string `json:"id"`
}
