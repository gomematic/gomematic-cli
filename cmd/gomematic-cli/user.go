package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/go-openapi/strfmt"
	"github.com/gomematic/gomematic-go/gomematic/user"
	"github.com/gomematic/gomematic-go/models"
	"gopkg.in/urfave/cli.v2"
)

// tmplUserList represents a row within user listing.
var tmplUserList = "Slug: \x1b[33m{{ .Slug }} \x1b[0m" + `
ID: {{ .ID }}
Username: {{ .Username }}
`

// tmplUserShow represents a user within details view.
var tmplUserShow = "Slug: \x1b[33m{{ .Slug }} \x1b[0m" + `
ID: {{ .ID }}
Username: {{ .Username }}
Email: {{ .Email }}
Active: {{ .Active }}
Admin: {{ .Admin }}
Created: {{ .CreatedAt }}
Updated: {{ .UpdatedAt }}
`

// tmplUserTeamList represents a row within user team listing.
var tmplUserTeamList = "Slug: \x1b[33m{{ .Team.Slug }} \x1b[0m" + `
ID: {{ .Team.ID }}
Name: {{ .Team.Name }}
Permission: {{ .Perm }}
`

// User provides the sub-command for the user API.
func User() *cli.Command {
	return &cli.Command{
		Name:  "user",
		Usage: "User related sub-commands",
		Subcommands: []*cli.Command{
			{
				Name:      "list",
				Aliases:   []string{"ls"},
				Usage:     "list all users",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:   "format",
						Value:  tmplUserList,
						Usage:  "custom output format",
						Hidden: true,
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, UserList)
				},
			},
			{
				Name:      "show",
				Usage:     "show an user",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "user id or slug",
					},
					&cli.StringFlag{
						Name:   "format",
						Value:  tmplUserShow,
						Usage:  "custom output format",
						Hidden: true,
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, UserShow)
				},
			},
			{
				Name:      "delete",
				Aliases:   []string{"rm"},
				Usage:     "delete an user",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "user id or slug",
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, UserDelete)
				},
			},
			{
				Name:      "update",
				Usage:     "update an user",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "user id or slug",
					},
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
					&cli.BoolFlag{
						Name:  "active",
						Usage: "mark user as active",
					},
					&cli.BoolFlag{
						Name:  "admin",
						Usage: "mark user as admin",
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, UserUpdate)
				},
			},
			{
				Name:      "create",
				Usage:     "create an user",
				ArgsUsage: " ",
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
					&cli.BoolFlag{
						Name:  "active",
						Usage: "mark user as active",
					},
					&cli.BoolFlag{
						Name:  "admin",
						Usage: "mark user as admin",
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, UserCreate)
				},
			},
			{
				Name:  "team",
				Usage: "team assignments",
				Subcommands: []*cli.Command{
					{
						Name:      "list",
						Aliases:   []string{"ls"},
						Usage:     "list assigned teams for a user",
						ArgsUsage: " ",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "id, i",
								Value: "",
								Usage: "user id or slug",
							},
							&cli.StringFlag{
								Name:   "format",
								Value:  tmplUserTeamList,
								Usage:  "custom output format",
								Hidden: true,
							},
						},
						Action: func(c *cli.Context) error {
							return Handle(c, UserTeamList)
						},
					},
					{
						Name:      "append",
						Usage:     "append a team to an user",
						ArgsUsage: " ",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "id, i",
								Value: "",
								Usage: "user id or slug",
							},
							&cli.StringFlag{
								Name:  "team, t",
								Value: "",
								Usage: "team id or slug",
							},
							&cli.StringFlag{
								Name:  "perm",
								Value: "user",
								Usage: "permission, can be user, admin or owner",
							},
						},
						Action: func(c *cli.Context) error {
							return Handle(c, UserTeamAppend)
						},
					},
					{
						Name:      "perm",
						Usage:     "update user team permissions",
						ArgsUsage: " ",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "id, i",
								Value: "",
								Usage: "user id or slug to update",
							},
							&cli.StringFlag{
								Name:  "team, t",
								Value: "",
								Usage: "team id or slug to update",
							},
							&cli.StringFlag{
								Name:  "perm",
								Value: "user",
								Usage: "permission, can be user, admin or owner",
							},
						},
						Action: func(c *cli.Context) error {
							return Handle(c, UserTeamPerm)
						},
					},
					{
						Name:      "remove",
						Aliases:   []string{"rm"},
						Usage:     "remove a team from an user",
						ArgsUsage: " ",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "id, i",
								Value: "",
								Usage: "user id or slug to remove from",
							},
							&cli.StringFlag{
								Name:  "team, t",
								Value: "",
								Usage: "team id or slug to remove",
							},
						},
						Action: func(c *cli.Context) error {
							return Handle(c, UserTeamRemove)
						},
					},
				},
			},
		},
	}
}

