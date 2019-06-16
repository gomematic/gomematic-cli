package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/go-openapi/strfmt"
	"github.com/gomematic/gomematic-go/gomematic/team"
	"github.com/gomematic/gomematic-go/models"
	"gopkg.in/urfave/cli.v2"
)

// tmplTeamList represents a row within user listing.
var tmplTeamList = "Slug: \x1b[33m{{ .Slug }} \x1b[0m" + `
ID: {{ .ID }}
Name: {{ .Name }}
`

// tmplTeamShow represents a user within details view.
var tmplTeamShow = "Slug: \x1b[33m{{ .Slug }} \x1b[0m" + `
ID: {{ .ID }}
Name: {{ .Name }}
Created: {{ .CreatedAt }}
Updated: {{ .UpdatedAt }}
`

// tmplTeamUserList represents a row within team user listing.
var tmplTeamUserList = "Slug: \x1b[33m{{ .User.Slug }} \x1b[0m" + `
ID: {{ .User.ID }}
Username: {{ .User.Username }}
Permission: {{ .Perm }}
`

// Team provides the sub-command for the team API.
func Team() *cli.Command {
	return &cli.Command{
		Name:  "team",
		Usage: "team commands",
		Subcommands: []*cli.Command{
			{
				Name:      "list",
				Aliases:   []string{"ls"},
				Usage:     "list all teams",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:   "format",
						Value:  tmplTeamList,
						Usage:  "custom output format",
						Hidden: true,
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, TeamList)
				},
			},
			{
				Name:      "show",
				Usage:     "show a team",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "team id or slug",
					},
					&cli.StringFlag{
						Name:   "format",
						Value:  tmplTeamShow,
						Usage:  "custom output format",
						Hidden: true,
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, TeamShow)
				},
			},
			{
				Name:      "delete",
				Aliases:   []string{"rm"},
				Usage:     "delete a team",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "team id or slug",
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, TeamDelete)
				},
			},
			{
				Name:      "update",
				Usage:     "update a team",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id, i",
						Value: "",
						Usage: "team id or slug",
					},
					&cli.StringFlag{
						Name:  "slug",
						Value: "",
						Usage: "provide a slug",
					},
					&cli.StringFlag{
						Name:  "name",
						Value: "",
						Usage: "provide a name",
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, TeamUpdate)
				},
			},
			{
				Name:      "create",
				Usage:     "create a team",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "slug",
						Value: "",
						Usage: "provide a slug",
					},
					&cli.StringFlag{
						Name:  "name",
						Value: "",
						Usage: "provide a name",
					},
				},
				Action: func(c *cli.Context) error {
					return Handle(c, TeamCreate)
				},
			},
			{
				Name:  "user",
				Usage: "user assignments",
				Subcommands: []*cli.Command{
					{
						Name:      "list",
						Aliases:   []string{"ls"},
						Usage:     "list assigned users for a team",
						ArgsUsage: " ",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "id, i",
								Value: "",
								Usage: "team id or slug",
							},
							&cli.StringFlag{
								Name:   "format",
								Value:  tmplTeamUserList,
								Usage:  "custom output format",
								Hidden: true,
							},
						},
						Action: func(c *cli.Context) error {
							return Handle(c, TeamUserList)
						},
					},
					{
						Name:      "append",
						Usage:     "append a user to team",
						ArgsUsage: " ",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "id, i",
								Value: "",
								Usage: "team id or slug",
							},
							&cli.StringFlag{
								Name:  "user, u",
								Value: "",
								Usage: "user id or slug",
							},
							&cli.StringFlag{
								Name:  "perm",
								Value: "user",
								Usage: "permission, can be user, admin or owner",
							},
						},
						Action: func(c *cli.Context) error {
							return Handle(c, TeamUserAppend)
						},
					},
					{
						Name:      "perm",
						Usage:     "update team user permissions",
						ArgsUsage: " ",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "id, i",
								Value: "",
								Usage: "team id or slug",
							},
							&cli.StringFlag{
								Name:  "user, u",
								Value: "",
								Usage: "user id or slug",
							},
							&cli.StringFlag{
								Name:  "perm",
								Value: "user",
								Usage: "permission, can be user, admin or owner",
							},
						},
						Action: func(c *cli.Context) error {
							return Handle(c, TeamUserPerm)
						},
					},
					{
						Name:      "remove",
						Aliases:   []string{"rm"},
						Usage:     "remove a user from a team",
						ArgsUsage: " ",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "id, i",
								Value: "",
								Usage: "team id or slug",
							},
							&cli.StringFlag{
								Name:  "user, u",
								Value: "",
								Usage: "user id or slug",
							},
						},
						Action: func(c *cli.Context) error {
							return Handle(c, TeamUserRemove)
						},
					},
				},
			},
		},
	}
}

