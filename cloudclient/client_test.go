package cloudclient_test

import (
	"context"
	"os"
	"testing"

	cloudservicev1 "go.temporal.io/cloud-sdk/api/cloudservice/v1"
	"go.temporal.io/cloud-sdk/cloudclient"
)

const (
	// The environment variable that contains the Temporal Cloud API key to use for testing.
	temporalCloudAPIKeyEnv = "TEST_TEMPORAL_CLOUD_SDK_API_KEY"
)

func TestClient(t *testing.T) {

	apikey := os.Getenv(temporalCloudAPIKeyEnv)
	if apikey == "" {
		t.Skipf("skipping test; environment variable %s is not set", temporalCloudAPIKeyEnv)
	}
	t.Run("New", func(t *testing.T) {
		client, err := cloudclient.New(cloudclient.Options{
			APIKey: apikey,
		})
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		defer client.Close()

		ctx := context.Background()
		_, err = client.CloudService().GetNamespaces(ctx, &cloudservicev1.GetNamespacesRequest{})
		if err != nil {
			t.Fatalf("failed to get namespaces: %v", err)
		}
	})
}
