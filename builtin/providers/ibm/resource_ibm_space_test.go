package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/IBM-Bluemix/bluemix-go/bmxerror"
)

func TestAccIBMSpace_Basic(t *testing.T) {
	var conf cfv2.SpaceFields
	name := fmt.Sprintf("terraform_%d", acctest.RandInt())
	updatedName := fmt.Sprintf("terraform_updated_%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMSpaceDestroy,
		Steps: []resource.TestStep{

			resource.TestStep{
				Config: testAccCheckIBMSpaceCreate(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMSpaceExists("ibm_space.space", &conf),
					resource.TestCheckResourceAttr("ibm_space.space", "org", cfOrganization),
					resource.TestCheckResourceAttr("ibm_space.space", "name", name),
				),
			},

			resource.TestStep{
				Config: testAccCheckIBMSpaceUpdate(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_space.space", "org", cfOrganization),
					resource.TestCheckResourceAttr("ibm_space.space", "name", updatedName),
				),
			},
		},
	})
}

func testAccCheckIBMSpaceExists(n string, obj *cfv2.SpaceFields) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		cfClient, err := testAccProvider.Meta().(ClientSession).CFAPI()
		if err != nil {
			return err
		}
		spaceGUID := rs.Primary.ID

		space, err := cfClient.Spaces().Get(spaceGUID)
		if err != nil {
			return err
		}

		*obj = *space
		return nil
	}
}

func testAccCheckIBMSpaceDestroy(s *terraform.State) error {
	cfClient, err := testAccProvider.Meta().(ClientSession).CFAPI()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_space" {
			continue
		}

		spaceGUID := rs.Primary.ID
		_, err := cfClient.Spaces().Get(spaceGUID)

		if err != nil {
			if apierr, ok := err.(bmxerror.RequestFailure); ok && apierr.StatusCode() != 404 {
				return fmt.Errorf("Error waiting for Space (%s) to be destroyed: %s", rs.Primary.ID, err)
			}
		}
	}
	return nil
}

func testAccCheckIBMSpaceCreate(name string) string {
	return fmt.Sprintf(`
	
resource "ibm_space" "space" {
    org = "%s"
	name = "%s"
}`, cfOrganization, name)

}

func testAccCheckIBMSpaceUpdate(updatedName string) string {
	return fmt.Sprintf(`
	
resource "ibm_space" "space" {
    org = "%s"
	name = "%s"
}`, cfOrganization, updatedName)

}
