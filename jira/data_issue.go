package jira

import (
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func dataIssue() *schema.Resource {
	return &schema.Resource{
		Read: dataIssueRead,

		Description: "Reads an issue",

		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Issue key",
			},
			"self": {
				Type:        schema.TypeString,
				Description: "self",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "description",
				Computed:    true,
			},
			"summary": {
				Type:        schema.TypeString,
				Description: "summary",
				Computed:    true,
			},
			"fields": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				Description: "fields",
				Computed:    true,
			},
		},
	}
}

func dataIssueRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	key := d.Get("key").(string)

	issue, res, err := config.jiraClient.Issue.Get(key, nil)
	if err != nil {
		body, _ := io.ReadAll(res.Body)
		return errors.Wrapf(err, "getting jira issue failed: %s", body)
	}

	d.SetId(issue.ID)
	d.Set("self", issue.Self)
	d.Set("description", issue.Fields.Description)
	d.Set("summary", issue.Fields.Summary)

	unknowns := map[string]string{}
	for field, value := range issue.Fields.Unknowns {
		unknowns[field] = fmt.Sprintf("%v", value)
	}
	d.Set("fields", unknowns)

	return nil
}
