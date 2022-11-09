package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTdcpgCluster_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdcpgCluster,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tdcpg_cluster.cluster", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_tdcpg_cluster.cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdcpgCluster = `

resource "tencentcloud_tdcpg_cluster" "cluster" {
  zone = ""
  master_user_password = ""
  CPU = ""
  memory = ""
  vpc_id = ""
  subnet_id = ""
  pay_mode = ""
  cluster_name = ""
  d_b_version = ""
  instance_count = ""
  period = ""
  storage = ""
  project_id = ""
}

`
