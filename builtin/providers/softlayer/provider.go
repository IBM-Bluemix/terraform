package softlayer

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOFTLAYER_USERNAME", nil),
				Description: "The user name for SoftLayer API operations.",
			},
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SOFTLAYER_API_KEY", nil),
				Description: "The API key for SoftLayer API operations.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"softlayer_virtual_guest":         resourceSoftLayerVirtualGuest(),
			"softlayer_ssh_key":               resourceSoftLayerSSHKey(),
			"softlayer_dns_domain_record":     resourceSoftLayerDnsDomainResourceRecord(),
			"softlayer_dns_domain":            resourceSoftLayerDnsDomain(),
			"softlayer_lb_vpx":                resourceSoftLayerNetworkApplicationDeliveryController(),
			"softlayer_lb_vpx_vip":            resourceSoftLayerNetworkLoadBalancerVirtualIpAddress(),
			"softlayer_lb_vpx_service":        resourceSoftLayerNetworkLoadBalancerService(),
			"softlayer_security_certificate":  resourceSoftLayerSecurityCertificate(),
			"softlayer_user":                  resourceSoftLayerUserCustomer(),
			"softlayer_objectstorage_account": resourceSoftLayerObjectStorageAccount(),
			"softlayer_provisioning_hook":     resourceSoftLayerProvisioningHook(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Username: d.Get("username").(string),
		ApiKey:   d.Get("api_key").(string),
	}

	return config.Client()
}
