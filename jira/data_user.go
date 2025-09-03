package jira

import (
	"io"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func dataUser() *schema.Resource {
	return &schema.Resource{
		Read: dataUserRead,

		Description: "Reads a user",
		Schema: map[string]*schema.Schema{
			"query": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name or Email of the user",
			},
			"id": {
				Type:        schema.TypeString,
				Description: "ID",
				Computed:    true,
			},
			"self": {
				Type:        schema.TypeString,
				Description: "self",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "name",
				Computed:    true,
			},
			"email_address": {
				Type:        schema.TypeString,
				Description: "email address",
				Computed:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "display name",
				Computed:    true,
			},
		},
	}
}

func dataUserRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	query := d.Get("query").(string)

	users, res, err := config.jiraClient.User.Find(query)
	if err != nil {
		body, _ := io.ReadAll(res.Body)
		return errors.Wrapf(err, "finding jira user failed: %s", body)
	}

	if len(users) > 1 {
		emails := make([]string, len(users))
		for idx, user := range users {
			emails[idx] = user.EmailAddress
		}
		return errors.Wrapf(err, "query returned multiple users: %s", emails)
	}

	user := users[0]
	d.SetId(user.AccountID)
	d.Set("self", user.Self)
	d.Set("name", user.Name)
	d.Set("email_address", user.EmailAddress)
	d.Set("display_name", user.DisplayName)

	return nil
}
