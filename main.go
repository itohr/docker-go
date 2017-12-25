package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func CreateContainer() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	imageName := "edgexfoundry/docker-edgex-volume"

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	config := container.Config{
		Image:    imageName,
		Hostname: "edgex-files",
		Volumes: map[string]struct{}{
			"/data/db":       struct{}{},
			"/edgex/logs":    struct{}{},
			"/consul/config": struct{}{},
			"/consul/data":   struct{}{},
		},
	}
	netConfig := network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"edgex-network": {},
		},
	}

	resp, err := cli.ContainerCreate(ctx, &config, nil, &netConfig, "edgex-files")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)
}

func CreateNetwork() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	opts := types.NetworkCreate{
		Driver: "bridge",
	}

	resp, err := cli.NetworkCreate(ctx, "edgex-network", opts)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func main() {
	CreateNetwork()
	CreateContainer()
}
