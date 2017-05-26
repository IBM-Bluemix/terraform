package ibmcloud

import (
	"fmt"

	v2 "github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIBMCloudCfPrivateDomain() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCloudCfPrivateDomainCreate,
		Read:     resourceIBMCloudCfPrivateDomainRead,
		Delete:   resourceIBMCloudCfPrivateDomainDelete,
		Exists:   resourceIBMCloudCfPrivateDomainExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The name of the domain",
				ValidateFunc: validateDomainName,
			},

			"org_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The organization that owns the domain.",
			},
		},
	}
}

func resourceIBMCloudCfPrivateDomainCreate(d *schema.ResourceData, meta interface{}) error {
	prdomainClient := meta.(ClientSession).CloudFoundryPrivateDomainClient()
	name := d.Get("name").(string)
	orgGUID := d.Get("org_guid").(string)

	params := v2.PrivateDomainRequest{
		Name:    name,
		OrgGUID: orgGUID,
	}

	prdomain, err := prdomainClient.Create(params)
	if err != nil {
		return fmt.Errorf("Error creating private domain: %s", err)
	}

	d.SetId(prdomain.Metadata.GUID)

	return resourceIBMCloudCfPrivateDomainRead(d, meta)
}

func resourceIBMCloudCfPrivateDomainRead(d *schema.ResourceData, meta interface{}) error {
	prdomainClient := meta.(ClientSession).CloudFoundryPrivateDomainClient()
	prdomainGUID := d.Id()

	prdomain, err := prdomainClient.Get(prdomainGUID)
	if err != nil {
		return fmt.Errorf("Error retrieving private domain: %s", err)
	}
	d.Set("name", prdomain.Entity.Name)
	d.Set("org_guid", prdomain.Entity.OwningOrganizationGUID)

	return nil
}

func resourceIBMCloudCfPrivateDomainDelete(d *schema.ResourceData, meta interface{}) error {
	prdomainClient := meta.(ClientSession).CloudFoundryPrivateDomainClient()

	prdomainGUID := d.Id()

	err := prdomainClient.Delete(prdomainGUID, true)
	if err != nil {
		return fmt.Errorf("Error deleting private domain: %s", err)
	}

	d.SetId("")

	return nil
}

func resourceIBMCloudCfPrivateDomainExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	prdomainClient := meta.(ClientSession).CloudFoundryPrivateDomainClient()
	prdomainGUID := d.Id()

	prdomain, err := prdomainClient.Get(prdomainGUID)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return prdomain.Metadata.GUID == prdomainGUID, nil
}
