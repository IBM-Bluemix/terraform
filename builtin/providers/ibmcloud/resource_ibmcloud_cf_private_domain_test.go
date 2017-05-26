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

func TestAccIBMCloudCFPrivateDomain_Basic(t *testing.T) {
	var conf cfv2.PrivateDomainFields
	name := fmt.Sprintf("terraform%d.com", acctest.RandInt())
	updateName := fmt.Sprintf("terraformnew%d.com", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudCFPrivateDomain_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCloudCFPrivateDomainExists("ibmcloud_cf_private_domain.domain", &conf),
					resource.TestCheckResourceAttr("ibmcloud_cf_private_domain.domain", "name", name),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCloudCFPrivateDomain_updateName(updateName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibmcloud_cf_private_domain.domain", "name", updateName),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCFPrivateDomainExists(n string, obj *cfv2.PrivateDomainFields) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		privateDomainRepo := testAccProvider.Meta().(ClientSession).CloudFoundryPrivateDomainClient()
		privateDomainGUID := rs.Primary.ID

		prdomain, err := privateDomainRepo.Get(privateDomainGUID)
		if err != nil {
			return err
		}

		*obj = *prdomain
		return nil
	}
}

func testAccCheckIBMCloudCFPrivateDomainDestroy(s *terraform.State) error {
	privateDomainRepo := testAccProvider.Meta().(ClientSession).CloudFoundryPrivateDomainClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibmcloud_cf_private_domain" {
			continue
		}

		privateDomainGUID := rs.Primary.ID

		// Try to find the private domain
		_, err := privateDomainRepo.Get(privateDomainGUID)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("Error waiting for CF private domain (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIBMCloudCFPrivateDomain_basic(name string) string {
	return fmt.Sprintf(`
		
		data "ibmcloud_cf_org" "orgdata" {
			org    = "%s"
		}

		resource "ibmcloud_cf_private_domain" "domain" {
			name = "%s"
			org_guid = "${data.ibmcloud_cf_org.orgdata.id}"
		}
	`, cfOrganization, name)
}

func testAccCheckIBMCloudCFPrivateDomain_updateName(updateName string) string {
	return fmt.Sprintf(`
		
		data "ibmcloud_cf_org" "orgdata" {
			org    = "%s"
		}

		resource "ibmcloud_cf_private_domain" "domain" {
			name = "%s"
			org_guid = "${data.ibmcloud_cf_org.orgdata.id}"
		}
	`, cfOrganization, updateName)
}
