package ibmcloud

import (
	"fmt"
	"log"
	"time"

	v2 "github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
	"github.com/IBM-Bluemix/bluemix-go/helpers"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIBMCloudCfApp() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCloudCfAppCreate,
		Read:     resourceIBMCloudCfAppRead,
		Update:   resourceIBMCloudCfAppUpdate,
		Delete:   resourceIBMCloudCfAppDelete,
		Exists:   resourceIBMCloudCfAppExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name for the app",
			},
			"memory": {
				Description: "The amount of memory each instance should have. In megabytes.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"instances": {
				Description: "The number of instances",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"disk_quota": {
				Description: "The maximum amount of disk available to an instance of an app. In megabytes.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"space_guid": {
				Description: "Define space guid to which app belongs",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"command": {
				Description: "The initial command for the app",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"buildpack": {
				Description: "Buildpack to build the app. 3 options: a) Blank means autodetection; b) A Git Url pointing to a buildpack; c) Name of an installed buildpack.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"diego": {
				Description: "Use diego to stage and to run when available",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"environment_json": {
				Description: "Key/value pairs of all the environment variables to run in your app. Does not include any system or service variables.",
				Type:        schema.TypeMap,
				Optional:    true,
			},
			"ports": {
				Description: "Ports on which application may listen. Overwrites previously configured ports. Ports must be in range 1024-65535. Supported for Diego only.",
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Set: func(v interface{}) int {
					return v.(int)
				},
			},
			"route_guid": {
				Description: "Define the route guids which should be bound to the application.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},
			"service_instance_guid": {
				Description: "Define the service instance guids that should be bound to this application.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},
			"wait_time_minutes": {
				Description: "Define timeout to wait for the app to start. Default is no wait",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
			},
			"app_path": {
				Description: "Define the  path of the zip file of the application.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"app_version": {
				Description: "Version of the application",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceIBMCloudCfAppCreate(d *schema.ResourceData, meta interface{}) error {
	appClient := meta.(ClientSession).CloudFoundryAppClient()

	name := d.Get("name").(string)
	spaceGUID := d.Get("space_guid").(string)

	appCreatePayload := v2.AppRequest{
		Name:      helpers.String(name),
		SpaceGUID: helpers.String(spaceGUID),
	}

	if memory, ok := d.GetOk("memory"); ok {
		appCreatePayload.Memory = memory.(int)
	}

	if instances, ok := d.GetOk("instances"); ok {
		appCreatePayload.Instances = instances.(int)
	}

	if diskQuota, ok := d.GetOk("disk_quota"); ok {
		appCreatePayload.DiskQuota = diskQuota.(int)
	}

	if command, ok := d.GetOk("command"); ok {
		appCreatePayload.Command = helpers.String(command.(string))
	}

	if buildpack, ok := d.GetOk("buildpack"); ok {
		appCreatePayload.BuildPack = helpers.String(buildpack.(string))
	}

	appCreatePayload.Diego = d.Get("diego").(bool)

	if portSet := d.Get("ports").(*schema.Set); len(portSet.List()) > 0 {
		ports := expandIntList(portSet.List())
		appCreatePayload.Ports = helpers.IntSlice(ports)
	}
	if environmentJSON, ok := d.GetOk("environment_json"); ok {
		appCreatePayload.EnvironmentJSON = helpers.Map(environmentJSON.(map[string]interface{}))

	}
	_, err := appClient.FindByName(spaceGUID, name)
	if err == nil {
		return fmt.Errorf("%s already exists in the given space %s", name, spaceGUID)
	}

	log.Println("[INFO] Creating Cloud Foundary Application")

	app, err := appClient.Create(&appCreatePayload)

	if err != nil {
		return fmt.Errorf("Error creating app: %s", err)
	}

	log.Println("[INFO] Cloud Foundary Application is created successfully")

	d.SetId(app.Metadata.GUID)

	log.Println("[INFO] Bind the route with cloud foundary application")

	if v, ok := d.Get("route_guid").(*schema.Set); ok && v.Len() > 0 {
		for _, routeID := range v.List() {
			_, err := appClient.BindRoute(app.Metadata.GUID, routeID.(string))
			if err != nil {
				return fmt.Errorf("Error binding route %s to app: %s", routeID.(string), err)
			}
		}
	}
	if v, ok := d.Get("service_instance_guid").(*schema.Set); ok && v.Len() > 0 {
		sbClient := meta.(ClientSession).CloudFoundryServiceBindingClient()
		for _, svcID := range v.List() {
			req := v2.ServiceBindingRequest{
				ServiceInstanceGUID: svcID.(string),
				AppGUID:             app.Metadata.GUID,
			}
			_, err := sbClient.Create(req)
			if err != nil {
				return fmt.Errorf("Error binding service %s to  app: %s", svcID.(string), err)
			}
		}
	}

	log.Println("[INFO] Upload the app bits to the cloud foundary application")

	if appPath, ok := d.GetOk("app_path"); ok {
		_, err = appClient.Upload(app.Metadata.GUID, appPath.(string))
		if err != nil {
			return fmt.Errorf("Error uploading app bits: %s", err)
		}
	}

	log.Println("[INFO] Start Cloud Foundary Application")

	waitTimeout := time.Duration(d.Get("wait_time_minutes").(int)) * time.Minute

	status, err := appClient.Start(app.Metadata.GUID, waitTimeout)
	if err != nil {
		return fmt.Errorf("Error while starting  app: %s", err)
	}
	//If you are explcity told to wait till the application has started
	if waitTimeout != 0 {
		if status.PackageState != v2.AppStagedState {
			return fmt.Errorf("Applications couldn't be staged  %s", err)
		}
		if status.InstanceState != v2.AppRunningState {
			return fmt.Errorf("Applications instances  %s", err)
		}
	}

	log.Printf("[INFO]Cloud Foundary Application: %s has started successfully", name)

	return resourceIBMCloudCfAppRead(d, meta)
}

func resourceIBMCloudCfAppRead(d *schema.ResourceData, meta interface{}) error {
	appClient := meta.(ClientSession).CloudFoundryAppClient()
	appGUID := d.Id()

	appData, err := appClient.Get(appGUID)
	if err != nil {
		return fmt.Errorf("Error retrieving app: %s", err)
	}

	d.SetId(appData.Metadata.GUID)
	d.Set("name", appData.Entity.Name)
	d.Set("memory", appData.Entity.Memory)
	d.Set("instances", appData.Entity.Instances)
	d.Set("space_guid", appData.Entity.SpaceGUID)
	d.Set("disk_quota", appData.Entity.DiskQuota)
	d.Set("ports", flattenPort(appData.Entity.Ports))
	d.Set("command", appData.Entity.Command)
	d.Set("buildpack", appData.Entity.BuildPack)
	d.Set("diego", appData.Entity.Diego)
	d.Set("environment_json", appData.Entity.EnvironmentJSON)

	route, err := appClient.ListRoutes(appGUID)
	if err != nil {
		return err
	}
	if len(route) > 0 {
		d.Set("route_guid", flattenRoute(route))
	}

	svcBindings, err := appClient.ListServiceBindings(appGUID)
	if err != nil {
		return err
	}
	if len(route) > 0 {
		d.Set("service_instance_guid", flattenServiceBindings(svcBindings))
	}

	return nil

}

func resourceIBMCloudCfAppUpdate(d *schema.ResourceData, meta interface{}) error {
	appClient := meta.(ClientSession).CloudFoundryAppClient()
	appGUID := d.Id()

	appUpdatePayload := v2.AppRequest{}

	if d.HasChange("name") {
		appUpdatePayload.Name = helpers.String(d.Get("name").(string))
	}

	if d.HasChange("memory") {
		appUpdatePayload.Memory = d.Get("memory").(int)
	}

	if d.HasChange("instances") {
		appUpdatePayload.Instances = d.Get("instances").(int)
	}

	if d.HasChange("disk_quota") {
		appUpdatePayload.DiskQuota = d.Get("disk_quota").(int)
	}

	if d.HasChange("buildpack") {
		appUpdatePayload.BuildPack = helpers.String(d.Get("buildpack").(string))
	}

	if d.HasChange("command") {
		appUpdatePayload.Command = helpers.String(d.Get("command").(string))
	}

	if d.HasChange("diego") {
		appUpdatePayload.Diego = d.Get("diego").(bool)
	}

	if d.HasChange("environment_json") {
		appUpdatePayload.EnvironmentJSON = helpers.Map(d.Get("environment_json").(map[string]interface{}))
	}

	if d.HasChange("ports") {
		portSet := d.Get("ports").(*schema.Set)
		if portSet.Len() == 0 {
			return fmt.Errorf("ports field can't be updated to have 0 elements")
		}
		ports := expandIntList(portSet.List())

		appUpdatePayload.Ports = helpers.IntSlice(ports)
	}

	log.Println("[INFO] Update cloud foundary application")

	_, err := appClient.Update(appGUID, &appUpdatePayload)
	if err != nil {
		return fmt.Errorf("Error updating application: %s", err)
	}

	if d.HasChange("app_path") {
		appPath := d.Get("app_path").(string)
		_, err = appClient.Upload(appGUID, appPath)
		if err != nil {
			return fmt.Errorf("Error uploading  app: %s", err)
		}
		appUpdatePayload := &v2.AppRequest{
			State: helpers.String(v2.AppStoppedState),
		}
		_, err := appClient.Update(appGUID, appUpdatePayload)
		if err != nil {
			return fmt.Errorf("Error updating application: %s", err)
		}
		waitTimeout := time.Duration(d.Get("wait_time_minutes").(int)) * time.Minute
		_, err = appClient.Start(appGUID, waitTimeout)
		if err != nil {
			return fmt.Errorf("Error while starting  app: %s", err)
		}

	}

	if d.HasChange("route_guid") {
		ors, nrs := d.GetChange("route_guid")
		or := ors.(*schema.Set)
		nr := nrs.(*schema.Set)

		remove := expandStringList(or.Difference(nr).List())
		add := expandStringList(nr.Difference(or).List())

		if len(add) > 0 {
			for i := range add {
				_, err = appClient.BindRoute(appGUID, add[i])
				if err != nil {
					return fmt.Errorf("Error while binding route %q to application %s: %q", add[i], appGUID, err)
				}
			}
		}
		if len(remove) > 0 {
			for i := range remove {
				err = appClient.UnBindRoute(appGUID, remove[i])
				if err != nil {
					return fmt.Errorf("Error while un-binding route %q from application %s: %q", add[i], appGUID, err)
				}
			}
		}

	}

	if d.HasChange("service_instance_guid") {
		oss, nss := d.GetChange("service_instance_guid")
		os := oss.(*schema.Set)
		ns := nss.(*schema.Set)
		remove := expandStringList(os.Difference(ns).List())
		add := expandStringList(ns.Difference(os).List())

		if len(add) > 0 {
			sbClient := meta.(ClientSession).CloudFoundryServiceBindingClient()
			for i := range add {
				sbPayload := v2.ServiceBindingRequest{
					ServiceInstanceGUID: add[i],
					AppGUID:             appGUID,
				}
				_, err = sbClient.Create(sbPayload)
				if err != nil {
					return fmt.Errorf("Error while binding service instance %s to application %s: %q", add[i], appGUID, err)
				}

			}
		}
		if len(remove) > 0 {
			for i := range remove {
				err = appClient.DeleteServiceBinding(appGUID, remove[i])
				if err != nil {
					return fmt.Errorf("Error while binding route %s to application %s: %q", remove[i], appGUID, err)
				}
			}
		}

	}

	return resourceIBMCloudCfAppRead(d, meta)
}

func resourceIBMCloudCfAppDelete(d *schema.ResourceData, meta interface{}) error {
	appClient := meta.(ClientSession).CloudFoundryAppClient()
	id := d.Id()

	err := appClient.Delete(id)
	if err != nil {
		return fmt.Errorf("Error deleting app: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceIBMCloudCfAppExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	appClient := meta.(ClientSession).CloudFoundryAppClient()
	id := d.Id()

	app, err := appClient.Get(id)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return app.Metadata.GUID == id, nil
}
