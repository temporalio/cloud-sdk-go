package cloudclient

import (
	"fmt"

	cloudservice "go.temporal.io/cloud-sdk/api/cloudservice/v1"
	"google.golang.org/grpc"
)

type (
	// Client is the client for cloud operations.
	//
	// WARNING: Cloud operations client is currently experimental.
	Client interface {
		// CloudService provides access to the underlying gRPC service.
		CloudService() cloudservice.CloudServiceClient

		// Close client and clean up underlying resources.
		Close()
	}

	client struct {
		conn               *grpc.ClientConn
		cloudServiceClient cloudservice.CloudServiceClient
	}
)

// New creates a client to perform cloud-management operations.
//
// WARNING: Cloud operations client is currently experimental.
func New(options ...Option) (Client, error) {
	return newClient(options)
}

func newClient(options []Option) (Client, error) {

	// compute the options provided by the user
	opts := computeOptions(options)

	// create a new gRPC client connection
	// note that the grpc.NewClient will not establish a connection to the server until the first call is made
	conn, err := grpc.NewClient(
		opts.hostPort.String(),
		opts.grpcDialOptions...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial `%s`: %v", opts.hostPort.String(), err)
	}

	return &client{
		conn:               conn,
		cloudServiceClient: cloudservice.NewCloudServiceClient(conn),
	}, nil
}

func (c *client) CloudService() cloudservice.CloudServiceClient {
	return c.cloudServiceClient
}

func (c *client) Close() {
	c.conn.Close()
}
