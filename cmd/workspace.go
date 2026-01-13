package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/EnSync-engine/CLI/app/api"
)

func newWorkspaceCmd(client *api.Client) *cobra.Command {
	var accessKey string

	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Manage workspaces",
		Long:  "Commands for listing and creating workspaces.",
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
		newWorkspaceListCmd(client),
		newWorkspaceCreateCmd(client),
	)

	return cmd
}

func newWorkspaceListCmd(client *api.Client) *cobra.Command {
	var (
		page    int
		limit   int
		order   string
		orderBy string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List workspaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := &api.ListParams{
				PageIndex: page,
				Limit:     limit,
				Order:     order,
				OrderBy:   orderBy,
			}

			workspaces, err := client.ListWorkspaces(cmd.Context(), params)
			if err != nil {
				return err
			}

			return printJSON(cmd.OutOrStdout(), workspaces)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "page index (0-based)")
	cmd.Flags().IntVar(&limit, "limit", 20, "items per page")
	cmd.Flags().StringVar(&order, "order", "DESC", "sort order (ASC or DESC)")
	cmd.Flags().StringVar(&orderBy, "order-by", "name", "field to order by")

	return cmd
}

func newWorkspaceCreateCmd(client *api.Client) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := client.CreateWorkspace(cmd.Context(), name); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Workspace %q created successfully\n", name)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "workspace name (required)")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}
