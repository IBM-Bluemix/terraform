package ibmcloud

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func resourceIBMCloudInfraProvisioningHook() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCloudInfraProvisioningHookCreate,
		Read:     resourceIBMCloudInfraProvisioningHookRead,
		Update:   resourceIBMCloudInfraProvisioningHookUpdate,
		Delete:   resourceIBMCloudInfraProvisioningHookDelete,
		Exists:   resourceIBMCloudInfraProvisioningHookExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"uri": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceIBMCloudInfraProvisioningHookCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetProvisioningHookService(sess)

	opts := datatypes.Provisioning_Hook{
		Name: sl.String(d.Get("name").(string)),
		Uri:  sl.String(d.Get("uri").(string)),
	}

	hook, err := service.CreateObject(&opts)
	if err != nil {
		return fmt.Errorf("Error creating Provisioning Hook: %s", err)
	}

	d.SetId(strconv.Itoa(*hook.Id))
	log.Printf("[INFO] Provisioning Hook ID: %d", *hook.Id)

	return resourceIBMCloudInfraProvisioningHookRead(d, meta)
}

func resourceIBMCloudInfraProvisioningHookRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetProvisioningHookService(sess)

	hookId, _ := strconv.Atoi(d.Id())

	hook, err := service.Id(hookId).GetObject()
	if err != nil {
		if err, ok := err.(sl.Error); ok {
			if err.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("Error retrieving Provisioning Hook: %s", err)
	}

	d.Set("id", hook.Id)
	d.Set("name", hook.Name)
	d.Set("uri", hook.Uri)

	return nil
}

func resourceIBMCloudInfraProvisioningHookUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetProvisioningHookService(sess)

	hookId, _ := strconv.Atoi(d.Id())

	opts := datatypes.Provisioning_Hook{}

	if d.HasChange("name") {
		opts.Name = sl.String(d.Get("name").(string))
	}

	if d.HasChange("uri") {
		opts.Uri = sl.String(d.Get("uri").(string))
	}

	opts.TypeId = sl.Int(1)
	_, err := service.Id(hookId).EditObject(&opts)

	if err != nil {
		return fmt.Errorf("Error editing Provisioning Hook: %s", err)
	}
	return nil
}

func resourceIBMCloudInfraProvisioningHookDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetProvisioningHookService(sess)

	hookId, err := strconv.Atoi(d.Id())
	log.Printf("[INFO] Deleting Provisioning Hook: %d", hookId)
	_, err = service.Id(hookId).DeleteObject()
	if err != nil {
		return fmt.Errorf("Error deleting Provisioning Hook: %s", err)
	}

	return nil
}

func resourceIBMCloudInfraProvisioningHookExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetProvisioningHookService(sess)

	hookId, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(hookId).GetObject()
	return result.Id != nil && err == nil && *result.Id == hookId, nil
}
