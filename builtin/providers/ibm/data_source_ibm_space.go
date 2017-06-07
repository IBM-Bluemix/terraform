package ibm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceIBMSpace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMSpaceRead,

		Schema: map[string]*schema.Schema{
			"space": {
				Description: "Space name, for example dev",
				Type:        schema.TypeString,
				Required:    true,
			},

			"org": {
				Description: "The org this space belongs to",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceIBMSpaceRead(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).CFAPI()
	if err != nil {
		return err
	}
	orgAPI := cfClient.Organizations()
	spaceAPI := cfClient.Spaces()

	space := d.Get("space").(string)
	org := d.Get("org").(string)

	orgFields, err := orgAPI.FindByName(org)
	if err != nil {
		return fmt.Errorf("Error retrieving org: %s", err)
	}
	spaceFields, err := spaceAPI.FindByNameInOrg(orgFields.GUID, space)
	if err != nil {
		return fmt.Errorf("Error retrieving space: %s", err)
	}

	d.SetId(spaceFields.GUID)

	return nil
}
