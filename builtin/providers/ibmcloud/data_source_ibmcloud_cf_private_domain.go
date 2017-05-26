package ibmcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceIBMCloudCfPrivateDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCloudCfPrivateDomainRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the private domain",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceIBMCloudCfPrivateDomainRead(d *schema.ResourceData, meta interface{}) error {

	domainName := d.Get("name").(string)

	privateDomain := meta.(ClientSession).CloudFoundryPrivateDomainClient()
	prdomain, err := privateDomain.FindByName(domainName)
	if err != nil {
		return fmt.Errorf("Error retrieving domain: %s", err)
	}
	d.SetId(prdomain.GUID)
	return nil

}
