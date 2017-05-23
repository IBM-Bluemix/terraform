package ibmcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceIBMCloudCsCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCloudCsClusterRead,

		Schema: map[string]*schema.Schema{
			"cluster_name_id": {
				Description: "Name or id of the cluster",
				Type:        schema.TypeString,
				Required:    true,
			},
			"worker_count": {
				Description: "Number of workers",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"workers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"org_guid": {
				Description: "The bluemix organization guid this cluster belongs to",
				Type:        schema.TypeString,
				Required:    true,
			},
			"space_guid": {
				Description: "The bluemix space guid this cluster belongs to",
				Type:        schema.TypeString,
				Required:    true,
			},
			"account_guid": {
				Description: "The bluemix account guid this cluster belongs to",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceIBMCloudCsClusterRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("cluster_name_id").(string)

	clusterClient := meta.(ClientSession).ClusterClient()
	workerClient := meta.(ClientSession).ClusterWorkerClient()

	targetEnv := getClusterTargetHeader(d)

	clusterFields, err := clusterClient.Find(name, targetEnv)
	if err != nil {
		return fmt.Errorf("Error retrieving cluster: %s", err)
	}
	workerFields, err := workerClient.List(name, targetEnv)
	if err != nil {
		return fmt.Errorf("Error retrieving workers for cluster: %s", err)
	}
	workers := make([]string, len(workerFields))
	for i, worker := range workerFields {
		workers[i] = worker.ID
	}
	d.SetId(clusterFields.ID)
	d.Set("worker_count", clusterFields.WorkerCount)
	d.Set("workers", workers)

	return nil
}