// TeamList provides the sub-command to list all teams.
func TeamList(c *cli.Context, client *Client) error {
	resp, err := client.Team.ListTeams(
		team.NewListTeamsParams(),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *team.ListTeamsForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *team.ListTeamsDefault:
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

// TeamShow provides the sub-command to show team details.
func TeamShow(c *cli.Context, client *Client) error {
	resp, err := client.Team.ShowTeam(
		team.NewShowTeamParams().WithTeamID(GetIdentifierParam(c)),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *team.ShowTeamForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *team.ShowTeamNotFound:
			return fmt.Errorf(*val.Payload.Message)
		case *team.ShowTeamDefault:
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

// TeamDelete provides the sub-command to delete a team.
func TeamDelete(c *cli.Context, client *Client) error {
	resp, err := client.Team.DeleteTeam(
		team.NewDeleteTeamParams().WithTeamID(GetIdentifierParam(c)),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *team.DeleteTeamForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *team.DeleteTeamNotFound:
			return fmt.Errorf(*val.Payload.Message)
		case *team.DeleteTeamBadRequest:
			return fmt.Errorf(*val.Payload.Message)
		case *team.DeleteTeamDefault:
			return fmt.Errorf(*val.Payload.Message)
		default:
			return PrettyError(err)
		}
	}

	fmt.Fprintln(os.Stderr, *resp.Payload.Message)
	return nil
}

// TeamUpdate provides the sub-command to update a team.
func TeamUpdate(c *cli.Context, client *Client) error {
	resp, err := client.Team.ShowTeam(
		team.NewShowTeamParams().WithTeamID(GetIdentifierParam(c)),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *team.ShowTeamForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *team.ShowTeamNotFound:
			return fmt.Errorf(*val.Payload.Message)
		case *team.ShowTeamDefault:
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

	if val := c.String("name"); c.IsSet("name") && val != *record.Name {
		record.Name = &val
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

		_, err := client.Team.UpdateTeam(
			team.NewUpdateTeamParams().WithTeamID(record.ID.String()).WithTeam(record),
			client.AuthInfo,
		)

		if err != nil {
			switch val := err.(type) {
			case *team.UpdateTeamForbidden:
				return fmt.Errorf(*val.Payload.Message)
			case *team.UpdateTeamNotFound:
				return fmt.Errorf(*val.Payload.Message)

				//
				//
				//

			case *team.UpdateTeamUnprocessableEntity:
				return fmt.Errorf(*val.Payload.Message)

				//
				//
				//

			case *team.UpdateTeamDefault:
				return fmt.Errorf(*val.Payload.Message)
			default:
				return PrettyError(err)
			}
		}

		fmt.Fprintln(os.Stderr, "successfully updated")
	} else {
		fmt.Fprintln(os.Stderr, "nothing to update...")
	}

	return nil
}

// TeamCreate provides the sub-command to create a team.
func TeamCreate(c *cli.Context, client *Client) error {
	record := &models.Team{}

	if val := c.String("slug"); c.IsSet("slug") && val != "" {
		record.Slug = &val
	}

	if val := c.String("name"); c.IsSet("name") && val != "" {
		record.Name = &val
	} else {
		return fmt.Errorf("you must provide a name")
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

	_, err := client.Team.CreateTeam(
		team.NewCreateTeamParams().WithTeam(record),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *team.CreateTeamForbidden:
			return fmt.Errorf(*val.Payload.Message)

		//
		//
		//

		case *team.CreateTeamUnprocessableEntity:
			return fmt.Errorf(*val.Payload.Message)

		//
		//
		//

		case *team.CreateTeamDefault:
			return fmt.Errorf(*val.Payload.Message)
		default:
			return PrettyError(err)
		}
	}

	fmt.Fprintln(os.Stderr, "successfully created")
	return nil
}

// TeamUserList provides the sub-command to list users of the team.
func TeamUserList(c *cli.Context, client *Client) error {
	resp, err := client.Team.ListTeamUsers(
		team.NewListTeamUsersParams().WithTeamID(GetIdentifierParam(c)),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *team.ListTeamUsersForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *team.ListTeamUsersNotFound:
			return fmt.Errorf(*val.Payload.Message)
		case *team.ListTeamUsersDefault:
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

// TeamUserAppend provides the sub-command to append a user to the team.
func TeamUserAppend(c *cli.Context, client *Client) error {
	userID := GetUserParam(c)
	perm := GetPermParam(c)

	resp, err := client.Team.AppendTeamToUser(
		team.NewAppendTeamToUserParams().WithTeamID(GetIdentifierParam(c)).WithTeamUser(
			&models.TeamUserParams{
				User: &userID,
				Perm: &perm,
			},
		),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *team.AppendTeamToUserForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *team.AppendTeamToUserNotFound:
			return fmt.Errorf(*val.Payload.Message)
		case *team.AppendTeamToUserPreconditionFailed:
			return fmt.Errorf(*val.Payload.Message)

		case *team.AppendTeamToUserUnprocessableEntity:
			return fmt.Errorf(*val.Payload.Message)

		case *team.AppendTeamToUserDefault:
			return fmt.Errorf(*val.Payload.Message)
		default:
			return PrettyError(err)
		}
	}

	fmt.Fprintln(os.Stderr, *resp.Payload.Message)
	return nil
}

// TeamUserPerm provides the sub-command to update team user permissions.
func TeamUserPerm(c *cli.Context, client *Client) error {
	userID := GetUserParam(c)
	perm := GetPermParam(c)

	resp, err := client.Team.PermitTeamUser(
		team.NewPermitTeamUserParams().WithTeamID(GetIdentifierParam(c)).WithTeamUser(
			&models.TeamUserParams{
				User: &userID,
				Perm: &perm,
			},
		),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *team.PermitTeamUserForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *team.PermitTeamUserNotFound:
			return fmt.Errorf(*val.Payload.Message)
		case *team.PermitTeamUserPreconditionFailed:
			return fmt.Errorf(*val.Payload.Message)

		case *team.PermitTeamUserUnprocessableEntity:
			return fmt.Errorf(*val.Payload.Message)

		case *team.PermitTeamUserDefault:
			return fmt.Errorf(*val.Payload.Message)
		default:
			return PrettyError(err)
		}
	}

	fmt.Fprintln(os.Stderr, *resp.Payload.Message)
	return nil
}

// TeamUserRemove provides the sub-command to remove a user from the team.
func TeamUserRemove(c *cli.Context, client *Client) error {
	userID := GetUserParam(c)
	perm := "user"

	resp, err := client.Team.DeleteTeamFromUser(
		team.NewDeleteTeamFromUserParams().WithTeamID(GetIdentifierParam(c)).WithTeamUser(
			&models.TeamUserParams{
				User: &userID,
				Perm: &perm,
			},
		),
		client.AuthInfo,
	)

	if err != nil {
		switch val := err.(type) {
		case *team.DeleteTeamFromUserForbidden:
			return fmt.Errorf(*val.Payload.Message)
		case *team.DeleteTeamFromUserNotFound:
			return fmt.Errorf(*val.Payload.Message)
		case *team.DeleteTeamFromUserPreconditionFailed:
			return fmt.Errorf(*val.Payload.Message)
		case *team.DeleteTeamFromUserDefault:
			return fmt.Errorf(*val.Payload.Message)
		default:
			return PrettyError(err)
		}
	}

	fmt.Fprintln(os.Stderr, *resp.Payload.Message)
	return nil
}
