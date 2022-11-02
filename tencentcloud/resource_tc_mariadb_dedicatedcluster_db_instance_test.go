package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMariadbDedicatedclusterDbInstance_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbDedicatedclusterDbInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_dedicatedcluster_db_instance.dedicatedcluster_db_instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_mariadb_dedicatedcluster_db_instance.dedicatedclusterDbInstance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbDedicatedclusterDbInstance = `

resource "tencentcloud_mariadb_dedicatedcluster_db_instance" "dedicatedcluster_db_instance" {
  goods_num = ""
  memory = ""
  storage = ""
  cluster_id = ""
  vpc_id = ""
  subnet_id = ""
  db_version_id = ""
  instance_name = ""
}

`
