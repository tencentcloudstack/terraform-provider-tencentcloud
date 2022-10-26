package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMariadbDedicatedClusterDBInstance_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbDedicatedClusterDBInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_dedicatedcluster_db_instance.dedicatedcluster_db_instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_mariadb_dedicatedcluster_db_instance.dedicatedClusterDBInstance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbDedicatedClusterDBInstance = `

resource "tencentcloud_mariadb_dedicatedcluster_db_instance" "dedicatedcluster_db_instance" {
  goods_num = ""
  memory = ""
  storage = ""
  cluster_id = ""
  vpc_id = ""
  subnet_id = ""
  db_version_id = ""
  instance_name = ""
  tags = {
    "createdBy" = "terraform"
  }
}

`
