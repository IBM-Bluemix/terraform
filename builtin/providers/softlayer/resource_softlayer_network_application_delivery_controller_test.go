package softlayer

import (
	"fmt"
	"strconv"
	"testing"

	datatypes "github.com/TheWeatherCompany/softlayer-go/data_types"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSoftLayerNetworkApplicationDeliveryController_Basic(t *testing.T) {
	var nadc datatypes.SoftLayer_Network_Application_Delivery_Controller

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSoftLayerNetworkApplicationDeliveryControllerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSoftLayerNetworkApplicationDeliveryControllerExists("softlayer_lb_vpx.testacc_foobar_vpx", &nadc),
					resource.TestCheckResourceAttr(
						"softlayer_lb_vpx.testacc_foobar_vpx", "type", "NetScaler VPX"),
					resource.TestCheckResourceAttr(
						"softlayer_lb_vpx.testacc_foobar_vpx", "datacenter", "dal06"),
					resource.TestCheckResourceAttr(
						"softlayer_lb_vpx.testacc_foobar_vpx", "speed", "10"),
					resource.TestCheckResourceAttr(
						"softlayer_lb_vpx.testacc_foobar_vpx", "plan", "Standard"),
					resource.TestCheckResourceAttr(
						"softlayer_lb_vpx.testacc_foobar_vpx", "ip_count", "2"),
					resource.TestCheckResourceAttr(
						"softlayer_lb_vpx.testacc_foobar_vpx", "version", "10.1"),
					resource.TestCheckResourceAttr(
						"softlayer_lb_vpx.testacc_foobar_vpx", "front_end_vlan.vlan_number", "1251"),
					resource.TestCheckResourceAttr(
						"softlayer_lb_vpx.testacc_foobar_vpx", "front_end_vlan.primary_router_hostname", "fcr01a.dal06"),
					resource.TestCheckResourceAttr(
						"softlayer_lb_vpx.testacc_foobar_vpx", "back_end_vlan.vlan_number", "1540"),
					resource.TestCheckResourceAttr(
						"softlayer_lb_vpx.testacc_foobar_vpx", "back_end_vlan.primary_router_hostname", "bcr01a.dal06"),
					resource.TestCheckResourceAttr(
						"softlayer_lb_vpx.testacc_foobar_vpx", "front_end_subnet", "23.246.226.248/29"),
					resource.TestCheckResourceAttr(
						"softlayer_lb_vpx.testacc_foobar_vpx", "back_end_subnet", "10.107.180.0/26"),
					resource.TestCheckResourceAttr(
						"softlayer_lb_vpx.testacc_foobar_vpx", "vip_pool.#", "2"),
				),
			},
		},
	})
}

func testAccCheckSoftLayerNetworkApplicationDeliveryControllerExists(n string, nadc *datatypes.SoftLayer_Network_Application_Delivery_Controller) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		nadcId, _ := strconv.Atoi(rs.Primary.ID)

		client := testAccProvider.Meta().(*Client).networkApplicationDeliveryControllerService
		found, err := client.GetObject(nadcId)

		if err != nil {
			return err
		}

		if strconv.Itoa(int(found.Id)) != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*nadc = found

		return nil
	}
}

const testAccCheckSoftLayerNetworkApplicationDeliveryControllerConfig_basic = `
resource "softlayer_lb_vpx" "testacc_foobar_vpx" {
    datacenter = "dal06"
    speed = 10
    version = "10.1"
    plan = "Standard"
    ip_count = 2
    front_end_vlan {
       vlan_number = 1251
       primary_router_hostname = "fcr01a.dal06"
    }
    back_end_vlan {
       vlan_number = 1540
       primary_router_hostname = "bcr01a.dal06"
    }
    front_end_subnet = "23.246.226.248/29"
    back_end_subnet = "10.107.180.0/26"
}`
