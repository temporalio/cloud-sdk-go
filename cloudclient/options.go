package cloudclient

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
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

type Options struct {
	// The API key to use when making requests to the cloud operations API.
	// At least one of APIKey and APIKeyReader must be provided, but not both.
	APIKey string

	// The API key reader to dynamically retrieve apikey to use when making requests to the cloud operations API.
	// At least one of APIKey and APIKeyReader must be provided, but not both.
	APIKeyReader APIKeyReader

	// The hostport to use when connecting to the cloud operations API.
	// If not provided, the default hostport of `saas-api.tmprl.cloud:443` will be used.
	HostPort url.URL

	// Allow the client to connect to the cloud operations API using an insecure connection.
	// This should only be used for testing purposes.
	AllowInsecure bool

	// The TLS configuration to use when connecting to the cloud operations API.
	// If not provided, a default TLS configuration will be used.
	// Will be ignored if AllowInsecure is set to true.
	TLSConfig tls.Config

	// The API version to use when making requests to the cloud operations API.
	// If not provided, the latest API version  will be used.
	APIVersion string

	// Enable the default retry policy.
	// The default retry policy is an exponential backoff with jitter with a maximum of 7 retries for retriable errors.
	EnableRetry bool

	// Add additional gRPC dial options.
	// This can be used to set custom timeouts, interceptors, etc.
	GRPCDialOptions []grpc.DialOption
}

type APIKeyReader interface {
	// Get the API key to use when making requests to the cloud operations API.
	// If an error is returned, the request will fail.
	// The GetAPIKey function will be called every time a request is made to the cloud operations API.
	GetAPIKey(ctx context.Context) (string, error)
}

type staticAPIKeyReader struct {
	// The API key to use when making requests to the cloud operations API.
	APIKey string
}

func (r staticAPIKeyReader) GetAPIKey(ctx context.Context) (string, error) {
	return r.APIKey, nil
}

func (o *Options) compute() (
	hostPort url.URL,
	grpcDialOptions []grpc.DialOption,
	err error,
) {

	grpcDialOptions = make([]grpc.DialOption, 0, len(o.GRPCDialOptions)+4)
	// set the default host port if not provided
	if o.HostPort.String() == "" {
		var defaultHostPort *url.URL
		defaultHostPort, err = url.Parse(defaultCloudOpsAPIHostPort)
		if err != nil {
			return url.URL{}, nil, fmt.Errorf("failed to parse default host port: %w", err)
		}
		hostPort = *defaultHostPort
	} else {
		hostPort = o.HostPort
	}

	var transport credentials.TransportCredentials
	// setup the transport
	if o.AllowInsecure {
		// allow insecure transport
		transport = insecure.NewCredentials()
	} else {
		// use the provided tls config, or the zero value if not provided
		transport = credentials.NewTLS(&o.TLSConfig)
	}
	grpcDialOptions = append(grpcDialOptions,
		grpc.WithTransportCredentials(transport),
	)

	if o.APIKey != "" && o.APIKeyReader != nil {
		return url.URL{}, nil, errors.New("only one of APIKey and APIKeyReader can be provided")
	}
	// setup the api key credentials
	creds := apikeyCreds{
		allowInsecureTransport: o.AllowInsecure,
	}
	if o.APIKey != "" {
		creds.reader = staticAPIKeyReader{APIKey: o.APIKey}
	} else if o.APIKeyReader != nil {
		creds.reader = o.APIKeyReader
	}
	if creds.reader == nil {
		return url.URL{}, nil, errors.New("either APIKey or APIKeyReader must be provided")
	} else {
		grpcDialOptions = append(grpcDialOptions,
			grpc.WithPerRPCCredentials(creds),
		)
	}

	// setup the api version header
	version := o.APIVersion
	if version == "" {
		version = defaultAPIVersion
	}
	grpcDialOptions = append(grpcDialOptions, grpc.WithChainUnaryInterceptor(
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
		},
	))

	if o.EnableRetry {
		// setup the default retry policy
		retryOpts := []retry.CallOption{
			retry.WithBackoff(
				retry.BackoffExponentialWithJitter(500*time.Millisecond, 0.5),
			),
			retry.WithMax(7),
		}
		grpcDialOptions = append(grpcDialOptions, grpc.WithChainUnaryInterceptor(
			// retry the request on retriable errors
			retry.UnaryClientInterceptor(retryOpts...),
		))
	}

	grpcDialOptions = append(grpcDialOptions, o.GRPCDialOptions...)
	return hostPort, grpcDialOptions, nil
}
