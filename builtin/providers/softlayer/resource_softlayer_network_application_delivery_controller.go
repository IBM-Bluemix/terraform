package softlayer

import (
	"fmt"
	"log"
	"strconv"

//	"github.com/TheWeatherCompany/softlayer-go/services"
	softlayer "github.com/TheWeatherCompany/softlayer-go/softlayer"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

const (
	NETSCALER_VPX_TYPE = "Netscaler VPX"
)

func resourceSoftLayerNetworkApplicationDeliveryController() *schema.Resource {
	return &schema.Resource{
		Create: resourceSoftLayerNetworkApplicationDeliveryControllerCreate,
		Read:   resourceSoftLayerNetworkApplicationDeliveryControllerRead,
		Update: resourceSoftLayerNetworkApplicationDeliveryControllerUpdate,
		Delete: resourceSoftLayerNetworkApplicationDeliveryControllerDelete,
		Exists: resourceSoftLayerNetworkApplicationDeliveryControllerExists,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"datacenter": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"speed": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"plan": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ip_count": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceSoftLayerNetworkApplicationDeliveryControllerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).networkApplicationDeliveryControllerService
	if client == nil {
		return fmt.Errorf("The client is nil.")
	}

	nadcType := d.Get("type").(string)

	switch nadcType {
	default:
		return fmt.Errorf("[ERROR] Network application delivery controller type %s is not supported", nadcType)
	case NETSCALER_VPX_TYPE:
		// create Netscaler VPX
		opts := softlayer.NetworkApplicationDeliveryControllerCreateOptions{
			Speed:    d.Get("speed").(int),
			Version:  d.Get("version").(string),
			Plan:     d.Get("plan").(string),
			IpCount:  d.Get("ip_count").(int),
			Location: d.Get("datacenter").(string),
		}

		log.Printf("[INFO] Creating network application delivery controller")

		netscalerVPX, err := client.CreateNetscalerVPX(&opts)

		if err != nil {
			return fmt.Errorf("Error creating network application delivery controller: %s", err)
		}

		d.SetId(fmt.Sprintf("%d", netscalerVPX.Id))

		log.Printf("[INFO] Netscaler VPX ID: %s", d.Id())
	}

	return resourceSoftLayerNetworkApplicationDeliveryControllerRead(d, meta)
}

func resourceSoftLayerNetworkApplicationDeliveryControllerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).networkApplicationDeliveryControllerService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}
	getObjectResult, err := client.GetObject(id)
	if err != nil {
		return fmt.Errorf("Error retrieving network application delivery controller: %s", err)
	}

	d.Set("name", getObjectResult.Name)
	d.Set("type", getObjectResult.Type)
	if getObjectResult.Datacenter != nil {
		d.Set("location", getObjectResult.Datacenter.Name)
	}

	return nil
}

// these parameters are not contained in the output object, so should be parsed from the
// description string
// example string to be parsed CITRIX_NETSCALER_VPX_10_1_10MBPS_STANDARD
// 10_1 -> version 10.1
// 10MBPS -> speed 10
// STANDARD -> plan STANDARD
func getVersionSpeedPlanFromDescription(description string) (string, int, string) {
	strs := strings.Split(description, " ")
	version := strings.Join([]string{strs[3], strs[4]}, ".")
	speedString := strings.Trim(strs[4], "MBPS")
	speed, _ := strconv.Atoi(speedString)
	plan := strings.Trim(strs[5], "")

	return version, speed, plan
}

func resourceSoftLayerNetworkApplicationDeliveryControllerUpdate(d *schema.ResourceData, meta interface{}) error {
	//	client := meta.(*Client).networkApplicationDeliveryControllerService
	return fmt.Errorf("Update is not supported yet")
}

func resourceSoftLayerNetworkApplicationDeliveryControllerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).networkApplicationDeliveryControllerService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	_, err = client.DeleteObject(id)

	if err != nil {
		return fmt.Errorf("Error deleting network application delivery controller: %s", err)
	}

	return nil
}

func resourceSoftLayerNetworkApplicationDeliveryControllerExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*Client).networkApplicationDeliveryControllerService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	nadc, err := client.GetObject(id)

	if err != nil {
		return false, fmt.Errorf("Error fetching network application delivery controller: %s", err)
	}

	return nadc.Id == id && err == nil, nil
}
