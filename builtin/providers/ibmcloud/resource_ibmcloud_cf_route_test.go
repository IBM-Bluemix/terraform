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

func TestAccIBMCloudCFRoute_Basic(t *testing.T) {
	var conf cfv2.RouteFields
	host := fmt.Sprintf("terraform_%d", acctest.RandInt())
	updateHost := fmt.Sprintf("terraform_%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCloudCFRouteDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudCFRoute_basic(host),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCloudCFRouteExists("ibmcloud_cf_route.route", &conf),
					resource.TestCheckResourceAttr("ibmcloud_cf_route.route", "host", host),
					resource.TestCheckResourceAttr("ibmcloud_cf_route.route", "path", "/app"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCloudCFRoute_updatePath(host),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCloudCFRouteExists("ibmcloud_cf_route.route", &conf),
					resource.TestCheckResourceAttr("ibmcloud_cf_route.route", "host", host),
					resource.TestCheckResourceAttr("ibmcloud_cf_route.route", "path", "/app1"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCloudCFRoute_updateHost(updateHost),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibmcloud_cf_route.route", "host", updateHost),
					resource.TestCheckResourceAttr("ibmcloud_cf_route.route", "path", ""),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCFRouteDestroy(s *terraform.State) error {
	routeRepo := testAccProvider.Meta().(ClientSession).CloudFoundryRouteClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibmcloud_cf_route" {
			continue
		}

		routeGuid := rs.Primary.ID

		// Try to find the key
		_, err := routeRepo.Get(routeGuid)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("Error waiting for CF route (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIBMCloudCFRouteExists(n string, obj *cfv2.RouteFields) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		routeRepo := testAccProvider.Meta().(ClientSession).CloudFoundryRouteClient()
		routeGuid := rs.Primary.ID

		route, err := routeRepo.Get(routeGuid)
		if err != nil {
			return err
		}

		*obj = *route
		return nil
	}
}

func testAccCheckIBMCloudCFRoute_basic(host string) string {
	return fmt.Sprintf(`
	
		data "ibmcloud_cf_space" "spacedata" {
			org    = "%s"
			space  = "%s"
		}
		
		data "ibmcloud_cf_domain" "domain" {
			name        = "mybluemix.net"
		}
		
		resource "ibmcloud_cf_route" "route" {
			domain_guid       = "${data.ibmcloud_cf_domain.domain.id}"
			space_guid        = "${data.ibmcloud_cf_space.spacedata.id}"
			host              = "%s"
			path              = "/app"
		}
	`, cfOrganization, cfSpace, host)
}

func testAccCheckIBMCloudCFRoute_updatePath(host string) string {
	return fmt.Sprintf(`
	
		data "ibmcloud_cf_space" "spacedata" {
			org    = "%s"
			space  = "%s"
		}
		
		data "ibmcloud_cf_domain" "domain" {
			name        = "mybluemix.net"
		}
		
		resource "ibmcloud_cf_route" "route" {
			domain_guid       = "${data.ibmcloud_cf_domain.domain.id}"
			space_guid        = "${data.ibmcloud_cf_space.spacedata.id}"
			host              = "%s"
			path              = "/app1"
		}
	`, cfOrganization, cfSpace, host)
}

func testAccCheckIBMCloudCFRoute_updateHost(updateHost string) string {
	return fmt.Sprintf(`
		
		data "ibmcloud_cf_space" "spacedata" {
			org    = "%s"
			space  = "%s"
		}
		
		data "ibmcloud_cf_domain" "domain" {
			name        = "mybluemix.net"
		}
		
		resource "ibmcloud_cf_route" "route" {
			domain_guid       = "${data.ibmcloud_cf_domain.domain.id}"
			space_guid        = "${data.ibmcloud_cf_space.spacedata.id}"
			host              = "%s"
		}
	`, cfOrganization, cfSpace, updateHost)
}
