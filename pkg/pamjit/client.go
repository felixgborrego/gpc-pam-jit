package pamjit

import (
	privilegedaccessmanager "cloud.google.com/go/privilegedaccessmanager/apiv1"
	"cloud.google.com/go/privilegedaccessmanager/apiv1/privilegedaccessmanagerpb"
	"context"
	"fmt"
)

type Client struct {
	gcpClient *privilegedaccessmanager.Client
	projectID string
	location  string
}

func NewPamJitClient(ctx context.Context, projectID, location string) (*Client, error) {
	pamClient, err := privilegedaccessmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create PAM client: %w", err)
	}

	client := &Client{
		gcpClient: pamClient,
		projectID: projectID,
		location:  location,
	}

	if err := client.CheckOnboardingStatus(ctx); err != nil {
		return nil, fmt.Errorf("onboarding status check failed: %w", err)
	}

	return client, nil
}

// CheckOnboardingStatus checks if the user is onboarded to PAM JIT and returns an error if not.
func (c *Client) CheckOnboardingStatus(ctx context.Context) error {
	req := &privilegedaccessmanagerpb.CheckOnboardingStatusRequest{
		Parent: c.parent(),
	}

	resp, err := c.gcpClient.CheckOnboardingStatus(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to check onboarding status: %w", err)
	}

	if len(resp.Findings) > 0 {
		var findings []string
		for _, f := range resp.Findings {
			findings = append(findings, f.String())
		}
		return fmt.Errorf("user is not onboarded: %s, findings: %s", resp.String(), findings)
	}

	return nil
}

// parent returns the resource name of the project and location.
func (c *Client) parent() string {
	return fmt.Sprintf("projects/%s/locations/%s", c.projectID, c.location)
}