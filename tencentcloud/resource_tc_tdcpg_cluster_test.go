package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdcpgClusterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdcpgCluster,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tdcpg_cluster.cluster", "id")),
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
  zone = &lt;nil&gt;
  master_user_password = &lt;nil&gt;
  c_p_u = &lt;nil&gt;
  memory = &lt;nil&gt;
  vpc_id = &lt;nil&gt;
  subnet_id = &lt;nil&gt;
  pay_mode = &lt;nil&gt;
  cluster_name = &lt;nil&gt;
  d_b_version = &lt;nil&gt;
  instance_count = &lt;nil&gt;
  period = &lt;nil&gt;
  storage = &lt;nil&gt;
  project_id = &lt;nil&gt;
}

`
