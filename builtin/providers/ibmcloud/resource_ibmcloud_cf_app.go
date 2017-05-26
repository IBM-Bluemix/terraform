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

const (
	appStopState = "STOPPED"
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
				Description:  "The amount of memory each instance should have. In megabytes.",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      512,
				ValidateFunc: validateAppQuota,
			},
			"instances": {
				Description:  "The number of instances",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validateAppInstance,
			},
			"disk_quota": {
				Description:  "The maximum amount of disk available to an instance of an app. In megabytes.",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1024,
				ValidateFunc: validateAppQuota,
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
				Description:   "Buildpack to build the app. 3 options: a) Blank means autodetection; b) A Git Url pointing to a buildpack; c) Name of an installed buildpack.",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"docker_image"},
			},
			"diego": {
				Description: "Use diego to stage and to run when available",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"docker_image": {
				Description:   "Name of the Docker image containing the app",
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"buildpack"},
			},
			"docker_credentials_json": {
				Description: "Docker credentials for pulling docker image.",
				Type:        schema.TypeMap,
				Optional:    true,
			},
			"environment_json": {
				Description: "Key/value pairs of all the environment variables to run in your app. Does not include any system or service variables.t",
				Type:        schema.TypeMap,
				Optional:    true,
			},
			"ports": {
				Description: "Ports on which application may listen. Overwrites previously configured ports. Ports must be in range 1024-65535. Supported for Diego only.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"route_guid": {
				Description: "Define the route guid needs to be attached to application.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},
			"service_instance_guid": {
				Description: "Define the service guid needs to be attached to application.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},
			"wait_timeout": {
				Description: "Define timeout to wait for the app to start",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
			},
			"app_path": {
				Description: "Define the path of the zip file of the application.",
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

	if diego, ok := d.GetOk("diego"); ok {
		appCreatePayload.Diego = diego.(bool)
	}

	if dockerImage, ok := d.GetOk("docker_image"); ok {
		appCreatePayload.DockerImage = helpers.String(dockerImage.(string))
	}

	if port, ok := d.GetOk("ports"); ok {
		p := port.([]interface{})
		ports := make([]int, len(p))
		for i := range p {
			ports[i] = p[i].(int)
		}
		appCreatePayload.Ports = helpers.Array(ports)
	}

	if dockerCredentialsJSON, ok := d.GetOk("docker_credentials_json"); ok {
		appCreatePayload.DockerCredentialsJSON = helpers.Map(dockerCredentialsJSON.(map[string]interface{}))
	}
	if environmentJSON, ok := d.GetOk("environment_json"); ok {
		appCreatePayload.EnvironmentJSON = helpers.Map(environmentJSON.(map[string]interface{}))

	}

	_, err := appClient.FindByName(spaceGUID, name)

	if err == nil {
		return fmt.Errorf("%s already exists", name)
	}

	log.Println("[INFO] Creating Cloud Foundary Application")

	app, err := appClient.Create(&appCreatePayload)

	if err != nil {
		return fmt.Errorf("Error creating app: %s", err)
	}

	log.Println("[INFO] Cloud Foundary Application is created successfully")

	d.SetId(app.Metadata.GUID)

	log.Println("[INFO] Bind the route with cloud foundary application")

	routeIDs := d.Get("route_guid").(*schema.Set)
	for _, routeID := range routeIDs.List() {
		if routeID != "" {
			_, err := appClient.BindRoute(app.Metadata.GUID, routeID.(string))
			if err != nil {
				return fmt.Errorf("Error binding route %s to  app: %s", routeID.(string), err)
			}
		}
	}

	log.Println("[INFO] Upload the app bits to the cloud foundary application")

	appPath := d.Get("app_path").(string)
	if appPath != "" {
		_, err = appClient.Upload(app.Metadata.GUID, appPath)
		if err != nil {
			return fmt.Errorf("Error uploading  app: %s", err)
		}
	}

	log.Println("[INFO] Start Cloud Foundary Application")

	waitTimeout := time.Duration(d.Get("wait_timeout").(int)) * time.Minute
	fmt.Println("time out", waitTimeout)
	_, err = appClient.Start(app.Metadata.GUID, waitTimeout)
	if err != nil {
		return fmt.Errorf("Error while starting  app: %s", err)
	}

	log.Println("[INFO]Cloud Foundary Application has started successfully")

	log.Println("[INFO] Bind the service instance with cloud foundary application")

	sbClient := meta.(ClientSession).CloudFoundryServiceBindingClient()

	serviceIDs := d.Get("service_instance_guid").(*schema.Set)
	for _, serviceID := range serviceIDs.List() {
		if serviceID != "" {
			sbPayload := v2.ServiceBindingRequest{
				ServiceInstanceGUID: serviceID.(string),
				AppGUID:             app.Metadata.GUID,
			}
			_, err := sbClient.Create(sbPayload)
			if err != nil {
				return fmt.Errorf("Error binding service %s to  app: %s", serviceID.(string), err)
			}
		}
	}

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
	d.Set("ports", appData.Entity.Ports)
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
		p := d.Get("ports").([]interface{})
		ports := make([]int, len(p))
		for i := range p {
			ports[i] = p[i].(int)
		}
		appUpdatePayload.Ports = helpers.Array(ports)
	}

	log.Println("[INFO] Update cloud foundary application")

	_, err := appClient.Update(appGUID, &appUpdatePayload)
	if err != nil {
		return fmt.Errorf("Error updating space: %s", err)
	}

	var appPath string
	if d.HasChange("app_path") {
		appPath = d.Get("app_path").(string)
		_, err = appClient.Upload(appGUID, appPath)
		if err != nil {
			return fmt.Errorf("Error uploading  app: %s", err)
		}
		appUpdatePayload := &v2.AppRequest{
			State: helpers.String(appStopState),
		}
		_, err := appClient.Update(appGUID, appUpdatePayload)
		if err != nil {
			return fmt.Errorf("Error updating space: %s", err)
		}
		waitTimeout := time.Duration(d.Get("wait_timeout").(int)) * time.Minute
		_, err = appClient.Start(appGUID, waitTimeout)
		if err != nil {
			return fmt.Errorf("Error while starting  app: %s", err)
		}

	}

	if d.HasChange("route_guid") {
		oldroutes, newroutes := d.GetChange("route_guid")
		oldRoute := oldroutes.(*schema.Set)
		newRoute := newroutes.(*schema.Set)

		remove := expandStringList(oldRoute.Difference(newRoute).List())
		add := expandStringList(newRoute.Difference(oldRoute).List())

		if len(add) > 0 {
			for i := range add {
				_, err = appClient.BindRoute(appGUID, add[i])
				if err != nil {
					return fmt.Errorf("Error while binding route : %s", err)
				}
			}
		}
		if len(remove) > 0 {
			for i := range remove {
				err = appClient.UnBindRoute(appGUID, remove[i])
				if err != nil {
					return fmt.Errorf("Error while unbinding route: %s", err)
				}
			}
		}

	}

	if d.HasChange("service_instance_guid") {
		oldServices, newServices := d.GetChange("service_instance_guid")
		oldService := oldServices.(*schema.Set)
		newService := newServices.(*schema.Set)
		remove := expandStringList(oldService.Difference(newService).List())
		add := expandStringList(newService.Difference(oldService).List())

		if len(add) > 0 {
			for i := range add {
				sbClient := meta.(ClientSession).CloudFoundryServiceBindingClient()
				sbPayload := v2.ServiceBindingRequest{
					ServiceInstanceGUID: add[i],
					AppGUID:             appGUID,
				}
				_, err = sbClient.Create(sbPayload)
				if err != nil {
					return fmt.Errorf("Error while binding service: %s", err)
				}

			}
		}
		if len(remove) > 0 {
			for i := range remove {
				err = appClient.DeleteServiceBinding(appGUID, remove[i])
				if err != nil {
					return fmt.Errorf("Error while unbinding  service : %s", err)
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
