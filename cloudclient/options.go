package cloudclient

import (
	"context"
	"crypto/tls"
	"net/url"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	defaultCloudOpsAPIHostPort = "saas-api.tmprl.cloud:443"
	defaultAPIVersion          = "v0.3.0"

	authorizationHeader           = "Authorization"
	authorizationBearer           = "Bearer"
	temporalCloudAPIVersionHeader = "temporal-cloud-api-version"
)

type (
	Options struct {
		// The hostport to use when connecting to the cloud operations API.
		// If not provided, the default hostport of `saas-api.tmprl.cloud:443` will be used.
		HostPort url.URL

		// Allow the client to connect to the cloud operations API using an insecure connection.
		// This should only be used for testing purposes.
		AllowInsecure bool

		// The API key to use when making requests to the cloud operations API.
		// If not provided, the request will fail to authenticate.
		APIKeyReader APIKeyReader

		// The API version to use when making requests to the cloud operations API.
		// If not provided, the latest API version  will be used.
		APIVersion string

		// Disable the default retry policy.
		// If not provided, the default retry policy will be used.
		// The default retry policy is an exponential backoff with jitter with a maximum of 7 retries for retriable errors.
		// The default retry policy will also set the operations id on the write requests, if not already set.
		// This is useful for ensuring that the write requests are idempotent in the case of a retry.
		DisableRetry bool

		// Add additional gRPC dial options.
		// This can be used to set custom timeouts, interceptors, etc.
		GRPCDialOptions []grpc.DialOption
	}

	APIKeyReader interface {
		// Get the API key to use when making requests to the cloud operations API.
		// If an error is returned, the request will fail.
		// The GetAPIKey function will be called every time a request is made to the cloud operations API.
		GetAPIKey(ctx context.Context) (string, error)
	}

	// StaticAPIKeyReader is an API key reader that always returns the same API key.
	StaticAPIKeyReader struct {
		// The API key to use when making requests to the cloud operations API.
		APIKey string
	}
)

func (r StaticAPIKeyReader) GetAPIKey(ctx context.Context) (string, error) {
	return r.APIKey, nil
}

func (o *Options) compute() (
	hostPort url.URL,
	grpcDialOptions []grpc.DialOption,
) {

	grpcDialOptions = make([]grpc.DialOption, 0, len(o.GRPCDialOptions)+4)
	// set the default host port if not provided
	if o.HostPort.String() == "" {
		defaultHostPort, err := url.Parse(defaultCloudOpsAPIHostPort)
		if err != nil {
			panic(err)
		}
		hostPort = *defaultHostPort
	} else {
		hostPort = o.HostPort
	}

	var transport credentials.TransportCredentials
	// setup the transport
	if o.AllowInsecure {
		transport = insecure.NewCredentials()
	} else {
		transport = credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: o.HostPort.Hostname(),
		})
	}
	grpcDialOptions = append(grpcDialOptions,
		grpc.WithTransportCredentials(transport),
	)

	// setup the api key credentials
	if o.APIKeyReader != nil {
		creds := apikeyCreds{
			reader:                 o.APIKeyReader,
			allowInsecureTransport: o.AllowInsecure,
		}
		grpcDialOptions = append(grpcDialOptions,
			grpc.WithPerRPCCredentials(creds),
		)
	}

	// setup the api version header
	version := o.APIVersion
	if version == "" {
		version = defaultAPIVersion
	}
	grpcDialOptions = append(grpcDialOptions, grpc.WithUnaryInterceptor(
		func(
			ctx context.Context,
			method string,
			req interface{}, reply interface{},
			conn *grpc.ClientConn,
			invoker grpc.UnaryInvoker,
			opts ...grpc.CallOption,
		) error {
			ctx = metadata.AppendToOutgoingContext(ctx, temporalCloudAPIVersionHeader, version)
			return invoker(ctx, method, req, reply, conn, opts...)
		}),
	)

	if !o.DisableRetry {
		// setup the default retry policy
		retryOpts := []retry.CallOption{
			retry.WithBackoff(
				retry.BackoffExponentialWithJitter(500*time.Millisecond, 0.5),
			),
			retry.WithMax(7),
		}
		grpcDialOptions = append(grpcDialOptions,
			grpc.WithChainUnaryInterceptor(
				// set the operation id on the write requests if not already set
				// this is useful for ensuring that the write requests are idempotent in the case of a retry
				SetOperationIDInterceptor,
				// retry the request on retriable errors
				retry.UnaryClientInterceptor(retryOpts...), // retry the request on retriable errors
			),
		)
	}

	grpcDialOptions = append(grpcDialOptions, o.GRPCDialOptions...)
	return
}