// UserList provides the sub-command to list all users.
func UserList(c *cli.Context, client *Client) error {
	resp, err := client.User.ListUsers(
		user.NewListUsersParams(),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *user.ListUsersForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *user.ListUsersDefault:
			return fmt.Errorf(*val.Payload.Message)
		default:
			return PrettyError(err)
		}
	}

	if len(resp.Payload) == 0 {
		fmt.Fprintln(os.Stderr, "empty result")
		return nil
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

	for _, record := range resp.Payload {
		if err := tmpl.Execute(os.Stdout, record); err != nil {
			return err
		}
	}

	return nil
}

// UserShow provides the sub-command to show user details.
func UserShow(c *cli.Context, client *Client) error {
	resp, err := client.User.ShowUser(
		user.NewShowUserParams().WithUserID(GetIdentifierParam(c)),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *user.ShowUserForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *user.ShowUserNotFound:
			return fmt.Errorf(*val.Payload.Message)
		case *user.ShowUserDefault:
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

// UserDelete provides the sub-command to delete a user.
func UserDelete(c *cli.Context, client *Client) error {
	resp, err := client.User.DeleteUser(
		user.NewDeleteUserParams().WithUserID(GetIdentifierParam(c)),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *user.DeleteUserForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *user.DeleteUserNotFound:
			return fmt.Errorf(*val.Payload.Message)
		case *user.DeleteUserBadRequest:
			return fmt.Errorf(*val.Payload.Message)
		case *user.DeleteUserDefault:
			return fmt.Errorf(*val.Payload.Message)
		default:
			return PrettyError(err)
		}
	}

	fmt.Fprintln(os.Stderr, *resp.Payload.Message)
	return nil
}

// UserUpdate provides the sub-command to update a user.
func UserUpdate(c *cli.Context, client *Client) error {
	resp, err := client.User.ShowUser(
		user.NewShowUserParams().WithUserID(GetIdentifierParam(c)),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *user.ShowUserForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *user.ShowUserNotFound:
			return fmt.Errorf(*val.Payload.Message)
		case *user.ShowUserDefault:
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

	if c.IsSet("active") {
		val := c.Bool("active")
		record.Active = &val
		changed = true
	}

	if c.IsSet("admin") {
		val := c.Bool("admin")
		record.Admin = &val
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

		_, err := client.User.UpdateUser(
			user.NewUpdateUserParams().WithUserID(record.ID.String()).WithUser(record),
			client.AuthInfo,
		)

		if err != nil {
			switch val := err.(type) {
			case *user.UpdateUserForbidden:
				return fmt.Errorf(*val.Payload.Message)
			case *user.UpdateUserNotFound:
				return fmt.Errorf(*val.Payload.Message)

				//
				//
				//

			case *user.UpdateUserUnprocessableEntity:
				return fmt.Errorf(*val.Payload.Message)

				//
				//
				//

			case *user.UpdateUserDefault:
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

// UserCreate provides the sub-command to create a user.
func UserCreate(c *cli.Context, client *Client) error {
	record := &models.User{}

	if val := c.String("slug"); c.IsSet("slug") && val != "" {
		record.Slug = &val
	}

	if val := c.String("email"); c.IsSet("email") && val != "" {
		record.Email = &val
	} else {
		return fmt.Errorf("you must provide an email")
	}

	if val := c.String("username"); c.IsSet("username") && val != "" {
		record.Username = &val
	} else {
		return fmt.Errorf("you must provide an username")
	}

	if val := c.String("password"); c.IsSet("password") && val != "" {
		password := strfmt.Password(val)
		record.Password = &password
	} else {
		return fmt.Errorf("you must provide a password")
	}

	if c.IsSet("active") {
		val := c.Bool("active")
		record.Active = &val
	}

	if c.IsSet("admin") {
		val := c.Bool("admin")
		record.Admin = &val
	}

	if err := record.Validate(strfmt.Default); err != nil {

		//
		//
		//

		return err

		//
		//
		//

	}

	_, err := client.User.CreateUser(
		user.NewCreateUserParams().WithUser(record),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *user.CreateUserForbidden:
			return fmt.Errorf(*val.Payload.Message)

		//
		//
		//

		case *user.CreateUserUnprocessableEntity:
			return fmt.Errorf(*val.Payload.Message)

		//
		//
		//

		case *user.CreateUserDefault:
			return fmt.Errorf(*val.Payload.Message)
		default:
			return PrettyError(err)
		}
	}

	fmt.Fprintln(os.Stderr, "successfully created")
	return nil
}

// UserTeamList provides the sub-command to list teams of the user.
func UserTeamList(c *cli.Context, client *Client) error {
	resp, err := client.User.ListUserTeams(
		user.NewListUserTeamsParams().WithUserID(GetIdentifierParam(c)),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *user.ListUserTeamsForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *user.ListUserTeamsNotFound:
			return fmt.Errorf(*val.Payload.Message)
		case *user.ListUserTeamsDefault:
			return fmt.Errorf(*val.Payload.Message)
		default:
			return PrettyError(err)
		}
	}

	if len(resp.Payload) == 0 {
		fmt.Fprintln(os.Stderr, "empty result")
		return nil
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

	for _, record := range resp.Payload {
		if err := tmpl.Execute(os.Stdout, record); err != nil {
			return err
		}
	}

	return nil
}

// UserTeamAppend provides the sub-command to append a team to the user.
func UserTeamAppend(c *cli.Context, client *Client) error {
	teamID := GetTeamParam(c)
	perm := GetPermParam(c)

	resp, err := client.User.AppendUserToTeam(
		user.NewAppendUserToTeamParams().WithUserID(GetIdentifierParam(c)).WithUserTeam(
			&models.UserTeamParams{
				Team: &teamID,
				Perm: &perm,
			},
		),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *user.AppendUserToTeamForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *user.AppendUserToTeamNotFound:
			return fmt.Errorf(*val.Payload.Message)
		case *user.AppendUserToTeamPreconditionFailed:
			return fmt.Errorf(*val.Payload.Message)

		case *user.AppendUserToTeamUnprocessableEntity:
			return fmt.Errorf(*val.Payload.Message)

		case *user.AppendUserToTeamDefault:
			return fmt.Errorf(*val.Payload.Message)
		default:
			return PrettyError(err)
		}
	}

	fmt.Fprintln(os.Stderr, *resp.Payload.Message)
	return nil
}

// UserTeamPerm provides the sub-command to update user team permissions.
func UserTeamPerm(c *cli.Context, client *Client) error {
	teamID := GetTeamParam(c)
	perm := GetPermParam(c)

	resp, err := client.User.PermitUserTeam(
		user.NewPermitUserTeamParams().WithUserID(GetIdentifierParam(c)).WithUserTeam(
			&models.UserTeamParams{
				Team: &teamID,
				Perm: &perm,
			},
		),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *user.PermitUserTeamForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *user.PermitUserTeamNotFound:
			return fmt.Errorf(*val.Payload.Message)
		case *user.PermitUserTeamPreconditionFailed:
			return fmt.Errorf(*val.Payload.Message)

		case *user.PermitUserTeamUnprocessableEntity:
			return fmt.Errorf(*val.Payload.Message)

		case *user.PermitUserTeamDefault:
			return fmt.Errorf(*val.Payload.Message)
		default:
			return PrettyError(err)
		}
	}

	fmt.Fprintln(os.Stderr, *resp.Payload.Message)
	return nil
}

// UserTeamRemove provides the sub-command to remove a team from the user.
func UserTeamRemove(c *cli.Context, client *Client) error {
	teamID := GetTeamParam(c)
	perm := "user"

	resp, err := client.User.DeleteUserFromTeam(
		user.NewDeleteUserFromTeamParams().WithUserID(GetIdentifierParam(c)).WithUserTeam(
			&models.UserTeamParams{
				Team: &teamID,
				Perm: &perm,
			},
		),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *user.DeleteUserFromTeamForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *user.DeleteUserFromTeamNotFound:
			return fmt.Errorf(*val.Payload.Message)
		case *user.DeleteUserFromTeamPreconditionFailed:
			return fmt.Errorf(*val.Payload.Message)
		case *user.DeleteUserFromTeamDefault:
			return fmt.Errorf(*val.Payload.Message)
		default:
			return PrettyError(err)
		}
	}

	fmt.Fprintln(os.Stderr, *resp.Payload.Message)
	return nil
}
