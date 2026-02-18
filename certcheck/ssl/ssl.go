// Package ssl provides functions to check SSL/TLS certificates.
package ssl

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"time"
)

// CertInfo contains information about an SSL certificate
type CertInfo struct {
	Host      string
	ExpiresAt time.Time
	DaysUntil int
}

// CheckCertificate connects to the given URL and returns certificate information
func CheckCertificate(rawURL string, timeout time.Duration) (*CertInfo, error) {
	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	host := parsedURL.Hostname()
	port := parsedURL.Port()
	if port == "" {
		port = "443"
	}

	// Connect with TLS
	dialer := &net.Dialer{
		Timeout: timeout,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", net.JoinHostPort(host, port), &tls.Config{
		InsecureSkipVerify: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	// Get the certificate chain
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, fmt.Errorf("no certificates found")
	}

	// Use the leaf certificate (first in chain)
	cert := certs[0]
	expiresAt := cert.NotAfter
	daysUntil := int(time.Until(expiresAt).Hours() / 24)

	return &CertInfo{
		Host:      host,
		ExpiresAt: expiresAt,
		DaysUntil: daysUntil,
	}, nil
}
