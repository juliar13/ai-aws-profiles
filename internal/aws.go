package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
)

type AWSClient struct {
	orgsClient *organizations.Client
}

type AccountInfo struct {
	ID   string
	Name string
}

func NewAWSClient(ctx context.Context) (*AWSClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &AWSClient{
		orgsClient: organizations.NewFromConfig(cfg),
	}, nil
}

func (c *AWSClient) ListAccounts(ctx context.Context) ([]AccountInfo, error) {
	var accounts []AccountInfo
	var nextToken *string

	for {
		input := &organizations.ListAccountsInput{
			NextToken: nextToken,
		}

		result, err := c.orgsClient.ListAccounts(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list accounts: %w", err)
		}

		for _, account := range result.Accounts {
			if account.Status == types.AccountStatusActive {
				accountInfo := AccountInfo{
					ID:   *account.Id,
					Name: sanitizeAccountName(*account.Name),
				}
				accounts = append(accounts, accountInfo)
			}
		}

		nextToken = result.NextToken
		if nextToken == nil {
			break
		}
	}

	return accounts, nil
}

func sanitizeAccountName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")
	
	var result strings.Builder
	prevWasDash := false
	
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
			prevWasDash = false
		} else if r == '-' && !prevWasDash {
			result.WriteRune(r)
			prevWasDash = true
		}
	}
	
	return strings.Trim(result.String(), "-")
}

func GenerateFromAWSAccounts(ctx context.Context) (*Generator, error) {
	client, err := NewAWSClient(ctx)
	if err != nil {
		return nil, err
	}

	accounts, err := client.ListAccounts(ctx)
	if err != nil {
		return nil, err
	}

	generator := NewGenerator()
	err = generator.ColorManager.LoadColorSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to load color settings: %w", err)
	}
	
	for _, account := range accounts {
		generator.GenerateProfiles(account.ID, account.Name)
	}

	return generator, nil
}