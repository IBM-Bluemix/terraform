package ibmcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceIBMCloudCfDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCloudCfDomainRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the domain",
				Type:        schema.TypeString,
				Required:    true,
			},
			"domain_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the domain. Accepted values are shared or private",
			},
		},
	}
}

func dataSourceIBMCloudCfDomainRead(d *schema.ResourceData, meta interface{}) error {

	domainName := d.Get("name").(string)
	domainType := d.Get("domain_type").(string)

	if domainType == "shared" {
		sharedDomain := meta.(ClientSession).CloudFoundrySharedDomainClient()
		shdomain, err := sharedDomain.FindByName(domainName)
		if err != nil {
			return fmt.Errorf("Error retrieving domain: %s", err)
		}
		d.SetId(shdomain.GUID)
		return nil
	}

	if domainType == "private" {
		privateDomain := meta.(ClientSession).CloudFoundryPrivateDomainClient()
		prdomain, err := privateDomain.FindByName(domainName)
		if err != nil {
			return fmt.Errorf("Error retrieving domain: %s", err)
		}
		d.SetId(prdomain.GUID)
		return nil
	}

	sharedDomain := meta.(ClientSession).CloudFoundrySharedDomainClient()
	shdomain, err := sharedDomain.FindByName(domainName)
	if err != nil {
		privateDomain := meta.(ClientSession).CloudFoundryPrivateDomainClient()
		prdomain, err := privateDomain.FindByName(domainName)
		if err != nil {
			return fmt.Errorf("Error retrieving domain: %s", domainName)
		}
		d.SetId(prdomain.GUID)
		return nil
	}

	d.SetId(shdomain.GUID)
	return nil

}
