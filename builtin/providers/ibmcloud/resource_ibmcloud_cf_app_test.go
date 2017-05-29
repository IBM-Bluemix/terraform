package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccIBMCloudCFApp_Basic(t *testing.T) {
	var conf cfv2.AppFields
	name := fmt.Sprintf("terraform_%d", acctest.RandInt())
	updatedName := fmt.Sprintf("terraform_updated_%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCloudCFAppDestroy,
		Steps: []resource.TestStep{

			resource.TestStep{
				Config: testAccCheckIBMCloudCFAppCreate(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCloudCFAppExists("ibmcloud_cf_app.app", &conf),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "name", name),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCloudCFAppUpdate(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "name", updatedName),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "instances", "2"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "memory", "128"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "disk_quota", "512"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "environment_json.test", "test1"),
				),
			},
		},
	})
}

func TestAccIBMCloudCFApp_with_routes(t *testing.T) {
	var conf cfv2.AppFields
	name := fmt.Sprintf("terraform_%d", acctest.RandInt())
	route1 := fmt.Sprintf("terraform_%d", acctest.RandInt())
	route2 := fmt.Sprintf("terraform_%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCloudCFAppDestroy,
		Steps: []resource.TestStep{

			resource.TestStep{
				Config: testAccCheckIBMCloudCFAppBindRoute(name, route1),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCloudCFAppExists("ibmcloud_cf_app.app", &conf),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "name", name),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "instances", "1"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "memory", "128"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "disk_quota", "512"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "environment_json.test", "test1"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "route_guid.#", "1"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCloudCFAppAddMultipleRoute(name, route1, route2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCloudCFAppExists("ibmcloud_cf_app.app", &conf),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "name", name),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "instances", "1"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "memory", "128"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "disk_quota", "512"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "environment_json.test", "test1"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "route_guid.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCloudCFAppUnBindRoute(name, route1, route2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCloudCFAppExists("ibmcloud_cf_app.app", &conf),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "name", name),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "instances", "1"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "memory", "128"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "disk_quota", "512"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "environment_json.test", "test1"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "route_guid.#", "1"),
				),
			},
		},
	})

}

func TestAccIBMCloudCFApp_with_service_instances(t *testing.T) {
	var conf cfv2.AppFields
	name := fmt.Sprintf("terraform_%d", acctest.RandInt())
	route := fmt.Sprintf("terraform_%d", acctest.RandInt())
	serviceName1 := fmt.Sprintf("terraform_%d", acctest.RandInt())
	serviceName2 := fmt.Sprintf("terraform_%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCloudCFAppDestroy,
		Steps: []resource.TestStep{

			resource.TestStep{
				Config: testAccCheckIBMCloudCFAppBindService(name, route, serviceName1),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCloudCFAppExists("ibmcloud_cf_app.app", &conf),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "name", name),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "instances", "1"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "memory", "128"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "disk_quota", "512"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "environment_json.test", "test1"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "route_guid.#", "1"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "service_instance_guid.#", "1"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCloudCFAppAddMultipleService(name, route, serviceName1, serviceName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCloudCFAppExists("ibmcloud_cf_app.app", &conf),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "name", name),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "instances", "1"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "memory", "128"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "disk_quota", "512"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "environment_json.test", "test1"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "route_guid.#", "1"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "service_instance_guid.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCloudCFAppUnBindService(name, route, serviceName1, serviceName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCloudCFAppExists("ibmcloud_cf_app.app", &conf),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "name", name),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "instances", "1"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "memory", "128"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "disk_quota", "512"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "environment_json.test", "test1"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "route_guid.#", "1"),
					resource.TestCheckResourceAttr("ibmcloud_cf_app.app", "service_instance_guid.#", "1"),
				),
			},
		},
	})

}

func testAccCheckIBMCloudCFAppDestroy(s *terraform.State) error {

	appClient := testAccProvider.Meta().(ClientSession).CloudFoundryAppClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibmcloud_cf_app" {
			continue
		}
		appGUID := rs.Primary.ID

		_, err := appClient.Get(appGUID)
		if err == nil {
			return fmt.Errorf("App still exists: %s", rs.Primary.ID)
		}
	}

	return nil

}

func testAccCheckIBMCloudCFAppExists(n string, obj *cfv2.AppFields) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		appClient := testAccProvider.Meta().(ClientSession).CloudFoundryAppClient()
		appGUID := rs.Primary.ID

		app, err := appClient.Get(appGUID)
		if err != nil {
			return err
		}

		*obj = *app
		return nil
	}
}

