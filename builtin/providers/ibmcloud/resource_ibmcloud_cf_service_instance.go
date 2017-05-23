package ibmcloud

import (
	"fmt"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIBMCloudCfServiceInstance() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCloudCfServiceInstanceCreate,
		Read:     resourceIBMCloudCfServiceInstanceRead,
		Update:   resourceIBMCloudCfServiceInstanceUpdate,
		Delete:   resourceIBMCloudCfServiceInstanceDelete,
		Exists:   resourceIBMCloudCfServiceInstanceExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A name for the service instance",
			},

			"space_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The guid of the space in which the instance will be created",
			},

			"service": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the service",
			},

			"credentials": {
				Description: "Credentials asociated with the key",
				Computed:    true,
				Type:        schema.TypeMap,
			},

			"service_plan_guid": {
				Description: "The uniquie identifier of the service offering plan type",
				Computed:    true,
				Type:        schema.TypeString,
			},

			"parameters": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Arbitrary parameters to pass along to the service broker. Must be a JSON object",
			},

			"plan": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The plan type of the service",
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceIBMCloudCfServiceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	serviceName := d.Get("service").(string)
	plan := d.Get("plan").(string)

	srOff := meta.(ClientSession).CloudFoundryServiceOfferingClient()
	serviceOff, err := srOff.FindByLabel(serviceName)
	if err != nil {
		return fmt.Errorf("Error retrieving service offering: %s", err)
	}

	srPlan := meta.(ClientSession).CloudFoundryServicePlanClient()
	servicePlan, err := srPlan.GetServicePlan(serviceOff.GUID, plan)
	if err != nil {
		return fmt.Errorf("Error retrieving plan: %s", err)
	}

	serviceInst := meta.(ClientSession).CloudFoundryServiceInstanceClient()

	name := d.Get("name").(string)
	spaceGUID := d.Get("space_guid").(string)
	var parameters map[string]interface{}
	var tags []string

	if parameters, ok := d.GetOk("parameters"); ok {
		parameters = parameters.(map[string]interface{})
	}

	if _, ok := d.GetOk("tags"); ok {
		tags = getServiceTags(d)
	}

	service, err := serviceInst.Create(name, servicePlan.GUID, spaceGUID, parameters, tags)
	if err != nil {
		return fmt.Errorf("Error creating service: %s", err)
	}

	d.SetId(service.Metadata.GUID)

	return resourceIBMCloudCfServiceInstanceRead(d, meta)
}

func resourceIBMCloudCfServiceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	serviceClient := meta.(ClientSession).CloudFoundryServiceInstanceClient()
	serviceGUID := d.Id()

	service, err := serviceClient.Get(serviceGUID)
	if err != nil {
		return fmt.Errorf("Error retrieving service: %s", err)
	}

	d.Set("service_plan_guid", service.Entity.ServicePlanGUID)
	d.Set("credentials", service.Entity.Credentials)

	return nil
}

func resourceIBMCloudCfServiceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	serviceClient := meta.(ClientSession).CloudFoundryServiceInstanceClient()

	serviceGUID := d.Id()

	var name, planguid string
	var parameters map[string]interface{}
	var tags []string

	name = d.Get("name").(string)

	if d.HasChange("plan") {
		plan := d.Get("plan").(string)
		service := d.Get("service").(string)
		srOff := meta.(ClientSession).CloudFoundryServiceOfferingClient()
		serviceOff, err := srOff.FindByLabel(service)
		if err != nil {
			return fmt.Errorf("Error retrieving service offering: %s", err)
		}

		srPlan := meta.(ClientSession).CloudFoundryServicePlanClient()
		servicePlan, err := srPlan.GetServicePlan(serviceOff.GUID, plan)
		if err != nil {
			return fmt.Errorf("Error retrieving plan: %s", err)
		}
		planguid = servicePlan.GUID

	}

	if d.HasChange("parameters") {
		parameters = d.Get("parameters").(map[string]interface{})
	}

	if d.HasChange("tags") {
		tags = getServiceTags(d)
	}

	_, err := serviceClient.Update(name, serviceGUID, planguid, parameters, tags)
	if err != nil {
		return fmt.Errorf("Error updating service: %s", err)
	}

	return resourceIBMCloudCfServiceInstanceRead(d, meta)
}

func resourceIBMCloudCfServiceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	serviceClient := meta.(ClientSession).CloudFoundryServiceInstanceClient()

	id := d.Id()

	err := serviceClient.Delete(id)
	if err != nil {
		return fmt.Errorf("Error deleting service: %s", err)
	}

	d.SetId("")

	return nil
}
func resourceIBMCloudCfServiceInstanceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	serviceClient := meta.(ClientSession).CloudFoundryServiceInstanceClient()
	serviceGUID := d.Id()

	service, err := serviceClient.Get(serviceGUID)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return service.Metadata.GUID == serviceGUID, nil
}

func getServiceTags(d *schema.ResourceData) []string {
	tagSet := d.Get("tags").(*schema.Set)

	if tagSet.Len() == 0 {
		empty := []string{}
		return empty
	}

	tags := make([]string, 0, tagSet.Len())
	for _, elem := range tagSet.List() {
		tag := elem.(string)
		tags = append(tags, tag)
	}
	return tags
}
