package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alertbingo/cli/api"
	"github.com/alertbingo/cli/hoststats"
	"github.com/urfave/cli/v3"
)

var version = "dev"

func main() {
	cmd := &cli.Command{
		Name:    "alertbingo",
		Usage:   "CLI tool for sending checks to Alert Bingo",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:  "hoststats",
				Usage: "Send host statistics checks (memory, uptime, CPU) to Alert Bingo",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "dashboard",
						Aliases:  []string{"d"},
						Usage:    "Dashboard name",
						Sources:  cli.EnvVars("ALERTBINGO_DASHBOARD"),
						Required: true,
					},
					&cli.StringFlag{
						Name:     "site",
						Aliases:  []string{"s"},
						Usage:    "Site identifier (e.g., myapp-prod)",
						Sources:  cli.EnvVars("ALERTBINGO_SITE"),
						Required: true,
					},
					&cli.StringFlag{
						Name:     "service",
						Usage:    "Service name (e.g., host)",
						Sources:  cli.EnvVars("ALERTBINGO_SERVICE"),
						Required: true,
					},
					&cli.StringFlag{
						Name:    "message",
						Aliases: []string{"m"},
						Usage:   "Optional long-form status message",
						Sources: cli.EnvVars("ALERTBINGO_MESSAGE"),
					},
					&cli.StringFlag{
						Name:    "inactive-expire",
						Usage:   "Optional duration string for inactive expiry (e.g., 48h or 30m)",
						Sources: cli.EnvVars("ALERTBINGO_INACTIVE_EXPIRE"),
					},
					&cli.StringFlag{
						Name:    "inactive-escalate",
						Usage:   "Optional duration string for inactive escalation (e.g., 1h or 30m)",
						Sources: cli.EnvVars("ALERTBINGO_INACTIVE_ESCALATE"),
					},
					&cli.StringFlag{
						Name:    "highlighted",
						Usage:   "Optional highlighted status",
						Sources: cli.EnvVars("ALERTBINGO_HIGHLIGHTED"),
					},
					&cli.StringFlag{
						Name:     "token",
						Aliases:  []string{"t"},
						Usage:    "API Bearer token",
						Sources:  cli.EnvVars("ALERTBINGO_TOKEN"),
						Required: true,
					},
					&cli.StringFlag{
						Name:    "api-url",
						Usage:   "API URL",
						Sources: cli.EnvVars("ALERTBINGO_API_URL"),
						Value:   "https://app.alert.bingo/api/v1/checks",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					cfg := hoststats.Config{
						Dashboard:        cmd.String("dashboard"),
						Site:             cmd.String("site"),
						Service:          cmd.String("service"),
						Message:          cmd.String("message"),
						InactiveExpire:   cmd.String("inactive-expire"),
						InactiveEscalate: cmd.String("inactive-escalate"),
						Highlighted:      cmd.String("highlighted"),
					}

					checks, err := hoststats.Collect(ctx, cfg)
					if err != nil {
						return fmt.Errorf("failed to collect host stats: %w", err)
					}

					client := api.NewClient(cmd.String("api-url"), cmd.String("token"))
					responses, err := client.SendChecks(ctx, checks)
					if err != nil {
						return err
					}

					fmt.Println("Host stats checks sent successfully")

					// Report any non-OK statuses and errors
					if summary := formatResponses(responses); summary != "" {
						fmt.Println(summary)
					}

					return nil
				},
			},
			{
				Name:  "check",
				Usage: "Send a check to Alert Bingo",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "dashboard",
						Aliases:  []string{"d"},
						Usage:    "Dashboard name",
						Sources:  cli.EnvVars("ALERTBINGO_DASHBOARD"),
						Required: true,
					},
					&cli.StringFlag{
						Name:     "site",
						Aliases:  []string{"s"},
						Usage:    "Site identifier (e.g., myapp-prod)",
						Sources:  cli.EnvVars("ALERTBINGO_SITE"),
						Required: true,
					},
					&cli.StringFlag{
						Name:     "service",
						Usage:    "Service name (e.g., postgres)",
						Sources:  cli.EnvVars("ALERTBINGO_SERVICE"),
						Required: true,
					},
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Usage:    "Check name (e.g., postgres-rds-space-free)",
						Sources:  cli.EnvVars("ALERTBINGO_NAME"),
						Required: true,
					},
					&cli.StringFlag{
						Name:    "alert-level",
						Aliases: []string{"l"},
						Usage:   "Alert level: ok, warn, or alert",
						Sources: cli.EnvVars("ALERTBINGO_ALERT_LEVEL"),
						Value:   "ok",
					},
					&cli.StringFlag{
						Name:    "message",
						Aliases: []string{"m"},
						Usage:   "Optional long-form status message",
						Sources: cli.EnvVars("ALERTBINGO_MESSAGE"),
					},
					&cli.StringFlag{
						Name:    "value",
						Aliases: []string{"v"},
						Usage:   "Short-form status value",
						Sources: cli.EnvVars("ALERTBINGO_VALUE"),
					},
					&cli.StringFlag{
						Name:    "inactive-expire",
						Usage:   "Optional duration string for inactive expiry (e.g., 48h or 30m)",
						Sources: cli.EnvVars("ALERTBINGO_INACTIVE_EXPIRE"),
					},
					&cli.StringFlag{
						Name:    "inactive-escalate",
						Usage:   "Optional duration string for inactive escalation (e.g., 1h or 30m)",
						Sources: cli.EnvVars("ALERTBINGO_INACTIVE_ESCALATE"),
					},
					&cli.StringFlag{
						Name:    "highlighted",
						Usage:   "Optional highlighted status",
						Sources: cli.EnvVars("ALERTBINGO_HIGHLIGHTED"),
					},
					&cli.StringFlag{
						Name:     "token",
						Aliases:  []string{"t"},
						Usage:    "API Bearer token",
						Sources:  cli.EnvVars("ALERTBINGO_TOKEN"),
						Required: true,
					},
					&cli.StringFlag{
						Name:    "api-url",
						Usage:   "API URL",
						Sources: cli.EnvVars("ALERTBINGO_API_URL"),
						Value:   "https://app.alert.bingo/api/v1/checks",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					alertLevel, err := api.ParseAlertLevel(cmd.String("alert-level"))
					if err != nil {
						return err
					}

					payload := api.CheckPayload{
						Dashboard:        cmd.String("dashboard"),
						Site:             cmd.String("site"),
						Service:          cmd.String("service"),
						Name:             cmd.String("name"),
						AlertLevel:       alertLevel,
						Message:          cmd.String("message"),
						Value:            cmd.String("value"),
						InactiveExpire:   cmd.String("inactive-expire"),
						InactiveEscalate: cmd.String("inactive-escalate"),
						Highlighted:      cmd.String("highlighted"),
					}

					client := api.NewClient(cmd.String("api-url"), cmd.String("token"))
					responses, err := client.SendChecks(ctx, []api.CheckPayload{payload})
					if err != nil {
						return err
					}

					fmt.Println("Check sent successfully")

					// Report any non-OK statuses and errors
					if summary := formatResponses(responses); summary != "" {
						fmt.Println(summary)
					}

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// formatResponses returns a formatted summary of non-OK statuses and errors
func formatResponses(responses []api.Response) string {
	// Collect unique non-OK statuses
	statusSet := make(map[string]struct{})
	for _, resp := range responses {
		if resp.Status != "OK" && resp.Status != "ok" {
			statusSet[resp.Status] = struct{}{}
		}
	}

	// Collect unique errors
	errorSet := make(map[string]struct{})
	for _, resp := range responses {
		for _, e := range resp.Errors {
			errorSet[e] = struct{}{}
		}
	}

	var parts []string

	// Format unique non-OK statuses
	if len(statusSet) > 0 {
		var statuses []string
		for s := range statusSet {
			statuses = append(statuses, s)
		}
		parts = append(parts, fmt.Sprintf("Statuses: %s", strings.Join(statuses, ", ")))
	}

	// Format unique errors
	if len(errorSet) > 0 {
		var errors []string
		for e := range errorSet {
			errors = append(errors, e)
		}
		parts = append(parts, fmt.Sprintf("Errors: %s", strings.Join(errors, ", ")))
	}

	return strings.Join(parts, "; ")
}
