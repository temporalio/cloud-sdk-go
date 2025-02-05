package cloudclient_test

import (
	"context"
	"os"
	"testing"

	cloudservicev1 "go.temporal.io/cloud-sdk/api/cloudservice/v1"
	"go.temporal.io/cloud-sdk/cloudclient"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// TemporalCloudAPIKeyEnv is the environment variable that contains the Temporal Cloud API key.
	TemporalCloudAPIKeyEnv = "TEMPORAL_API_KEY"
)

func TestClient(t *testing.T) {

	// Create a new client
	client, err := cloudclient.New()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	_, err = client.CloudService().GetNamespaces(ctx, &cloudservicev1.GetNamespacesRequest{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("expected error code %v, got %v", codes.Unauthenticated, status.Code(err))
	}

	apikey := os.Getenv("TEMPORAL_API_KEY")
	if apikey == "" {
		return
	}

	client, err = cloudclient.New(cloudclient.WithAPIKey(func() (string, error) {
		return apikey, nil
	}))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.CloudService().GetNamespaces(ctx, &cloudservicev1.GetNamespacesRequest{})
	if err != nil {
		t.Fatalf("failed to get namespaces: %v", err)
	}
}
