package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"syscall"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/gomematic/gomematic-go/gomematic"
	"gopkg.in/urfave/cli.v2"

	transport "github.com/go-openapi/runtime/client"
)

// HandleFunc is the real handle implementation.
type HandleFunc func(c *cli.Context, client *Client) error

// Client simply wraps the openapi client including authentication.
type Client struct {
	*gomematic.GomematicOpen
	AuthInfo runtime.ClientAuthInfoWriter
}

// Handle wraps the command function handler.
func Handle(c *cli.Context, fn HandleFunc) error {
	if c.String("server") == "" {
		fmt.Fprintf(os.Stderr, "error: you must provide the server address.\n")
		os.Exit(1)
	}

	server, err := url.Parse(c.String("server"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: invalid server address, bad format?.\n")
		os.Exit(1)
	}

	client := &Client{
		GomematicOpen: gomematic.NewHTTPClientWithConfig(
			strfmt.Default,
			&gomematic.TransportConfig{
				Host: server.Host,
				BasePath: path.Join(
					gomematic.DefaultBasePath,
					server.Path,
				),
				Schemes: []string{
					server.Scheme,
				},
			},
		),
	}

	if c.String("token") != "" {
		client.AuthInfo = transport.APIKeyAuth(
			"X-API-Key",
			"header",
			c.String("token"),
		)
	} else {
		client.AuthInfo = transport.PassThroughAuth
	}

	if err := fn(c, client); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(2)
	}

	return nil
}

// PrettyError catches regular networking errors and prints it.
func PrettyError(err error) error {
	if val, ok := err.(net.Error); ok && val.Timeout() {
		return fmt.Errorf("connection to server timed out")
	}

	switch val := err.(type) {
	case *net.OpError:
		switch val.Op {
		case "dial":
			return fmt.Errorf("unknown host for server connection")
		case "read":
			return fmt.Errorf("connection to server had been refused")
		default:
			return fmt.Errorf("failed to connect to the server")
		}
	case syscall.Errno:
		switch val {
		case syscall.ECONNREFUSED:
			return fmt.Errorf("connection to server had been refused")
		default:
			return fmt.Errorf("failed to connect to the server")
		}
	case net.Error:
		return fmt.Errorf("failed to connect to the server")
	default:
		return err
	}
}
