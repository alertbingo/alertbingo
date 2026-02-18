// Package certcheck provides functions to check SSL/TLS certificates.
package certcheck

import (
	"context"
	"fmt"
	"time"

	"github.com/alertbingo/alertbingo/api"
	"github.com/alertbingo/alertbingo/certcheck/ssl"
)

// Config holds the common configuration for certificate checks
type Config struct {
	Dashboard        string
	Site             string
	Name             string
	Message          string
	InactiveExpire   string
	InactiveEscalate string
	Highlighted      string
	Timeout          time.Duration
}

// Collect gathers certificate information for the given URLs and returns check payloads
func Collect(ctx context.Context, cfg Config, urls []string) []api.CheckPayload {
	var checks []api.CheckPayload

	for _, u := range urls {
		check := checkCertificate(ctx, cfg, u)
		checks = append(checks, check)
	}

	return checks
}

func checkCertificate(ctx context.Context, cfg Config, rawURL string) api.CheckPayload {
	certInfo, err := ssl.CheckCertificate(rawURL, cfg.Timeout)
	if err != nil {
		return api.CheckPayload{
			Dashboard:        cfg.Dashboard,
			Site:             cfg.Site,
			Service:          rawURL,
			Name:             cfg.Name,
			AlertLevel:       2, // alert
			Value:            "Error",
			Message:          err.Error(),
			InactiveExpire:   cfg.InactiveExpire,
			InactiveEscalate: cfg.InactiveEscalate,
			Highlighted:      cfg.Highlighted,
		}
	}

	// Determine alert level based on days until expiry
	alertLevel := 0
	message := cfg.Message
	var value string

	if certInfo.DaysUntil < 0 {
		// Expired
		alertLevel = 2
		value = "Expired"
		message = appendAlertReason(cfg.Message, fmt.Sprintf("Certificate expired %d days ago", -certInfo.DaysUntil))
	} else if certInfo.DaysUntil <= 14 {
		// 2 weeks or less remaining
		alertLevel = 1
		value = fmt.Sprintf("%dd", certInfo.DaysUntil)
		message = appendAlertReason(cfg.Message, fmt.Sprintf("Certificate expires in %d days", certInfo.DaysUntil))
	} else {
		// OK
		alertLevel = 0
		value = fmt.Sprintf("%dd", certInfo.DaysUntil)
	}

	return api.CheckPayload{
		Dashboard:        cfg.Dashboard,
		Site:             cfg.Site,
		Service:          certInfo.Host,
		Name:             cfg.Name,
		AlertLevel:       alertLevel,
		Value:            value,
		Message:          message,
		InactiveExpire:   cfg.InactiveExpire,
		InactiveEscalate: cfg.InactiveEscalate,
		Highlighted:      cfg.Highlighted,
	}
}

// appendAlertReason appends an alert reason to an existing message
func appendAlertReason(message, reason string) string {
	if message == "" {
		return reason
	}
	return message + " - " + reason
}
