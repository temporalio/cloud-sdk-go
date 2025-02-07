##### Compile proto files for go #####
proto: clean go-grpc

go-grpc:
	buf generate
	mv -f api/temporal/api/cloud/* api && rm -rf api/temporal

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
	go install github.com/bufbuild/buf/cmd/buf@v1.49.0

grpc-install:
	# "Install/update grpc and plugins..."
	go install -modfile=go.mod google.golang.org/protobuf/cmd/protoc-gen-go
	go install -modfile=go.mod google.golang.org/grpc/cmd/protoc-gen-go-grpc

##### Clean #####
clean:
	rm -rf api/*
