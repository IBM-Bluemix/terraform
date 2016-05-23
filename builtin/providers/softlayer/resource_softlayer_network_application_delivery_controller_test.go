package softlayer

import (
	"fmt"
	"strconv"
	"testing"

	datatypes "github.com/TheWeatherCompany/softlayer-go/data_types"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	// "time"
)

func TestAccSoftLayerNetworkApplicationDeliveryController_Basic(t *testing.T) {
	var nadc datatypes.SoftLayer_Network_Application_Delivery_Controller

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSoftLayerNetworkApplicationDeliveryControllerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSoftLayerNetworkApplicationDeliveryControllerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSoftLayerNetworkApplicationDeliveryControllerExists("softlayer_network_application_delivery_controller.testacc_foobar_nadc", &nadc),
					resource.TestCheckResourceAttr(
						"softlayer_network_application_delivery_controller.testacc_foobar_nadc", "name", "nadc_test_name"),
					resource.TestCheckResourceAttr(
						"softlayer_network_application_delivery_controller.testacc_foobar_nadc", "type", "Netscaler VPX"),
					resource.TestCheckResourceAttr(
						"softlayer_network_application_delivery_controller.testacc_foobar_nadc", "datacenter", "DALLAS06"),
					resource.TestCheckResourceAttr(
						"softlayer_network_application_delivery_controller.testacc_foobar_nadc", "speed", "10"),
					resource.TestCheckResourceAttr(
						"softlayer_network_application_delivery_controller.testacc_foobar_nadc", "version", "10.1"),
					resource.TestCheckResourceAttr(
						"softlayer_network_application_delivery_controller.testacc_foobar_nadc", "plan", "Standard"),
					resource.TestCheckResourceAttr(
						"softlayer_network_application_delivery_controller.testacc_foobar_nadc", "ip_count", "2"),
				),
			},
		},
	})
}

func testAccCheckSoftLayerNetworkApplicationDeliveryControllerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).networkApplicationDeliveryControllerService

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "softlayer_network_application_delivery_controller" {
			continue
		}

		id, _ := strconv.Atoi(rs.Primary.ID)

		// Try to find the application delivery controller
		_, err := client.GetObject(id)

		if err == nil {
			return fmt.Errorf("Application Delivery Controller still exists")
		}
	}

	return nil
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
resource "softlayer_network_application_delivery_controller" "testacc_foobar_nadc" {
    name = "nadc_test_name"
    type = "Netscaler VPX"
    datacenter = "DALLAS06"
    speed = 10
    version = "10.1"
    plan = "Standard"
    ip_count = 2
}`
