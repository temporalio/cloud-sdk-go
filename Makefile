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
install:
	go generate ./internal/build/tools.go

##### Clean #####
clean:
	rm -rf api/*
