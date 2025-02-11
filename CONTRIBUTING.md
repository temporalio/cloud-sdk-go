# Developing Temporal Cloud Go SDK

This doc is intended for contributors to Go SDK (hopefully that's you!)

**Note:** All contributors also need to fill out the [Temporal Contributor License Agreement](https://gist.github.com/samarabbas/7dcd41eb1d847e12263cc961ccfdb197) before we can merge in any of your changes.

## Development Environment

* [Go Lang](https://golang.org/): Follow instructions to install Go from [here](https://go.dev/doc/install).

## Checking out the code

Temporal Cloud GO SDK uses go modules, there is no dependency on `$GOPATH` variable. Clone the repo into the preferred location:

```bash
git clone https://github.com/temporalio/cloud-sdk-go.git
```

## Running protocol buffer compiler

Temporal Cloud Go SDK uses protocol buffers to define the API. To generate the Go code from the proto files, run the following commands:

```bash
# Load the proto submodule
git submodule update --init --recursive
# Install the dependencies tools at the correct version
go generate ./internal/build/tools.go
# Run buf to generate
buf generate

```

## Getting the latest protos
To get the latest protos, run the following commands:

```bash
git submodule update --recursive --remote --merge
```


## Testing

Run the tests:

```bash
export TEMPORAL_API_KEY=<your_api_key>
go test ./cloudclient/...
```

You can generate an api key by following the instructions [here](https://docs.temporal.io/cloud/api-keys#generate-an-api-key).


## Updating go mod files

Sometimes all go.mod files need to be tidied. For an easy way to do this on linux or (probably) mac,
run:

```bash
find . -name go.mod -execdir go mod tidy \;
```

