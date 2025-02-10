all: clean go-grpc tests

##### Compile proto files for go #####
go-grpc:
	buf generate
	mv -f api/temporal/api/cloud/* api && rm -rf api/temporal

##### Tests #####

tests:
	go test -v ./cloudclient

##### api-cloud Submodule #####
load-submodule:
	printf "Load the api-cloud submodule..."
	git submodule update --init --recursive

update-submodule:
	# "Update api-cloud submodule..."
	git submodule update --recursive --remote --merge

##### Plugins & tools #####
install: buf-install grpc-install

buf-install:
	# "Install/update buf..."
	go install -modfile internal/build/go.mod github.com/bufbuild/buf/cmd/buf

grpc-install:
	# "Install/update grpc and plugins..."
	go install -modfile internal/build/go.mod google.golang.org/protobuf/cmd/protoc-gen-go
	go install -modfile internal/build/go.mod google.golang.org/grpc/cmd/protoc-gen-go-grpc

##### Clean #####
clean:
	rm -rf api/*
