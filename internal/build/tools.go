package build

import (
	// TODO: Move these to using -tool flag once go 1.24 is released (https://tip.golang.org/doc/go1.24#tools)
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
