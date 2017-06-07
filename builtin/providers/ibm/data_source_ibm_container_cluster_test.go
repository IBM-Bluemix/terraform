package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIBMContainerClusterDataSource_basic(t *testing.T) {
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandInt())
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMContainerClusterDataSource(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ibm_container_cluster.testacc_ds_cluster", "worker_count", "1"),
				),
			},
		},
	})
}

func testAccCheckIBMContainerClusterDataSource(clusterName string) string {
	return fmt.Sprintf(`
data "ibm_org" "testacc_ds_org" {
    org = "%s"
}

data "ibm_space" "testacc_ds_space" {
    org = "%s"
	space = "%s"
}

data "ibm_account" "testacc_acc" {
    org_guid = "${data.ibm_org.testacc_ds_org.id}"
}


resource "ibm_container_cluster" "testacc_cluster" {
    name = "%s"
    datacenter = "dal10"
	org_guid = "${data.ibm_org.testacc_ds_org.id}"
	space_guid = "${data.ibm_space.testacc_ds_space.id}"
	account_guid = "${data.ibm_account.testacc_acc.id}"

   workers = [{
    name = "worker1"

    action = "add"
  }]
	machine_type = "free"
	isolation = "public"
	public_vlan_id = "vlan"
	private_vlan_id = "vlan"
}
data "ibm_container_cluster" "testacc_ds_cluster" {
	org_guid = "${data.ibm_org.testacc_ds_org.id}"
	space_guid = "${data.ibm_space.testacc_ds_space.id}"
	account_guid = "${data.ibm_account.testacc_acc.id}"
    cluster_name_id = "${ibm_container_cluster.testacc_cluster.id}"
}
`, cfOrganization, cfOrganization, cfSpace, clusterName)
}
