package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/go-openapi/strfmt"
	"github.com/gomematic/gomematic-go/gomematic/auth"
	"github.com/gomematic/gomematic-go/gomematic/profile"
	"github.com/gomematic/gomematic-go/models"
	"gopkg.in/urfave/cli.v2"
)

// tmplProfileLogin represents a expiring login token.
var tmplProfileLogin = "Token: \x1b[33m{{ .Token }} \x1b[0m" + `
Expires: {{ .ExpiresAt }}
`

// tmplProfileToken represents a permanent login token.
var tmplProfileToken = "Token: \x1b[33m{{ .Token }} \x1b[0m" + `
`

// tmplProfileShow represents a profile within details view.
var tmplProfileShow = "Slug: \x1b[33m{{ .Slug }} \x1b[0m" + `
ID: {{ .ID }}
Username: {{ .Username }}
Email: {{ .Email }}
Active: {{ .Active }}
Admin: {{ .Admin }}
Created: {{ .CreatedAt }}
Updated: {{ .UpdatedAt }}
`

// Profile provides the sub-command for the profile API.
func Profile() *cli.Command {
	return &cli.Command{
		Name:  "profile",
		Usage: "profile commands",
		Subcommands: []*cli.Command{
			{
				Name:  "login",
				Usage: "login by credentials",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "username",
						Value: "",
						Usage: "username for authentication",
					},
					&cli.StringFlag{
						Name:  "password",
						Value: "",
						Usage: "password for authentication",
					},
					&cli.StringFlag{
						Name:   "format",
						Value:  tmplProfileLogin,
						Usage:  "custom output format",
						Hidden: true,
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, ProfileLogin)
				},
			},
			{
				Name:  "token",
				Usage: "show your token",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:   "format",
						Value:  tmplProfileToken,
						Usage:  "custom output format",
						Hidden: true,
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, ProfileToken)
				},
			},
			{
				Name:  "show",
				Usage: "show profile details",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:   "format",
						Value:  tmplProfileShow,
						Usage:  "custom output format",
						Hidden: true,
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, ProfileShow)
				},
			},
			{
				Name:  "update",
				Usage: "update profile details",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "slug",
						Value: "",
						Usage: "provide a slug",
					},
					&cli.StringFlag{
						Name:  "email",
						Value: "",
						Usage: "provide an email",
					},
					&cli.StringFlag{
						Name:  "username",
						Value: "",
						Usage: "provide an username",
					},
					&cli.StringFlag{
						Name:  "password",
						Value: "",
						Usage: "provide a password",
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, ProfileUpdate)
				},
			},
		},
	}
}

// ProfileLogin provides the sub-command to login by credentials.
func ProfileLogin(c *cli.Context, client *Client) error {
	if !c.IsSet("username") {
		return fmt.Errorf("please provide a username")
	}

	if !c.IsSet("password") {
		return fmt.Errorf("please provide a password")
	}

	username := c.String("username")
	password := strfmt.Password(c.String("password"))

	resp, err := client.Auth.LoginUser(
		auth.NewLoginUserParams().WithAuthLogin(&models.AuthLogin{
			Username: &username,
			Password: &password,
		}),
	)

	if err != nil {
		switch val := err.(type) {
		case *auth.LoginUserUnauthorized:
			return fmt.Errorf(*val.Payload.Message)
		case *auth.LoginUserDefault:
			return fmt.Errorf(*val.Payload.Message)
		default:
			return PrettyError(err)
		}
	}

	tmpl, err := template.New(
		"_",
	).Funcs(
		globalFuncMap,
	).Funcs(
		sprigFuncMap,
	).Parse(
		fmt.Sprintln(c.String("format")),
	)

	if err != nil {
		return err
	}

	return tmpl.Execute(os.Stdout, resp.Payload)
}

// ProfileToken provides the sub-command to show your token.
func ProfileToken(c *cli.Context, client *Client) error {
	resp, err := client.Profile.TokenProfile(
		profile.NewTokenProfileParams(),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *profile.TokenProfileForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *profile.TokenProfileInternalServerError:
			return fmt.Errorf(*val.Payload.Message)
		case *profile.TokenProfileDefault:
			return fmt.Errorf(*val.Payload.Message)
		default:
			return PrettyError(err)
		}
	}

	tmpl, err := template.New(
		"_",
	).Funcs(
		globalFuncMap,
	).Funcs(
		sprigFuncMap,
	).Parse(
		fmt.Sprintln(c.String("format")),
	)

	if err != nil {
		return err
	}

	return tmpl.Execute(os.Stdout, resp.Payload)
}

// ProfileShow provides the sub-command to show profile details.
func ProfileShow(c *cli.Context, client *Client) error {
	resp, err := client.Profile.ShowProfile(
		profile.NewShowProfileParams(),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *profile.ShowProfileForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *profile.ShowProfileDefault:
			return fmt.Errorf(*val.Payload.Message)
		default:
			return PrettyError(err)
		}
	}

	tmpl, err := template.New(
		"_",
	).Funcs(
		globalFuncMap,
	).Funcs(
		sprigFuncMap,
	).Parse(
		fmt.Sprintln(c.String("format")),
	)

	if err != nil {
		return err
	}

	return tmpl.Execute(os.Stdout, resp.Payload)
}

// ProfileUpdate provides the sub-command to update the profile.
func ProfileUpdate(c *cli.Context, client *Client) error {
	resp, err := client.Profile.ShowProfile(
		profile.NewShowProfileParams(),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *profile.ShowProfileForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *profile.ShowProfileDefault:
			return fmt.Errorf(*val.Payload.Message)
		default:
			return PrettyError(err)
		}
	}

	record := resp.Payload
	changed := false

	if val := c.String("slug"); c.IsSet("slug") && val != *record.Slug {
		record.Slug = &val
		changed = true
	}

	if val := c.String("email"); c.IsSet("email") && val != *record.Email {
		record.Email = &val
		changed = true
	}

	if val := c.String("username"); c.IsSet("username") && val != *record.Username {
		record.Username = &val
		changed = true
	}

	if val := c.String("password"); c.IsSet("password") {
		password := strfmt.Password(val)
		record.Password = &password
		changed = true
	}

	if changed {
		if err := record.Validate(strfmt.Default); err != nil {

			//
			//
			//

			return err

			//
			//
			//

		}

		_, err := client.Profile.UpdateProfile(
			profile.NewUpdateProfileParams().WithProfile(record),
			client.AuthInfo,
		)

		if err != nil {
			switch val := err.(type) {
			case *profile.UpdateProfileForbidden:
				return fmt.Errorf(*val.Payload.Message)

				//
				//
				//

			case *profile.UpdateProfileUnprocessableEntity:
				return fmt.Errorf(*val.Payload.Message)

				//
				//
				//

			case *profile.UpdateProfileDefault:
				return fmt.Errorf(*val.Payload.Message)
			default:
				return PrettyError(err)
			}
		}

		fmt.Fprintln(os.Stderr, "successfully update")
	} else {
		fmt.Fprintln(os.Stderr, "nothing to update...")
	}

	return nil
}
