package cloudclient

import (
	"context"
	"reflect"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type (
	requestWithAsyncOperationID interface {
		GetAsyncOperationId() string
	}
)

func SetOperationIDInterceptor(
	ctx context.Context,
	method string,
	req interface{}, reply interface{},
	conn *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	if reqWithOperationID, ok := req.(requestWithAsyncOperationID); ok && reqWithOperationID.GetAsyncOperationId() == "" {
		// The request does not have an operation ID set, set a random one
		// This is best effort, if the request does not have the field, it will be ignored
		req = setAsyncOperationID(uuid.NewString(), req)
	}
	return invoker(ctx, method, req, reply, conn, opts...)
}
func setAsyncOperationID(operationID string, request interface{}) interface{} {
	// Get the reflect.Value of the request object
	objValue := reflect.ValueOf(request)
	// Check if the value is addressable (can be set)
	if objValue.Kind() == reflect.Ptr && !objValue.IsNil() {
		// Get the element pointed to by the pointer
		elemValue := objValue.Elem()
		// Get the AsyncOperationId field by name
		nameField := elemValue.FieldByName("AsyncOperationId")
		// Check if the field exists and can be set
		if nameField.IsValid() && nameField.CanSet() {
			nameField.SetString(operationID)
		}
	}
	return request
}
