package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/EnSync-engine/CLI/app/api"
	"github.com/EnSync-engine/CLI/app/domain"
	"github.com/spf13/cobra"
)

// newAccessKeyCmd creates and returns the `access-key` command with its subcommands.
func newAccessKeyCmd(client *api.Client) *cobra.Command {
	var accessKey string

	cmd := &cobra.Command{
		Use:   "access-key",
		Short: "Manage access keys",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate the access key
			if accessKey == "" {
				return fmt.Errorf("access key is required")
			}

			// Set the access key in the client for authentication
			client.SetAccessKey(accessKey)
			return nil
		},
	}

	// Add the access key flag to all access-key subcommands
	cmd.PersistentFlags().StringVar(&accessKey, "access-key", "", "Access key for API authentication")
	cmd.MarkPersistentFlagRequired("access-key")

	cmd.AddCommand(
		newAccessKeyListCmd(client),
		newAccessKeyCreateCmd(client),
		newAccessKeyPermissionsCmd(client),
	)

	return cmd
}

// newAccessKeyListCmd creates and returns the `access-key list` command.
func newAccessKeyListCmd(client *api.Client) *cobra.Command {
	var (
		pageIndex int
		limit     int
		order     string
		orderBy   string
		filterKey string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List access keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := &api.ListParams{
				PageIndex: pageIndex,
				Limit:     limit,
				Order:     order,
				OrderBy:   orderBy,
				Filter:    map[string]string{"accessKey": filterKey},
			}
			keys, err := client.ListAccessKeys(context.Background(), params)
			if err != nil {
				return fmt.Errorf("failed to list access keys: %w", err)
			}

			return printJSON(cmd.OutOrStdout(), keys)
		},
	}

	cmd.Flags().IntVar(&pageIndex, "page", 0, "Page index")
	cmd.Flags().IntVar(&limit, "limit", 10, "Number of items per page")
	cmd.Flags().StringVar(&order, "order", "DESC", "Sort order (ASC/DESC)")
	cmd.Flags().StringVar(&orderBy, "order-by", "createdAt", "Field to order by")
	cmd.Flags().StringVar(&filterKey, "key", "", "Filter by access key")

	return cmd
}

// newAccessKeyCreateCmd creates and returns the `access-key create` command.
func newAccessKeyCreateCmd(client *api.Client) *cobra.Command {
	var permissionsJSON string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new access key with permissions",
		RunE: func(cmd *cobra.Command, args []string) error {
			var permissions *domain.Permissions
			if permissionsJSON != "" {
				if err := json.Unmarshal([]byte(permissionsJSON), &permissions); err != nil {
					return fmt.Errorf("failed to parse permissions JSON: %w", err)
				}
			}

			createdKey, err := client.CreateAccessKey(context.Background(), permissions)
			if err != nil {
				return fmt.Errorf("failed to create access key: %w", err)
			}

			return printJSON(cmd.OutOrStdout(), createdKey)
		},
	}

	cmd.Flags().StringVar(&permissionsJSON, "permissions", "", "JSON string representing the permissions")
	return cmd
}

// newAccessKeyPermissionsCmd creates and returns the `access-key permissions` command.
func newAccessKeyPermissionsCmd(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "permissions",
		Short: "Manage access key permissions",
	}

	cmd.AddCommand(
		newAccessKeyGetPermissionsCmd(client),
		newAccessKeySetPermissionsCmd(client),
	)

	return cmd
}

// newAccessKeyGetPermissionsCmd creates and returns the `access-key permissions get` command.
func newAccessKeyGetPermissionsCmd(client *api.Client) *cobra.Command {
	var accessKey string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get access key permissions",
		RunE: func(cmd *cobra.Command, args []string) error {
			if accessKey == "" {
				return fmt.Errorf("access key is required")
			}

			permissions, err := client.GetAccessKeyPermissions(context.Background(), accessKey)
			if err != nil {
				return fmt.Errorf("failed to get permissions: %w", err)
			}

			return printJSON(cmd.OutOrStdout(), permissions)
		},
	}

	cmd.Flags().StringVar(&accessKey, "key", "", "Access key")
	cmd.MarkFlagRequired("key")

	return cmd
}

// newAccessKeySetPermissionsCmd creates and returns the `access-key permissions set` command.
func newAccessKeySetPermissionsCmd(client *api.Client) *cobra.Command {
	var (
		accessKey       string
		permissionsJSON string
	)

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set access key permissions",
		RunE: func(cmd *cobra.Command, args []string) error {
			if accessKey == "" {
				return fmt.Errorf("access key is required")
			}

			if permissionsJSON == "" {
				return fmt.Errorf("permissions JSON is required")
			}

			var permissions domain.Permissions
			if err := json.Unmarshal([]byte(permissionsJSON), &permissions); err != nil {
				return fmt.Errorf("failed to parse permissions JSON: %w", err)
			}

			err := client.SetAccessKeyPermissions(context.Background(), accessKey, &permissions)
			if err != nil {
				return fmt.Errorf("failed to set permissions: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Permissions updated successfully")
			return nil
		},
	}

	cmd.Flags().StringVar(&accessKey, "key", "", "Access key")
	cmd.Flags().StringVar(&permissionsJSON, "permissions", "", "JSON string representing permissions")
	cmd.MarkFlagRequired("key")
	cmd.MarkFlagRequired("permissions")

	return cmd
}
