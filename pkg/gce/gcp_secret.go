// Copyright 2021 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package gce

import (
	"context"
	"fmt"

	"cloud.google.com/go/compute/metadata"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// GcpSecret returns the GCP Secret Manager blob as a []byte data.
func GcpSecret(name string) ([]byte, error) {
	return GcpSecretWithContext(context.Background(), name)
}

func GcpSecretWithContext(ctx context.Context, name string) ([]byte, error) {
	// name := "projects/my-project/secrets/my-secret/versions/5"
	// name := "projects/my-project/secrets/my-secret/versions/latest"

	// Create the client.
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, err
	}

	return result.Payload.Data, nil
}

// LatestGcpSecret returns the latest secret value.
func LatestGcpSecret(ctx context.Context, projectName, key string) ([]byte, error) {
	return GcpSecretWithContext(ctx,
		fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectName, key))
}

// ProjectName returns the name of the GCP project the code is running on.
func ProjectName(ctx context.Context) (string, error) {
	if !metadata.OnGCE() {
		return "", fmt.Errorf("not running on GKE/GCE")
	}
	projectID, err := metadata.ProjectIDWithContext(ctx)
	if err != nil {
		return "", err
	}
	return projectID, nil
}
