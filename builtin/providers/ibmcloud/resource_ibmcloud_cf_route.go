package ibmcloud

import (
	"fmt"

	v2 "github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/IBM-Bluemix/bluemix-go/helpers"

	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIBMCloudCfRoute() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCloudCfRouteCreate,
		Read:     resourceIBMCloudCfRouteRead,
		Update:   resourceIBMCloudCfRouteUpdate,
		Delete:   resourceIBMCloudCfRouteDelete,
		Exists:   resourceIBMCloudCfRouteExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The host portion of the route. Required for shared-domains.",
			},

			"space_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The guid of the associated space",
			},

			"domain_guid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The guid of the associated domain",
			},

			"port": {
				Description:  "The port of the route. Supported for domains of TCP router groups only.",
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validateRoutePort,
			},

			"path": {
				Description:  "The path for a route as raw text.Paths must be between 2 and 128 characters.Paths must start with a forward slash '/'.Paths must not contain a '?'",
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validateRoutePath,
			},
		},
	}
}

func resourceIBMCloudCfRouteCreate(d *schema.ResourceData, meta interface{}) error {
	routeClient := meta.(ClientSession).CloudFoundryRouteClient()

	spaceGUID := d.Get("space_guid").(string)
	domainGUID := d.Get("domain_guid").(string)

	params := v2.RouteRequest{
		SpaceGUID:  spaceGUID,
		DomainGUID: domainGUID,
	}

	if host, ok := d.GetOk("host"); ok {
		params.Host = host.(string)
	}

	if port, ok := d.GetOk("port"); ok {
		params.Port = port.(int)
	}

	if path, ok := d.GetOk("path"); ok {
		params.Path = path.(string)
	}

	route, err := routeClient.Create(params)
	if err != nil {
		return fmt.Errorf("Error creating route: %s", err)
	}

	d.SetId(route.Metadata.GUID)

	return resourceIBMCloudCfRouteRead(d, meta)
}

func resourceIBMCloudCfRouteRead(d *schema.ResourceData, meta interface{}) error {
	routeClient := meta.(ClientSession).CloudFoundryRouteClient()
	routeGUID := d.Id()

	route, err := routeClient.Get(routeGUID)
	if err != nil {
		return fmt.Errorf("Error retrieving route: %s", err)
	}

	d.Set("host", route.Entity.Host)
	d.Set("space_guid", route.Entity.SpaceGUID)
	d.Set("domain_guid", route.Entity.DomainGUID)
	d.Set("port", route.Entity.Port)
	d.Set("path", route.Entity.Path)

	return nil
}

func resourceIBMCloudCfRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	routeClient := meta.(ClientSession).CloudFoundryRouteClient()

	routeGUID := d.Id()

	params := v2.RouteUpdateRequest{}

	if d.HasChange("host") {
		params.Host = helpers.String(d.Get("host").(string))
	}

	if d.HasChange("port") {
		params.Port = helpers.Int(d.Get("port").(int))
	}

	if d.HasChange("path") {
		params.Path = helpers.String(d.Get("path").(string))
	}

	_, err := routeClient.Update(routeGUID, params)
	if err != nil {
		return fmt.Errorf("Error updating route: %s", err)
	}

	return resourceIBMCloudCfRouteRead(d, meta)
}

func resourceIBMCloudCfRouteDelete(d *schema.ResourceData, meta interface{}) error {
	routeClient := meta.(ClientSession).CloudFoundryRouteClient()
	routeGUID := d.Id()

	err := routeClient.Delete(routeGUID, true)
	if err != nil {
		return fmt.Errorf("Error deleting route: %s", err)
	}

	d.SetId("")

	return nil
}
func resourceIBMCloudCfRouteExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	routeClient := meta.(ClientSession).CloudFoundryRouteClient()
	routeGUID := d.Id()

	route, err := routeClient.Get(routeGUID)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return route.Metadata.GUID == routeGUID, nil
}
