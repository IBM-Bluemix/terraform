package ibmcloud

import (
	"fmt"

	v1 "github.com/IBM-Bluemix/bluemix-go/api/k8scluster/k8sclusterv1"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIBMCloudClusterBindService() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCloudClusterBindServiceCreate,
		Read:     resourceIBMCloudClusterBindServiceRead,
		Delete:   resourceIBMCloudClusterBindServiceDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"cluster_name_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"service_instance_space_guid": {
				Type:        schema.TypeString,
				Description: "The space guid the service instance belongs to",
				ForceNew:    true,
				Required:    true,
			},
			"service_instance_name_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"secret_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_guid": {
				Description: "The bluemix organization guid this cluster belongs to",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"space_guid": {
				Description: "The bluemix space guid this cluster belongs to",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"account_guid": {
				Description: "The bluemix account guid this cluster belongs to",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func getClusterTargetHeader(d *schema.ResourceData) *v1.ClusterTargetHeader {
	orgGUID := d.Get("org_guid").(string)
	spaceGUID := d.Get("space_guid").(string)
	accountGUID := d.Get("account_guid").(string)

	targetEnv := &v1.ClusterTargetHeader{
		OrgID:     orgGUID,
		SpaceID:   spaceGUID,
		AccountID: accountGUID,
	}
	return targetEnv
}

func resourceIBMCloudClusterBindServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ClientSession).ClusterClient()
	clusterNameID := d.Get("cluster_name_id").(string)
	serviceInstanceSpaceGUID := d.Get("service_instance_space_guid").(string)
	serviceInstanceNameID := d.Get("service_instance_name_id").(string)
	namespaceID := d.Get("namespace_id").(string)

	bindService := &v1.ServiceBindRequest{
		ClusterNameOrID:         clusterNameID,
		SpaceGUID:               serviceInstanceSpaceGUID,
		ServiceInstanceNameOrID: serviceInstanceNameID,
		NamespaceID:             namespaceID,
	}

	targetEnv := getClusterTargetHeader(d)
	bindResp, err := client.BindService(bindService, targetEnv)
	if err != nil {
		return err
	}
	d.SetId(clusterNameID)
	d.Set("service_instance_name_id", serviceInstanceNameID)
	d.Set("namespace_id", namespaceID)
	d.Set("space_guid", serviceInstanceSpaceGUID)
	d.Set("secret_name", bindResp.SecretName)

	return resourceIBMCloudClusterBindServiceRead(d, meta)
}

func resourceIBMCloudClusterBindServiceRead(d *schema.ResourceData, meta interface{}) error {
	//No API to read back the credentials so leave schema as it is
	return nil
}

func resourceIBMCloudClusterBindServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(ClientSession).ClusterClient()
	clusterID := d.Id()
	namespace := d.Get("namespace_id").(string)
	serviceInstanceNameId := d.Get("service_instance_name_id").(string)
	targetEnv := getClusterTargetHeader(d)

	err := client.UnBindService(clusterID, namespace, serviceInstanceNameId, targetEnv)
	if err != nil {
		return fmt.Errorf("Error unbinding service: %s", err)
	}
	return nil
}
