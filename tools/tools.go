// +build tools

package tools

import (
	_ "github.com/fullstorydev/grpcurl/cmd/grpcurl"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/momotaro98/strictimportsort/cmd/strictimportsort"
	_ "golang.org/x/tools/cmd/goimports"
)
