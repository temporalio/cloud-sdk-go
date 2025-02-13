package build

//go:generate go install -modfile go.mod github.com/bufbuild/buf/cmd/buf
//go:generate go install -modfile go.mod google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go install -modfile go.mod google.golang.org/grpc/cmd/protoc-gen-go-grpc

import (
	// TODO: Move these to using -tool flag once go 1.24 is released (https://tip.golang.org/doc/go1.24#tools)
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