func testAccCheckIBMCloudCFAppCreate(name string) string {
	return fmt.Sprintf(`

data "ibmcloud_cf_space" "space" {
  org    = "%s"
  space  = "%s"
}
resource "ibmcloud_cf_app" "app" {
	name = "%s"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	app_path = "test-fixtures/app1.zip"
	wait_time_minutes = 90
	buildpack = "sdk-for-nodejs"
}`, cfOrganization, cfSpace, name)

}

func testAccCheckIBMCloudCFAppUpdate(name string) string {
	return fmt.Sprintf(`

data "ibmcloud_cf_space" "space" {
   org    = "%s"
  space  = "%s"
}
resource "ibmcloud_cf_app" "app" {
	name = "%s"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	app_path = "test-fixtures/app1.zip"
	wait_time_minutes = 20
	buildpack = "sdk-for-nodejs"
	instances = 1
	disk_quota = 512
	memory = 128
	instances = 2
	environment_json = {
		"test" = "test1"
	}
}`, cfOrganization, cfSpace, name)

}

func testAccCheckIBMCloudCFAppBindRoute(name, route1 string) string {
	return fmt.Sprintf(`

data "ibmcloud_cf_space" "space" {
   org    = "%s"
  space  = "%s"
}

data "ibmcloud_cf_shared_domain" "domain" {
			name        = "mybluemix.net"
}
		
resource "ibmcloud_cf_route" "route" {
		domain_guid       = "${data.ibmcloud_cf_shared_domain.domain.id}"
		space_guid        = "${data.ibmcloud_cf_space.space.id}"
		host              = "%s"
		
}

resource "ibmcloud_cf_app" "app" {
	name = "%s"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	app_path = "test-fixtures/app1.zip"
	wait_time_minutes = 20
	buildpack = "sdk-for-nodejs"
	instances = 1
	route_guid = ["${ibmcloud_cf_route.route.id}"]
	disk_quota = 512
	memory = 128
	instances = 1
	environment_json = {
		"test" = "test1"
	}
}`, cfOrganization, cfSpace, route1, name)

}

func testAccCheckIBMCloudCFAppAddMultipleRoute(name, route1, route2 string) string {
	return fmt.Sprintf(`

data "ibmcloud_cf_space" "space" {
   org    = "%s"
  space  = "%s"
}

data "ibmcloud_cf_shared_domain" "domain" {
			name        = "mybluemix.net"
}
		
resource "ibmcloud_cf_route" "route" {
		domain_guid       = "${data.ibmcloud_cf_shared_domain.domain.id}"
		space_guid        = "${data.ibmcloud_cf_space.space.id}"
		host              = "%s"
		
}

resource "ibmcloud_cf_route" "route1" {
		domain_guid       = "${data.ibmcloud_cf_shared_domain.domain.id}"
		space_guid        = "${data.ibmcloud_cf_space.space.id}"
		host              = "%s"
		
}

resource "ibmcloud_cf_app" "app" {
	name = "%s"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	app_path = "test-fixtures/app1.zip"
	wait_time_minutes = 20
	buildpack = "sdk-for-nodejs"
	instances = 1
	route_guid = ["${ibmcloud_cf_route.route.id}","${ibmcloud_cf_route.route1.id}"]
	disk_quota = 512
	memory = 128
	instances = 1
	disk_quota = 512
	environment_json = {
		"test" = "test1"
	}
}`, cfOrganization, cfSpace, route1, route2, name)

}

func testAccCheckIBMCloudCFAppUnBindRoute(name, route1, route2 string) string {
	return fmt.Sprintf(`

data "ibmcloud_cf_space" "space" {
   org    = "%s"
  space  = "%s"
}

data "ibmcloud_cf_shared_domain" "domain" {
			name        = "mybluemix.net"
}
		
resource "ibmcloud_cf_route" "route" {
		domain_guid       = "${data.ibmcloud_cf_shared_domain.domain.id}"
		space_guid        = "${data.ibmcloud_cf_space.space.id}"
		host              = "%s"
		
}

resource "ibmcloud_cf_route" "route1" {
		domain_guid       = "${data.ibmcloud_cf_shared_domain.domain.id}"
		space_guid        = "${data.ibmcloud_cf_space.space.id}"
		host              = "%s"
		
}

resource "ibmcloud_cf_app" "app" {
	name = "%s"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	app_path = "test-fixtures/app1.zip"
	wait_time_minutes = 20
	buildpack = "sdk-for-nodejs"
	instances = 1
	route_guid = ["${ibmcloud_cf_route.route.id}"]
	disk_quota = 512
	memory = 128
	instances = 1
	disk_quota = 512
	environment_json = {
		"test" = "test1"
	}
}`, cfOrganization, cfSpace, route1, route2, name)

}

