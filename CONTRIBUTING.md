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
# Clean the generated files
rm -rf api/*
# Run buf to generate
buf generate
# Move the generated files to the correct location
mv -f api/temporal/api/cloud/* api && rm -rf api/temporal
# Update the default API version in cloudclient/options.go
sed -i '' 's/defaultAPIVersion = ".*"/defaultAPIVersion = "'$(cat proto/cloud-api/VERSION)'"/' cloudclient/options.go

```

## Getting the latest protos
To get the latest protos, run the following commands:

```bash
git submodule update --recursive --remote --merge
```

### Automated Proto Updates
For convenience, there's a GitHub Action workflow that can automatically update the protos and create a PR. This workflow:

1. Updates the proto submodule to get the latest changes or a specific release
2. Regenerates all Go code from the proto files
3. Updates the default API version in `cloudclient/options.go`
4. Increments the patch version of the SDK version in `cloudclient/options.go`
5. Creates a new branch and commits the changes
6. Opens a pull request with the updates

To use this workflow:
1. Go to the "Actions" tab in the GitHub repository
2. Select "Update Protos and Create PR" from the workflow list
3. Click "Run workflow"
4. Provide the following optional inputs:
   - **Release version**: Specify a cloud-api release version (e.g., `v1.2.3`) to checkout that specific release. Leave empty to get the latest.
   - **Branch name**: Custom branch name for the PR (default: auto-generated)
   - **PR title**: Custom PR title (default: auto-generated)
5. Click "Run workflow" to execute

The workflow will automatically handle all the steps from the manual process above and create a PR for review.

## Testing

Run the tests:

```bash
export TEST_TEMPORAL_CLOUD_SDK_API_KEY=<your_api_key>
go test ./cloudclient/...
```

You can generate an api key by following the instructions [here](https://docs.temporal.io/cloud/api-keys#generate-an-api-key).


## Updating go mod files

Sometimes all go.mod files need to be tidied. For an easy way to do this on linux or (probably) mac,
run:

```bash
find . -name go.mod -execdir go mod tidy \;
```

