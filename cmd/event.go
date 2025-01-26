package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/EnSync-engine/CLI/app/api"
	"github.com/EnSync-engine/CLI/app/domain"
)

func newEventCmd(client *api.Client) *cobra.Command {
	var accessKey string

	cmd := &cobra.Command{
		Use:   "event",
		Short: "Manage events",
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

	// Add the access key flag to all event subcommands
	cmd.PersistentFlags().StringVar(&accessKey, "access-key", "", "Access key for API authentication")
	cmd.MarkPersistentFlagRequired("access-key")

	// Add subcommands
	cmd.AddCommand(
		newEventListCmd(client),
		newEventCreateCmd(client),
		newEventUpdateCmd(client),
		newEventGetByNameCmd(client),
	)

	return cmd
}

func newEventListCmd(client *api.Client) *cobra.Command {
	var (
		pageIndex int
		limit     int
		order     string
		orderBy   string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List events",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := &api.ListParams{
				PageIndex: pageIndex,
				Limit:     limit,
				Order:     order,
				OrderBy:   orderBy,
			}

			events, err := client.ListEvents(context.Background(), params)
			if err != nil {
				return fmt.Errorf("failed to list events: %w", err)
			}

			return printJSON(cmd.OutOrStdout(), events)
		},
	}

	cmd.Flags().IntVar(&pageIndex, "page", 0, "Page index")
	cmd.Flags().IntVar(&limit, "limit", 10, "Number of items per page")
	cmd.Flags().StringVar(&order, "order", "DESC", "Sort order (ASC/DESC)")
	cmd.Flags().StringVar(&orderBy, "order-by", "createdAt", "Field to order by")

	return cmd
}

func newEventCreateCmd(client *api.Client) *cobra.Command {
	var (
		name    string
		payload string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new event definition",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("name is required")
			}

			var payloadMap map[string]string
			if err := json.Unmarshal([]byte(payload), &payloadMap); err != nil {
				return fmt.Errorf("invalid payload JSON: %w", err)
			}

			event := &domain.Event{
				Name:    name,
				Payload: payloadMap,
			}

			ctx := context.Background()
			err := client.CreateEvent(ctx, event)
			if err != nil {
				return fmt.Errorf("failed to create event: %w", err)
			}

			fmt.Printf("Event '%s' created successfully\n", name)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Event name")
	cmd.Flags().StringVar(&payload, "payload", "{}", "Event payload in JSON format")
	cmd.MarkFlagRequired("name")

	return cmd
}

func newEventUpdateCmd(client *api.Client) *cobra.Command {
	var (
		id      int64
		name    string
		payload string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing event definition",
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == 0 {
				return fmt.Errorf("id is required")
			}

			var payloadMap map[string]string
			if err := json.Unmarshal([]byte(payload), &payloadMap); err != nil {
				return fmt.Errorf("invalid payload JSON: %w", err)
			}

			event := &domain.Event{
				ID:      id,
				Name:    name,
				Payload: payloadMap,
			}

			ctx := context.Background()
			err := client.UpdateEvent(ctx, event)
			if err != nil {
				return fmt.Errorf("failed to update event: %w", err)
			}

			fmt.Printf("Event '%d' updated successfully\n", id)
			return nil
		},
	}

	cmd.Flags().Int64Var(&id, "id", 0, "Event ID")
	cmd.Flags().StringVar(&name, "name", "", "New event name")
	cmd.Flags().StringVar(&payload, "payload", "{}", "Event payload in JSON format")
	cmd.MarkFlagRequired("id")

	return cmd
}

func newEventGetByNameCmd(client *api.Client) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get event by name",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("name is required")
			}

			ctx := context.Background()
			event, err := client.GetEventByName(ctx, name)
			if err != nil {
				return fmt.Errorf("failed to get event: %w", err)
			}

			return printJSON(cmd.OutOrStdout(), event)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Event name")
	cmd.MarkFlagRequired("name")

	return cmd
}
