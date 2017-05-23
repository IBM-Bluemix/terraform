package ibmcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceIBMCloudCfSpace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCloudCfSpaceRead,

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

func dataSourceIBMCloudCfSpaceRead(d *schema.ResourceData, meta interface{}) error {
	or := meta.(ClientSession).CloudFoundryOrgClient()
	sp := meta.(ClientSession).CloudFoundrySpaceClient()

	space := d.Get("space").(string)
	org := d.Get("org").(string)

	orgFields, err := or.FindByName(org)
	if err != nil {
		return fmt.Errorf("Error retrieving org: %s", err)
	}

	spaceFields, err := sp.FindByNameInOrg(orgFields.GUID, space)
	if err != nil {
		return fmt.Errorf("Error retrieving space: %s", err)
	}

	d.SetId(spaceFields.GUID)

	return nil
}