func testAccCheckIBMCloudCFAppBindService(name, route1, serviceName string) string {
	return fmt.Sprintf(`

data "ibmcloud_cf_space" "space" {
   org    = "%s"
  space  = "%s"
}

data "ibmcloud_cf_shared_domain" "domain" {
			name        = "mybluemix.net"
}
		
resource "ibmcloud_cf_route" "route" {
		domain_guid       = "${data.ibmcloud_cf_shared_domain.domain.id}"
		space_guid        = "${data.ibmcloud_cf_space.space.id}"
		host              = "%s"
		
}

resource "ibmcloud_cf_service_instance" "service" {
			name              = "%s"
			space_guid        = "${data.ibmcloud_cf_space.space.id}"
			service           = "cleardb"
			plan              = "spark"
			tags               = ["cluster-service","cluster-bind"]
		}

resource "ibmcloud_cf_app" "app" {
	name = "%s"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	app_path = "test-fixtures/app1.zip"
	wait_time_minutes = 20
	buildpack = "sdk-for-nodejs"
	instances = 1
	route_guid = ["${ibmcloud_cf_route.route.id}"]
	service_instance_guid = ["${ibmcloud_cf_service_instance.service.id}"]
	disk_quota = 512
	memory = 128
	instances = 1
	environment_json = {
		"test" = "test1"
	}
}`, cfOrganization, cfSpace, route1, serviceName, name)

}

func testAccCheckIBMCloudCFAppAddMultipleService(name, route, serviceName1, serviceName2 string) string {
	return fmt.Sprintf(`

data "ibmcloud_cf_space" "space" {
   org    = "%s"
  space  = "%s"
}

data "ibmcloud_cf_shared_domain" "domain" {
			name        = "mybluemix.net"
}
		
resource "ibmcloud_cf_route" "route" {
		domain_guid       = "${data.ibmcloud_cf_shared_domain.domain.id}"
		space_guid        = "${data.ibmcloud_cf_space.space.id}"
		host              = "%s"
		
}

resource "ibmcloud_cf_service_instance" "service" {
			name              = "%s"
			space_guid        = "${data.ibmcloud_cf_space.space.id}"
			service           = "cleardb"
			plan              = "spark"
			tags               = ["cluster-service","cluster-bind"]
		}
		resource "ibmcloud_cf_service_instance" "service1" {
			name              = "%s"
			space_guid        = "${data.ibmcloud_cf_space.space.id}"
			service           = "cloudantNOSQLDB"
			plan              = "Lite"
			tags               = ["cluster-service"]
		}

resource "ibmcloud_cf_app" "app" {
	name = "%s"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	app_path = "test-fixtures/app1.zip"
	wait_time_minutes = 20
	buildpack = "sdk-for-nodejs"
	instances = 1
	route_guid = ["${ibmcloud_cf_route.route.id}"]
	service_instance_guid = ["${ibmcloud_cf_service_instance.service.id}","${ibmcloud_cf_service_instance.service1.id}"]
	disk_quota = 512
	memory = 128
	instances = 1
	disk_quota = 512
	environment_json = {
		"test" = "test1"
	}
}`, cfOrganization, cfSpace, route, serviceName1, serviceName2, name)

}

func testAccCheckIBMCloudCFAppUnBindService(name, route1, serviceName1, serviceName2 string) string {
	return fmt.Sprintf(`

data "ibmcloud_cf_space" "space" {
   org    = "%s"
  space  = "%s"
}

data "ibmcloud_cf_shared_domain" "domain" {
			name        = "mybluemix.net"
}
		
resource "ibmcloud_cf_route" "route" {
		domain_guid       = "${data.ibmcloud_cf_shared_domain.domain.id}"
		space_guid        = "${data.ibmcloud_cf_space.space.id}"
		host              = "%s"
		
}

resource "ibmcloud_cf_service_instance" "service" {
			name              = "%s"
			space_guid        = "${data.ibmcloud_cf_space.space.id}"
			service           = "cleardb"
			plan              = "spark"
			tags               = ["cluster-service","cluster-bind"]
}
resource "ibmcloud_cf_service_instance" "service1" {
			name              = "%s"
			space_guid        = "${data.ibmcloud_cf_space.space.id}"
			service           = "cloudantNOSQLDB"
			plan              = "Lite"
			tags               = ["cluster-service"]
}

resource "ibmcloud_cf_app" "app" {
	name = "%s"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	app_path = "test-fixtures/app1.zip"
	wait_time_minutes = 20
	buildpack = "sdk-for-nodejs"
	instances = 1
	route_guid = ["${ibmcloud_cf_route.route.id}"]
	service_instance_guid = ["${ibmcloud_cf_service_instance.service.id}"]
	disk_quota = 512
	memory = 128
	instances = 1
	disk_quota = 512
	environment_json = {
		"test" = "test1"
	}
}`, cfOrganization, cfSpace, route1, serviceName1, serviceName2, name)

}
