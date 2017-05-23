package ibmcloud

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

func TestAccIBMCloudInfraBareMetal_Basic(t *testing.T) {
	var bareMetal datatypes.Hardware

	hostname := acctest.RandString(16)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCloudInfraBareMetalDestroy,
		Steps: []resource.TestStep{
			{
				Config:  testAccCheckIBMCloudInfraBareMetalConfig_basic(hostname),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCloudInfraBareMetalExists("ibmcloud_infra_bare_metal.terraform-acceptance-test-1", &bareMetal),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_bare_metal.terraform-acceptance-test-1", "hostname", hostname),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_bare_metal.terraform-acceptance-test-1", "domain", "terraformuat.ibm.com"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_bare_metal.terraform-acceptance-test-1", "os_reference_code", "UBUNTU_16_64"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_bare_metal.terraform-acceptance-test-1", "datacenter", "dal01"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_bare_metal.terraform-acceptance-test-1", "network_speed", "100"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_bare_metal.terraform-acceptance-test-1", "hourly_billing", "true"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_bare_metal.terraform-acceptance-test-1", "private_network_only", "false"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_bare_metal.terraform-acceptance-test-1", "user_metadata", "{\"value\":\"newvalue\"}"),
					resource.TestCheckResourceAttr(
						"ibmcloud_infra_bare_metal.terraform-acceptance-test-1", "fixed_config_preset", "S1270_32GB_1X1TBSATA_NORAID"),
					CheckStringSet(
						"ibmcloud_infra_bare_metal.terraform-acceptance-test-1",
						"tags", []string{"collectd"},
					),
				),
			},

			{
				Config:  testAccCheckIBMCloudInfraBareMetalConfig_update(hostname),
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCloudInfraBareMetalExists("ibmcloud_infra_bare_metal.terraform-acceptance-test-1", &bareMetal),
					CheckStringSet(
						"ibmcloud_infra_bare_metal.terraform-acceptance-test-1",
						"tags", []string{"mesos-master"},
					),
				),
			},
		},
	})
}

func testAccCheckIBMCloudInfraBareMetalDestroy(s *terraform.State) error {
	service := services.GetHardwareService(testAccProvider.Meta().(ClientSession).SoftLayerSession())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibmcloud_infra_bare_metal" {
			continue
		}

		id, _ := strconv.Atoi(rs.Primary.ID)

		// Try to find the bare metal
		_, err := service.Id(id).GetObject()

		// Wait
		if err != nil {
			if apiErr, ok := err.(sl.Error); !ok || apiErr.StatusCode != 404 {
				return fmt.Errorf(
					"Error waiting for bare metal (%d) to be destroyed: %s",
					id, err)
			}
		}
	}

	return nil
}

func testAccCheckIBMCloudInfraBareMetalExists(n string, bareMetal *datatypes.Hardware) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No bare metal ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)

		if err != nil {
			return err
		}

		service := services.GetHardwareService(testAccProvider.Meta().(ClientSession).SoftLayerSession())
		bm, err := service.Id(id).GetObject()
		if err != nil {
			return err
		}

		fmt.Printf("The ID is %d", *bm.Id)

		if *bm.Id != id {
			return errors.New("Bare metal not found")
		}

		*bareMetal = bm

		return nil
	}
}

func testAccCheckIBMCloudInfraBareMetalConfig_basic(hostname string) string {
	return fmt.Sprintf(`
resource "ibmcloud_infra_bare_metal" "terraform-acceptance-test-1" {
    hostname = "%s"
    domain = "terraformuat.ibm.com"
    os_reference_code = "UBUNTU_16_64"
    datacenter = "dal01"
    network_speed = 100
    hourly_billing = true
	private_network_only = false
    user_metadata = "{\"value\":\"newvalue\"}"
    fixed_config_preset = "S1270_32GB_1X1TBSATA_NORAID"
    tags = ["collectd"]
}
`, hostname)
}

func testAccCheckIBMCloudInfraBareMetalConfig_update(hostname string) string {
	return fmt.Sprintf(`
resource "ibmcloud_infra_bare_metal" "terraform-acceptance-test-1" {
    hostname = "%s"
    domain = "terraformuat.ibm.com"
    os_reference_code = "UBUNTU_16_64"
    datacenter = "dal01"
    network_speed = 100
    hourly_billing = true
    private_network_only = false
    user_metadata = "{\"value\":\"newvalue\"}"
    fixed_config_preset = "S1270_32GB_1X1TBSATA_NORAID"
    tags = ["mesos-master"]
}
`, hostname)
}