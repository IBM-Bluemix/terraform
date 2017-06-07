package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccIBMStorageObject_Basic(t *testing.T) {
	var accountName string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMStorageObjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMStorageObjectConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMStorageObjectExists("ibm_storage_object.testacc_foobar", &accountName),
					testAccCheckIBMStorageObjectAttributes(&accountName),
				),
			},
		},
	})
}

func testAccCheckIBMStorageObjectDestroy(s *terraform.State) error {
	return nil
}

func testAccCheckIBMStorageObjectExists(n string, accountName *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		*accountName = rs.Primary.ID

		return nil
	}
}

func testAccCheckIBMStorageObjectAttributes(accountName *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if *accountName == "" {
			return fmt.Errorf("No object storage account name")
		}

		return nil
	}
}

var testAccCheckIBMStorageObjectConfig_basic = `
resource "ibm_storage_object" "testacc_foobar" {
}`
