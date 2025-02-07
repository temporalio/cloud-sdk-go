//go:build tools

package tools

import (
	// TODO: Move to using -tool flag once go 1.24 is released (https://tip.golang.org/doc/go1.24#tools)
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
)
