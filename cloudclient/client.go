package cloudclient

import (
	"fmt"

	cloudservice "go.temporal.io/cloud-sdk/api/cloudservice/v1"
	"google.golang.org/grpc"
)

type (
	Client struct {
		conn               *grpc.ClientConn
		cloudServiceClient cloudservice.CloudServiceClient
	}
)

// New creates a client to perform cloud-management operations.
//
// WARNING: Cloud operations client is currently experimental.
//
// The client will not establish a connection to the server until the first call is made.
// The client is safe for concurrent use by multiple goroutines.
// The client must be closed when it is no longer needed to clean up resources.
func New(options Options) (*Client, error) {

	// compute the options provided by the user
	hostPort, grpcDialOptions := options.compute()

	// create a new gRPC client connection
	// note that the grpc.NewClient will not establish a connection to the server until the first call is made
	conn, err := grpc.NewClient(
		hostPort.String(),
		grpcDialOptions...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial `%s`: %w", hostPort.String(), err)
	}

	return &Client{
		conn:               conn,
		cloudServiceClient: cloudservice.NewCloudServiceClient(conn),
	}, nil
}

func (c *Client) CloudService() cloudservice.CloudServiceClient {
	return c.cloudServiceClient
}

func (c *Client) Close() error {
	return c.conn.Close()
}
