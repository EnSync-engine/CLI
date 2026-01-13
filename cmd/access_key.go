package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/EnSync-engine/CLI/app/api"
	"github.com/EnSync-engine/CLI/app/domain"
)

func newAccessKeyCmd(client *api.Client) *cobra.Command {
	var accessKey string

	cmd := &cobra.Command{
		Use:   "access-key",
		Short: "Manage access keys",
		Long:  "Commands for listing, creating, and managing access key permissions.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if accessKey == "" {
				return fmt.Errorf("--access-key is required")
			}
			client.SetAccessKey(accessKey)
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&accessKey, "access-key", "", "access key for API authentication (required)")
	_ = cmd.MarkPersistentFlagRequired("access-key")

	cmd.AddCommand(
		newAccessKeyListCmd(client),
		newAccessKeyGetCmd(client),
		newAccessKeyCreateCmd(client),
		newAccessKeyDeleteCmd(client),
		newAccessKeyPermissionsCmd(client),
		newAccessKeyRotateCmd(client),
	)

	return cmd
}

func newAccessKeyListCmd(client *api.Client) *cobra.Command {
	var (
		page      int
		limit     int
		order     string
		orderBy   string
		filterKey string
		name      string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List access keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := &api.ListParams{
				PageIndex: page,
				Limit:     limit,
				Order:     order,
				OrderBy:   orderBy,
			}
			if filterKey != "" {
				params.Filter = map[string]string{"accessKey": filterKey}
			}
			if name != "" {
				if params.Filter == nil {
					params.Filter = make(map[string]string)
				}
				params.Filter["name"] = name
			}

			keys, err := client.ListAccessKeys(cmd.Context(), params)
			if err != nil {
				return err
			}

			return printJSON(cmd.OutOrStdout(), keys)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "page index (0-based)")
	cmd.Flags().IntVar(&limit, "limit", 20, "items per page")
	cmd.Flags().StringVar(&order, "order", "DESC", "sort order (ASC or DESC)")
	cmd.Flags().StringVar(&orderBy, "order-by", "key", "field to order by")
	cmd.Flags().StringVar(&filterKey, "filter-key", "", "filter by access key")
	cmd.Flags().StringVar(&name, "name", "", "filter by name")

	return cmd
}

func newAccessKeyGetCmd(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Get access key by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := client.GetAccessKeyByID(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return printJSON(cmd.OutOrStdout(), key)
		},
	}

	return cmd
}

func newAccessKeyCreateCmd(client *api.Client) *cobra.Command {
	var (
		keyType         string
		name            string
		permissionsJSON string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new access key",
		RunE: func(cmd *cobra.Command, args []string) error {
			var permissions *domain.Permissions
			if permissionsJSON != "" {
				if err := json.Unmarshal([]byte(permissionsJSON), &permissions); err != nil {
					return fmt.Errorf("invalid permissions JSON: %w", err)
				}
			}

			req := &domain.CreateAccessKeyRequest{
				Type:        keyType,
				Name:        name,
				Permissions: permissions,
			}

			key, err := client.CreateAccessKey(cmd.Context(), req)
			if err != nil {
				return err
			}

			return printJSON(cmd.OutOrStdout(), key)
		},
	}

	cmd.Flags().StringVar(&keyType, "type", "SERVICE", "access key type (SERVICE or ACCOUNT)")
	cmd.Flags().StringVar(&name, "name", "", "access key name (required)")
	cmd.Flags().StringVar(&permissionsJSON, "permissions", "", `permissions JSON`)
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func newAccessKeyDeleteCmd(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete an access key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := client.DeleteAccessKey(cmd.Context(), args[0]); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Access key %q deleted successfully\n", args[0])
			return nil
		},
	}

	return cmd
}

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

func newAccessKeyGetPermissionsCmd(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [key]",
		Short: "Get permissions for an access key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			permissions, err := client.GetAccessKeyPermissions(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return printJSON(cmd.OutOrStdout(), permissions)
		},
	}

	return cmd
}

func newAccessKeySetPermissionsCmd(client *api.Client) *cobra.Command {
	var permissionsJSON string

	cmd := &cobra.Command{
		Use:   "set [key]",
		Short: "Set permissions for an access key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var permissions domain.Permissions
			if err := json.Unmarshal([]byte(permissionsJSON), &permissions); err != nil {
				return fmt.Errorf("invalid permissions JSON: %w", err)
			}

			if err := client.SetAccessKeyPermissions(cmd.Context(), args[0], &permissions); err != nil {
				return err
			}

			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Permissions updated successfully")
			return nil
		},
	}

	cmd.Flags().StringVar(&permissionsJSON, "permissions", "", `permissions JSON (required)`)
	_ = cmd.MarkFlagRequired("permissions")

	return cmd
}

func newAccessKeyRotateCmd(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate [key]",
		Short: "Rotate service key pair for an access key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			keyPair, err := client.UpdateServiceKeyPair(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return printJSON(cmd.OutOrStdout(), keyPair)
		},
	}

	return cmd
}
