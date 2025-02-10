package cloudclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc/credentials"
)

type (
	apikeyCreds struct {
		reader                 APIKeyReader
		allowInsecureTransport bool
	}
)

func (c apikeyCreds) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	ri, ok := credentials.RequestInfoFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to retrieve request info from context")
	}

	if !c.allowInsecureTransport {
		// ensure that the API key, AKA bearer token, is sent over a secure connection - meaning TLS.
		if err := credentials.CheckSecurityLevel(ri.AuthInfo, credentials.PrivacyAndIntegrity); err != nil {
			return nil, fmt.Errorf("the connection's transport security level is too low for API keys: %v", err)
		}
	}

	apiKey, err := c.reader.GetAPIKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %v", err)
	}

	return map[string]string{
		authorizationHeader: fmt.Sprintf("%s %s", authorizationBearer, apiKey),
	}, nil
}

func (c apikeyCreds) RequireTransportSecurity() bool {
	return !c.allowInsecureTransport
}
