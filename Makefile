all: clean go-grpc tests

##### Compile proto files for go #####
go-grpc:
	buf generate
	mv -f api/temporal/api/cloud/* api && rm -rf api/temporal
	find api -type f -name '*.go' -exec sed -i '' 's/api\/temporal\/api\/cloud/api/g' {} +
	find api -type f -name '*.go' -exec sed -i '' 's/go.temporal.io\/cloud-sdk\/api\/temporal\/api/go.temporal.io\/api/g' {} +
	find api -type f -name '*.go' -exec sed -i '' 's/go.temporal.io\/cloud-sdk\/api\/google\/api/google.golang.org\/api/g' {} +

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
	go install github.com/bufbuild/buf/cmd/buf@v1.49.0

grpc-install:
	# "Install/update grpc and plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.4
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

##### Clean #####
clean:
	rm -rf api/*
