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

func TestAccIBMCloudCFServiceInstance_Basic(t *testing.T) {
	var conf cfv2.ServiceInstanceFields
	serviceName := fmt.Sprintf("terraform_%d", acctest.RandInt())
	updateName := fmt.Sprintf("terraform_%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCloudInfraScaleGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudCFServiceInstance_basic(serviceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCloudCFServiceInstanceExists("ibmcloud_cf_service_instance.service", &conf),
					resource.TestCheckResourceAttr("ibmcloud_cf_service_instance.service", "name", serviceName),
					resource.TestCheckResourceAttr("ibmcloud_cf_service_instance.service", "service", "cleardb"),
					resource.TestCheckResourceAttr("ibmcloud_cf_service_instance.service", "plan", "spark"),
					resource.TestCheckResourceAttr("ibmcloud_cf_service_instance.service", "tags.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCloudCFServiceInstance_updateWithSameName(serviceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCloudCFServiceInstanceExists("ibmcloud_cf_service_instance.service", &conf),
					resource.TestCheckResourceAttr("ibmcloud_cf_service_instance.service", "name", serviceName),
					resource.TestCheckResourceAttr("ibmcloud_cf_service_instance.service", "service", "cleardb"),
					resource.TestCheckResourceAttr("ibmcloud_cf_service_instance.service", "plan", "spark"),
					resource.TestCheckResourceAttr("ibmcloud_cf_service_instance.service", "tags.#", "3"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCloudCFServiceInstance_update(updateName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibmcloud_cf_service_instance.service", "name", updateName),
					resource.TestCheckResourceAttr("ibmcloud_cf_service_instance.service", "service", "cleardb"),
					resource.TestCheckResourceAttr("ibmcloud_cf_service_instance.service", "plan", "spark"),
					resource.TestCheckResourceAttr("ibmcloud_cf_service_instance.service", "tags.#", "1"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCloudCFServiceInstance_newServiceType(updateName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibmcloud_cf_service_instance.service", "name", updateName),
					resource.TestCheckResourceAttr("ibmcloud_cf_service_instance.service", "service", "cloudantNOSQLDB"),
					resource.TestCheckResourceAttr("ibmcloud_cf_service_instance.service", "plan", "Lite"),
					resource.TestCheckResourceAttr("ibmcloud_cf_service_instance.service", "tags.#", "1"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCFServiceInstanceDestroy(s *terraform.State) error {
	serviceRepo := testAccProvider.Meta().(ClientSession).CloudFoundryServiceInstanceClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibmcloud_cf_service_instance" {
			continue
		}

		serviceGuid := rs.Primary.ID

		// Try to find the key
		_, err := serviceRepo.Get(serviceGuid)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("Error waiting for CF service (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIBMCloudCFServiceInstanceExists(n string, obj *cfv2.ServiceInstanceFields) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		serviceRepo := testAccProvider.Meta().(ClientSession).CloudFoundryServiceInstanceClient()
		serviceGuid := rs.Primary.ID

		service, err := serviceRepo.Get(serviceGuid)
		if err != nil {
			return err
		}

		*obj = *service
		return nil
	}
}

func testAccCheckIBMCloudCFServiceInstance_basic(serviceName string) string {
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
	`, cfSpace, cfOrganization, serviceName)
}

func testAccCheckIBMCloudCFServiceInstance_updateWithSameName(serviceName string) string {
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
			tags               = ["cluster-service","cluster-bind","db"]
		}
	`, cfSpace, cfOrganization, serviceName)
}

func testAccCheckIBMCloudCFServiceInstance_update(updateName string) string {
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
			tags               = ["cluster-service"]
		}
	`, cfSpace, cfOrganization, updateName)
}

func testAccCheckIBMCloudCFServiceInstance_newServiceType(updateName string) string {
	return fmt.Sprintf(`
		data "ibmcloud_cf_space" "spacedata" {
			space  = "%s"
			org    = "%s"
		}

		resource "ibmcloud_cf_service_instance" "service" {
			name              = "%s"
			space_guid        = "${data.ibmcloud_cf_space.spacedata.id}"
			service           = "cloudantNOSQLDB"
			plan              = "Lite"
			tags               = ["cluster-service"]
		}
	`, cfSpace, cfOrganization, updateName)
}
