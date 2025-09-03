package jira

import (
	"fmt"

	jira "github.com/andygrunwald/go-jira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

func dataComponent() *schema.Resource {
	return &schema.Resource{
		Read: dataComponentRead,

		Description: "Reads a project component",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the component",
			},

			"project_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project Key (for example PRJ)",
			},
		},
	}
}

func dataComponentRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	name := d.Get("name").(string)
	projectKey := d.Get("project_key").(string)

	urlStr := fmt.Sprintf("%s/%s/components", projectAPIEndpoint, projectKey)

	components := [](*jira.ProjectComponent){}
	err := request(config.jiraClient, "GET", urlStr, nil, &components)

	if err != nil {
		return errors.Wrap(err, "getting project components failed")
	}

	for _, component := range components {
		if component.Name == name {
			d.SetId(component.ID)
			return nil
		}
	}

	return errors.Wrap(ResourceNotFoundError, "could not find component in project")
}
