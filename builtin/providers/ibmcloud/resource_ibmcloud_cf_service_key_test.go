package ibmcloud

import (
	"fmt"
	"testing"

	"strings"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
)

func TestAccIBMCloudCFServiceKey_Basic(t *testing.T) {
	var conf cfv2.ServiceKeyFields
	serviceName := fmt.Sprintf("terraform_%d", acctest.RandInt())
	serviceKey := fmt.Sprintf("terraform_%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudCFServiceKey_basic(serviceName, serviceKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCloudCFServiceKeyExists("ibmcloud_cf_service_key.serviceKey", &conf),
					resource.TestCheckResourceAttr("ibmcloud_cf_service_key.serviceKey", "name", serviceKey),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCFServiceKeyExists(n string, obj *cfv2.ServiceKeyFields) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		serviceRepo := testAccProvider.Meta().(ClientSession).CloudFoundryServiceKeyClient()
		serviceKeyGuid := rs.Primary.ID

		serviceKey, err := serviceRepo.Get(serviceKeyGuid)
		if err != nil {
			return err
		}

		*obj = *serviceKey
		return nil
	}
}

func testAccCheckIBMCloudCFServiceKeyDestroy(s *terraform.State) error {
	serviceKeyRepo := testAccProvider.Meta().(ClientSession).CloudFoundryServiceKeyClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibmcloud_cf_service_key" {
			continue
		}

		serviceKeyGuid := rs.Primary.ID

		// Try to find the key
		_, err := serviceKeyRepo.Get(serviceKeyGuid)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("Error waiting for CF service key (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIBMCloudCFServiceKey_basic(serviceName, serviceKey string) string {
	return fmt.Sprintf(`
		
		data "ibmcloud_cf_space" "spacedata" {
			space  = "%s"
			org    = "%s"
		}
		
		resource "ibmcloud_cf_service_instance" "service" {
			name              = "%s"
			space_guid        = "${data.ibmcloud_cf_space.spacedata.id}"
			service           = "cleardb"
			plan              = "spark"
			tags               = ["cluster-service","cluster-bind"]
		}

		resource "ibmcloud_cf_service_key" "serviceKey" {
			name = "%s"
			service_instance_guid = "${ibmcloud_cf_service_instance.service.id}"
		}
	`, cfSpace, cfOrganization, serviceName, serviceKey)
}
