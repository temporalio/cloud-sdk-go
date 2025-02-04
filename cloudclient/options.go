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
	DefaultCloudOpsAPIHostPort = "saas-api.tmprl.cloud:443"
	DefaultAPIVersion          = "v0.3.0"

	AuthorizationHeader           = "Authorization"
	AuthorizationBearer           = "Bearer"
	TemporalCloudAPIVersionHeader = "temporal-cloud-api-version"
)

type (
	clientOptions struct {
		hostPort          url.URL
		allowInsecure     bool
		perRPCCredentials perRPCCredentials
		apiVersion        string
		disableRetry      bool
		retryOpts         []retry.CallOption
		grpcDialOptions   []grpc.DialOption
	}

	Option func(*clientOptions)

	perRPCCredentials interface {
		credentials.PerRPCCredentials
		allowInsecure()
	}
)

// The hostport to use when connecting to the cloud operations API.
// If not provided, the default hostport of `saas-api.tmprl.cloud:443` will be used.
func WithHostPort(hostPort url.URL) Option {
	return func(o *clientOptions) {
		o.hostPort = hostPort
	}
}

// Allow the client to connect to the cloud operations API using an insecure connection.
// This should only be used for testing purposes.
func WithAllowInsecure() Option {
	return func(o *clientOptions) {
		o.allowInsecure = true
	}
}

// Use the provided API key to authenticate with the cloud operations API.
// The API key will be sent as a bearer token in the `Authorization` header.
func WithAPIKey(getAPIKeyFn func() (string, error)) Option {
	return func(o *clientOptions) {
		o.perRPCCredentials = &apikeyCreds{
			getAPIKeyFn: getAPIKeyFn,
		}
	}
}

// Add additional gRPC dial options.
// This can be used to set custom timeouts, interceptors, etc.
func WithGRPCDialOptions(opts ...grpc.DialOption) Option {
	return func(o *clientOptions) {
		o.grpcDialOptions = append(o.grpcDialOptions, opts...)
	}
}

// Set the API version to use when making requests to the cloud operations API.
// If not provided, the default API version of `v0.3.0` will be used.
func WithAPIVersion(version string) Option {
	return func(o *clientOptions) {
		o.apiVersion = version
	}
}

// Disable the retry policy for the client.
// By default, the client will retry requests with exponential backoff up to 7 times.
func WithDisableRetry() Option {
	return func(o *clientOptions) {
		o.disableRetry = true
	}
}

func computeOptions(opts []Option) *clientOptions {

	options := &clientOptions{}
	// apply user-provided options
	for _, opt := range opts {
		opt(options)
	}

	grpcDialOptions := make([]grpc.DialOption, 0, len(options.grpcDialOptions)+4)

	// set the default host port if not provided
	if options.hostPort.String() == "" {
		defaultHostPort, err := url.Parse(DefaultCloudOpsAPIHostPort)
		if err != nil {
			panic(err)
		}
		options.hostPort = *defaultHostPort
	}

	var transport credentials.TransportCredentials
	// setup the transport
	if options.allowInsecure {
		if options.perRPCCredentials != nil {
			options.perRPCCredentials.allowInsecure()
		}
		transport = insecure.NewCredentials()
	} else {
		transport = credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: options.hostPort.Hostname(),
		})
	}
	if transport != nil {
		grpcDialOptions = append(grpcDialOptions,
			grpc.WithTransportCredentials(transport),
		)
	}

	// setup the auth credentials
	if options.perRPCCredentials != nil {
		grpcDialOptions = append(grpcDialOptions,
			grpc.WithPerRPCCredentials(options.perRPCCredentials),
		)
	}

	// setup the api version header
	version := options.apiVersion
	if version == "" {
		version = DefaultAPIVersion
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
			ctx = metadata.AppendToOutgoingContext(ctx, TemporalCloudAPIVersionHeader, version)
			return invoker(ctx, method, req, reply, conn, opts...)
		}),
	)

	if !options.disableRetry {
		// setup the default retry policy
		retryOpts := []retry.CallOption{
			retry.WithBackoff(
				retry.BackoffExponentialWithJitter(500*time.Millisecond, 0.5),
			),
			retry.WithMax(7),
		}
		grpcDialOptions = append(grpcDialOptions,
			grpc.WithChainUnaryInterceptor(retry.UnaryClientInterceptor(retryOpts...)),
		)
	}

	grpcDialOptions = append(grpcDialOptions, options.grpcDialOptions...)
	options.grpcDialOptions = grpcDialOptions

	return options
}
