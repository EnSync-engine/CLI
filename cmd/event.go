package cmd

import (
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
		Long:  "Commands for listing, creating, updating, and retrieving events.",
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
		newEventListCmd(client),
		newEventGetCmd(client),
		newEventCreateCmd(client),
		newEventUpdateCmd(client),
	)

	return cmd
}

func newEventListCmd(client *api.Client) *cobra.Command {
	var (
		page    int
		limit   int
		order   string
		orderBy string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List events",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := &api.ListParams{
				PageIndex: page,
				Limit:     limit,
				Order:     order,
				OrderBy:   orderBy,
			}

			events, err := client.ListEvents(cmd.Context(), params)
			if err != nil {
				return err
			}

			return printJSON(cmd.OutOrStdout(), events)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "page index (0-based)")
	cmd.Flags().IntVar(&limit, "limit", 20, "items per page")
	cmd.Flags().StringVar(&order, "order", "DESC", "sort order (ASC or DESC)")
	cmd.Flags().StringVar(&orderBy, "order-by", "createdAt", "field to order by")

	return cmd
}

func newEventGetCmd(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [name]",
		Short: "Get event by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			event, err := client.GetEventByName(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return printJSON(cmd.OutOrStdout(), event)
		},
	}

	return cmd
}

func newEventCreateCmd(client *api.Client) *cobra.Command {
	var (
		name        string
		payloadJSON string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new event",
		RunE: func(cmd *cobra.Command, args []string) error {
			payload, err := parsePayloadJSON(payloadJSON)
			if err != nil {
				return err
			}

			event := &domain.Event{
				Name:    name,
				Payload: payload,
			}

			if err := client.CreateEvent(cmd.Context(), event); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Event %q created successfully\n", name)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "event name (required)")
	cmd.Flags().StringVar(&payloadJSON, "payload", "{}", "event payload as JSON")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func newEventUpdateCmd(client *api.Client) *cobra.Command {
	var (
		id          string
		name        string
		payloadJSON string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing event",
		RunE: func(cmd *cobra.Command, args []string) error {
			payload, err := parsePayloadJSON(payloadJSON)
			if err != nil {
				return err
			}

			event := &domain.Event{
				ID:      id,
				Name:    name,
				Payload: payload,
			}

			if err := client.UpdateEvent(cmd.Context(), event); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Event %s updated successfully\n", id)
			return nil
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "event ID (required)")
	cmd.Flags().StringVar(&name, "name", "", "new event name")
	cmd.Flags().StringVar(&payloadJSON, "payload", "{}", "event payload as JSON")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}

func parsePayloadJSON(s string) (map[string]any, error) {
	var payload map[string]any
	if err := json.Unmarshal([]byte(s), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload JSON: %w", err)
	}
	return payload, nil
}
