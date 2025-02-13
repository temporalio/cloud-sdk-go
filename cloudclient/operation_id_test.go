package cloudclient

import (
	"context"
	"testing"

	"go.temporal.io/cloud-sdk/api/cloudservice/v1"
	"google.golang.org/grpc"
)

func TestSetOperationIDGRPCInterceptorAlreadySet(t *testing.T) {

}

func TestSetOperationIDGRPCInterceptor(t *testing.T) {

	t.Run("AsyncOperationId Not Set", func(t *testing.T) {
		req := &cloudservice.UpdateNamespaceRequest{}
		err := setOperationIDGRPCInterceptor(context.Background(), "method", req, nil, nil,
			func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				return nil
			}, nil)
		if err != nil {
			t.Errorf("SetOperationIDInterceptor() error = %v", err)
		}
		if req.AsyncOperationId == "" {
			t.Errorf("SetOperationIDInterceptor() expected operation ID to be set")
		}
	})

	t.Run("AsyncOperationId Already Set", func(t *testing.T) {
		req := &cloudservice.UpdateNamespaceRequest{
			AsyncOperationId: "already-set",
		}
		err := setOperationIDGRPCInterceptor(context.Background(), "method", req, nil, nil,
			func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				return nil
			}, nil)
		if err != nil {
			t.Errorf("SetOperationIDInterceptor() error = %v", err)
		}
		if req.AsyncOperationId != "already-set" {
			t.Errorf("SetOperationIDInterceptor() expected operation ID to not be set")
		}
	})

	t.Run("Message Without AsyncOperationId", func(t *testing.T) {

		req := &cloudservice.GetNamespaceRequest{}
		err := setOperationIDGRPCInterceptor(context.Background(), "method", req, nil, nil,
			func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				return nil
			}, nil)
		if err != nil {
			t.Errorf("SetOperationIDInterceptor() error = %v", err)
		}
	})
}
