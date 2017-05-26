package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMCloudCFPrivateDomainDataSource_basic(t *testing.T) {
	name := fmt.Sprintf("terraform%d.com", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCloudCFPrivateDomainDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ibmcloud_cf_private_domain.testacc_domain", "id"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCFPrivateDomainDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	
		data "ibmcloud_cf_org" "orgdata" {
			org    = "%s"
		}

		resource "ibmcloud_cf_private_domain" "domain" {
			name = "%s"
			org_guid = "${data.ibmcloud_cf_org.orgdata.id}"
		}
	
		data "ibmcloud_cf_private_domain" "testacc_domain" {
			name = "${ibmcloud_cf_private_domain.domain.name}"
		}`, cfOrganization, name)

}
