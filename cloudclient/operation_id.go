package cloudclient

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type (
	requestWithProtoReflectMessage interface {
		ProtoReflect() protoreflect.Message
	}
)

func setOperationIDGRPCInterceptor(
	ctx context.Context,
	method string,
	req interface{}, reply interface{},
	conn *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {

	if msg, ok := req.(requestWithProtoReflectMessage); ok {
		// It is a proto message, check if it has an operation ID field.
		field := msg.ProtoReflect().Descriptor().Fields().ByTextName("async_operation_id")
		if field != nil {
			// The field exists, check if it is empty
			if val := msg.ProtoReflect().Get(field); val.IsValid() && val.String() == "" {
				// The field is empty, set a random value
				msg.ProtoReflect().Set(field, protoreflect.ValueOfString(uuid.NewString()))
			}
		}
	}
	return invoker(ctx, method, req, reply, conn, opts...)
}
