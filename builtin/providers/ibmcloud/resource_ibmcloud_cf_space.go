package ibmcloud

import (
	"fmt"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIBMCloudCfSpace() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCloudCfSpaceCreate,
		Read:     resourceIBMCloudCfSpaceRead,
		Update:   resourceIBMCloudCfSpaceUpdate,
		Delete:   resourceIBMCloudCfSpaceDelete,
		Exists:   resourceIBMCloudCfSpaceExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name for the space",
			},
			"org": {
				Description: "The org this space belongs to",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"space_quota": {
				Description: "The name of the Space Quota Definition",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceIBMCloudCfSpaceCreate(d *schema.ResourceData, meta interface{}) error {
	orgClient := meta.(ClientSession).CloudFoundryOrgClient()
	org := d.Get("org").(string)

	orgFields, err := orgClient.FindByName(org)
	if err != nil {
		return fmt.Errorf("Error retrieving org: %s", err)
	}

	spaceClient := meta.(ClientSession).CloudFoundrySpaceClient()
	name := d.Get("name").(string)

	var spaceQuotaGUID string

	if spaceQuota, ok := d.GetOk("space_quota"); ok {
		spaceQuotaClient := meta.(ClientSession).CloudFoundrySpaceQuotaClient()
		quota, err := spaceQuotaClient.FindByName(spaceQuota.(string), orgFields.GUID)
		if err != nil {
			return fmt.Errorf("Error retrieving space quota: %s", err)
		}
		spaceQuotaGUID = quota.GUID
	}

	space, err := spaceClient.Create(name, orgFields.GUID, spaceQuotaGUID)
	if err != nil {
		return fmt.Errorf("Error creating space: %s", err)
	}

	d.SetId(space.Metadata.GUID)
	return resourceIBMCloudCfSpaceRead(d, meta)
}

func resourceIBMCloudCfSpaceRead(d *schema.ResourceData, meta interface{}) error {
	spaceClient := meta.(ClientSession).CloudFoundrySpaceClient()
	spaceGUID := d.Id()

	_, err := spaceClient.Get(spaceGUID)
	if err != nil {
		return fmt.Errorf("Error retrieving space: %s", err)
	}
	return nil
}

func resourceIBMCloudCfSpaceUpdate(d *schema.ResourceData, meta interface{}) error {
	spaceClient := meta.(ClientSession).CloudFoundrySpaceClient()
	spaceGUID := d.Id()

	var name string
	if d.HasChange("name") {
		name = d.Get("name").(string)
	}

	_, err := spaceClient.Update(name, spaceGUID)
	if err != nil {
		return fmt.Errorf("Error updating space: %s", err)
	}

	return resourceIBMCloudCfSpaceRead(d, meta)
}

func resourceIBMCloudCfSpaceDelete(d *schema.ResourceData, meta interface{}) error {
	spaceClient := meta.(ClientSession).CloudFoundrySpaceClient()
	id := d.Id()

	err := spaceClient.Delete(id)
	if err != nil {
		return fmt.Errorf("Error deleting space: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceIBMCloudCfSpaceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	spaceClient := meta.(ClientSession).CloudFoundrySpaceClient()
	id := d.Id()

	space, err := spaceClient.Get(id)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return space.Metadata.GUID == id, nil
}
