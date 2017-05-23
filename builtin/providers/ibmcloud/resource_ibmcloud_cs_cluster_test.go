package ibmcloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	bluemix "github.com/IBM-Bluemix/bluemix-go"
	"github.com/IBM-Bluemix/bluemix-go/session"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"

	"github.com/IBM-Bluemix/bluemix-go/api/account/accountv2"
	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	v1 "github.com/IBM-Bluemix/bluemix-go/api/k8scluster/k8sclusterv1"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccIBMCloudCluster_basic(t *testing.T) {
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandInt())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCloudCSClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMCloudInfraCluster_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibmcloud_cs_cluster.testacc_cluster", "name", clusterName),
					resource.TestCheckResourceAttr(
						"ibmcloud_cs_cluster.testacc_cluster", "worker_num", "1"),
				),
			},
		},
	})
}

func testAccCheckIBMCloudCSClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(ClientSession).ClusterClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibmcloud_cs_cluster" {
			continue
		}

		targetEnv := getClusterTargetHeaderTestACC()
		// Try to find the key
		_, err := client.Find(rs.Primary.ID, targetEnv)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("Error waiting for cluster (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func getClusterTargetHeaderTestACC() *v1.ClusterTargetHeader {
	org := cfOrganization
	space := cfSpace
	c := new(bluemix.Config)
	sess, err := session.New(c)
	if err != nil {
		log.Fatal(err)
	}

	client, err := cfv2.New(sess)

	if err != nil {
		log.Fatal(err)
	}

	orgAPI := client.Organizations()
	myorg, err := orgAPI.FindByName(org)

	if err != nil {
		log.Fatal(err)
	}

	spaceAPI := client.Spaces()
	myspace, err := spaceAPI.FindByNameInOrg(myorg.GUID, space)

	if err != nil {
		log.Fatal(err)
	}

	accClient, err := accountv2.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	accountAPI := accClient.Accounts()
	myAccount, err := accountAPI.FindByOrg(myorg.GUID, c.Region)
	if err != nil {
		log.Fatal(err)
	}

	target := &v1.ClusterTargetHeader{
		OrgID:     myorg.GUID,
		SpaceID:   myspace.GUID,
		AccountID: myAccount.GUID,
	}

	return target
}

func testAccCheckIBMCloudInfraCluster_basic(clusterName string) string {
	return fmt.Sprintf(`

data "ibmcloud_cf_org" "org" {
    org = "%s"
}

data "ibmcloud_cf_space" "space" {
  org    = "%s"
  space  = "%s"
}

data "ibmcloud_cf_account" "acc" {
   org_guid = "${data.ibmcloud_cf_org.org.id}"
}

resource "ibmcloud_cs_cluster" "testacc_cluster" {
  name       = "%s"
  datacenter = "dal10"

  org_guid = "${data.ibmcloud_cf_org.org.id}"
	space_guid = "${data.ibmcloud_cf_space.space.id}"
	account_guid = "${data.ibmcloud_cf_account.acc.id}"

  workers = [{
    name = "worker1"

    action = "add"
  }]

  machine_type    = "free"
  isolation       = "public"
  public_vlan_id  = "vlan"
  private_vlan_id = "vlan"
}	`, cfOrganization, cfOrganization, cfSpace, clusterName)
}
