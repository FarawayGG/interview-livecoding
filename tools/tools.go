//go:build tools
// +build tools

package tools

// https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md
import (
	_ "github.com/envoyproxy/protoc-gen-validate"
	_ "github.com/fullstorydev/grpcui/cmd/grpcui"
	_ "github.com/gojuno/minimock/v3/cmd/minimock"
	_ "github.com/hexdigest/gowrap/cmd/gowrap"
	_ "golang.org/x/tools/cmd/goimports"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
