package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
)

func TestAccIBMCloudCFSpace_Basic(t *testing.T) {
	var conf cfv2.SpaceFields
	name := fmt.Sprintf("terraform_%d", acctest.RandInt())
	updatedName := fmt.Sprintf("terraform_updated_%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCloudCFSpaceDestroy,
		Steps: []resource.TestStep{

			resource.TestStep{
				Config: testAccCheckIBMCloudCFSpaceCreate(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCloudCFSpaceExists("ibmcloud_cf_space.space", &conf),
					resource.TestCheckResourceAttr("ibmcloud_cf_space.space", "org", cfOrganization),
					resource.TestCheckResourceAttr("ibmcloud_cf_space.space", "name", name),
				),
			},

			resource.TestStep{
				Config: testAccCheckIBMCloudCFSpaceUpdate(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibmcloud_cf_space.space", "org", cfOrganization),
					resource.TestCheckResourceAttr("ibmcloud_cf_space.space", "name", updatedName),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCFSpaceExists(n string, obj *cfv2.SpaceFields) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		spaceClient := testAccProvider.Meta().(ClientSession).CloudFoundrySpaceClient()
		spaceGUID := rs.Primary.ID

		space, err := spaceClient.Get(spaceGUID)
		if err != nil {
			return err
		}

		*obj = *space
		return nil
	}
}

func testAccCheckIBMCloudCFSpaceDestroy(s *terraform.State) error {
	spaceClient := testAccProvider.Meta().(ClientSession).CloudFoundrySpaceClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibmcloud_cf_space" {
			continue
		}

		spaceGUID := rs.Primary.ID
		_, err := spaceClient.Get(spaceGUID)

		if err != nil {
			if apierr, ok := err.(bmxerror.RequestFailure); ok && apierr.StatusCode() != 404 {
				return fmt.Errorf("Error waiting for Space (%s) to be destroyed: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}

func testAccCheckIBMCloudCFSpaceCreate(name string) string {
	return fmt.Sprintf(`
	
resource "ibmcloud_cf_space" "space" {
    org = "%s"
	name = "%s"
}`, cfOrganization, name)

}

func testAccCheckIBMCloudCFSpaceUpdate(updatedName string) string {
	return fmt.Sprintf(`
	
resource "ibmcloud_cf_space" "space" {
    org = "%s"
	name = "%s"
}`, cfOrganization, updatedName)

}
