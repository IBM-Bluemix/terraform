package ibmcloud

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"bluemix_api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Bluemix API Key",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"BM_API_KEY", "BLUEMIX_API_KEY"}, ""),
			},
			"bluemix_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout (in seconds) to set for any Bluemix API calls made.",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"BM_TIMEOUT", "BLUEMIX_TIMEOUT"}, 60),
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Bluemix Region (for example 'us-south').",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"BM_REGION", "BLUEMIX_REGION"}, "us-south"),
			},
			"softlayer_api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The SoftLayer API Key",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SL_API_KEY", "SOFTLAYER_API_KEY"}, ""),
			},
			"softlayer_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The SoftLayer user name",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SL_USERNAME", "SOFTLAYER_USERNAME"}, ""),
			},
			"softlayer_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout (in seconds) to set for any SoftLayer API calls made.",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SL_TIMEOUT", "SOFTLAYER_TIMEOUT"}, 60),
			},
			"skip_service_configuration": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "IBM Cloud has many services. At times you may not need to interact with some of them. This paramter allows to skip configuring clients for those",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"ibmcloud_cf_account":           dataSourceIBMCloudCfAccount(),
			"ibmcloud_cf_org":               dataSourceIBMCloudCfOrg(),
			"ibmcloud_cf_service_instance":  dataSourceIBMCloudCfServiceInstance(),
			"ibmcloud_cf_service_key":       dataSourceIBMCloudCfServiceKey(),
			"ibmcloud_cf_service_plan":      dataSourceIBMCloudCfServicePlan(),
			"ibmcloud_cf_space":             dataSourceIBMCloudCfSpace(),
			"ibmcloud_cs_cluster":           dataSourceIBMCloudCsCluster(),
			"ibmcloud_cs_cluster_config":    dataSourceIBMCloudArmadaClusterConfig(),
			"ibmcloud_cs_worker":            dataSourceIBMCloudCsWorker(),
			"ibmcloud_infra_dns_domain":     dataSourceIBMCloudInfraDNSDomain(),
			"ibmcloud_infra_image_template": dataSourceIBMCloudInfraImageTemplate(),
			"ibmcloud_infra_ssh_key":        dataSourceIBMCloudInfraSSHKey(),
			"ibmcloud_infra_virtual_guest":  dataSourceIBMCloudInfraVirtualGuest(),
			"ibmcloud_infra_vlan":           dataSourceIBMCloudInfraVlan(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"ibmcloud_cf_service_instance":               resourceIBMCloudCfServiceInstance(),
			"ibmcloud_cf_service_key":                    resourceIBMCloudCfServiceKey(),
			"ibmcloud_cf_space":                          resourceIBMCloudCfSpace(),
			"ibmcloud_cs_cluster":                        resourceIBMCloudArmadaCluster(),
			"ibmcloud_cs_cluster_service_bind":           resourceIBMCloudClusterBindService(),
			"ibmcloud_infra_bare_metal":                  resourceIBMCloudInfraBareMetal(),
			"ibmcloud_infra_basic_monitor":               resourceIBMCloudInfraBasicMonitor(),
			"ibmcloud_infra_block_storage":               resourceIBMCloudInfraBlockStorage(),
			"ibmcloud_infra_dns_domain":                  resourceIBMCloudInfraDNSDomain(),
			"ibmcloud_infra_dns_domain_record":           resourceIBMCloudInfraDNSDomainRecord(),
			"ibmcloud_infra_file_storage":                resourceIBMCloudInfraFileStorage(),
			"ibmcloud_infra_fw_hardware_dedicated":       resourceIBMCloudInfraFwHardwareDedicated(),
			"ibmcloud_infra_fw_hardware_dedicated_rules": resourceIBMCloudInfraFwHardwareDedicatedRules(),
			"ibmcloud_infra_global_ip":                   resourceIBMCloudInfraGlobalIp(),
			"ibmcloud_infra_lb_local":                    resourceIBMCloudInfraLbLocal(),
			"ibmcloud_infra_lb_local_service":            resourceIBMCloudInfraLbLocalService(),
			"ibmcloud_infra_lb_local_service_group":      resourceIBMCloudInfraLbLocalServiceGroup(),
			"ibmcloud_infra_lb_vpx":                      resourceIBMCloudInfraLbVpx(),
			"ibmcloud_infra_lb_vpx_ha":                   resourceIBMCloudInfraLbVpxHa(),
			"ibmcloud_infra_lb_vpx_service":              resourceIBMCloudInfraLbVpxService(),
			"ibmcloud_infra_lb_vpx_vip":                  resourceIBMCloudInfraLbVpxVip(),
			"ibmcloud_infra_objectstorage_account":       resourceIBMCloudInfraObjectStorageAccount(),
			"ibmcloud_infra_provisioning_hook":           resourceIBMCloudInfraProvisioningHook(),
			"ibmcloud_infra_scale_group":                 resourceIBMCloudInfraScaleGroup(),
			"ibmcloud_infra_scale_policy":                resourceIBMCloudInfraScalePolicy(),
			"ibmcloud_infra_security_certificate":        resourceIBMCloudInfraSecurityCertificate(),
			"ibmcloud_infra_ssh_key":                     resourceIBMCloudInfraSSHKey(),
			"ibmcloud_infra_user":                        resourceIBMCloudInfraUser(),
			"ibmcloud_infra_virtual_guest":               resourceIBMCloudInfraVirtualGuest(),
			"ibmcloud_infra_vlan":                        resourceIBMCloudInfraVlan(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	bluemixAPIKey := d.Get("bluemix_api_key").(string)
	softlayerUsername := d.Get("softlayer_username").(string)
	softlayerAPIKey := d.Get("softlayer_api_key").(string)
	softlayerTimeout := d.Get("softlayer_timeout").(int)
	bluemixTimeout := d.Get("bluemix_timeout").(int)
	region := d.Get("region").(string)

	skipServiceConfig := d.Get("skip_service_configuration").(*schema.Set)

	config := Config{
		BluemixAPIKey:        bluemixAPIKey,
		Region:               region,
		BluemixTimeout:       time.Duration(bluemixTimeout) * time.Second,
		SoftLayerTimeout:     time.Duration(softlayerTimeout) * time.Second,
		SoftLayerUserName:    softlayerUsername,
		SoftLayerAPIKey:      softlayerAPIKey,
		SkipServiceConfig:    skipServiceConfig,
		RetryCount:           3,
		RetryDelay:           30 * time.Millisecond,
		SoftLayerEndpointURL: SoftlayerRestEndpoint,
	}

	return config.ClientSession()
}
