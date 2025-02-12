package cloudclient_test

import (
	"context"
	"testing"

	"go.temporal.io/cloud-sdk/api/cloudservice/v1"
	"go.temporal.io/cloud-sdk/cloudclient"
	"google.golang.org/grpc"
)

func TestSetOperationIDInterceptorAlreadySet(t *testing.T) {
	req := &cloudservice.UpdateNamespaceRequest{
		AsyncOperationId: "already-set",
	}
	err := cloudclient.SetOperationIDInterceptor(context.Background(), "method", req, nil, nil,
		func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			return nil
		}, nil)
	if err != nil {
		t.Errorf("SetOperationIDInterceptor() error = %v", err)
	}
	if req.AsyncOperationId != "already-set" {
		t.Errorf("SetOperationIDInterceptor() expected operation ID to not be set")
	}
}

func TestSetOperationIDInterceptor(t *testing.T) {
	req := &cloudservice.UpdateNamespaceRequest{}
	err := cloudclient.SetOperationIDInterceptor(context.Background(), "method", req, nil, nil,
		func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			return nil
		}, nil)
	if err != nil {
		t.Errorf("SetOperationIDInterceptor() error = %v", err)
	}
	if req.AsyncOperationId == "" {
		t.Errorf("SetOperationIDInterceptor() expected operation ID to be set")
	}
}
