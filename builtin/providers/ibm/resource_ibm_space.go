package ibm

import (
	"fmt"

	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/helpers"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIBMSpace() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMSpaceCreate,
		Read:     resourceIBMSpaceRead,
		Update:   resourceIBMSpaceUpdate,
		Delete:   resourceIBMSpaceDelete,
		Exists:   resourceIBMSpaceExists,
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

func resourceIBMSpaceCreate(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).CFAPI()
	if err != nil {
		return err
	}
	org := d.Get("org").(string)
	name := d.Get("name").(string)

	req := cfv2.SpaceCreateRequest{
		Name: name,
	}

	orgFields, err := cfClient.Organizations().FindByName(org)
	if err != nil {
		return fmt.Errorf("Error retrieving org: %s", err)
	}
	req.OrgGUID = orgFields.GUID

	if spaceQuota, ok := d.GetOk("space_quota"); ok {
		quota, err := cfClient.SpaceQuotas().FindByName(spaceQuota.(string), orgFields.GUID)
		if err != nil {
			return fmt.Errorf("Error retrieving space quota: %s", err)
		}
		req.SpaceQuotaGUID = quota.GUID
	}

	space, err := cfClient.Spaces().Create(req)
	if err != nil {
		return fmt.Errorf("Error creating space: %s", err)
	}

	d.SetId(space.Metadata.GUID)
	return resourceIBMSpaceRead(d, meta)
}

func resourceIBMSpaceRead(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).CFAPI()
	if err != nil {
		return err
	}
	spaceGUID := d.Id()

	_, err = cfClient.Spaces().Get(spaceGUID)
	if err != nil {
		return fmt.Errorf("Error retrieving space: %s", err)
	}
	return nil
}

func resourceIBMSpaceUpdate(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).CFAPI()
	if err != nil {
		return err
	}
	id := d.Id()

	req := cfv2.SpaceUpdateRequest{}
	if d.HasChange("name") {
		req.Name = helpers.String(d.Get("name").(string))
	}

	_, err = cfClient.Spaces().Update(id, req)
	if err != nil {
		return fmt.Errorf("Error updating space: %s", err)
	}

	return resourceIBMSpaceRead(d, meta)
}

func resourceIBMSpaceDelete(d *schema.ResourceData, meta interface{}) error {
	cfClient, err := meta.(ClientSession).CFAPI()
	if err != nil {
		return err
	}
	id := d.Id()

	err = cfClient.Spaces().Delete(id)
	if err != nil {
		return fmt.Errorf("Error deleting space: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceIBMSpaceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cfClient, err := meta.(ClientSession).CFAPI()
	if err != nil {
		return false, err
	}
	id := d.Id()

	space, err := cfClient.Spaces().Get(id)
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
