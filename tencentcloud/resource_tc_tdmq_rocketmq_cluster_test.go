package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqRocketmqClusterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRocketmqCluster,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_cluster.cluster", "id")),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rocketmq_cluster.cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqRocketmqCluster = `

resource "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
  cluster_name = &lt;nil&gt;
  remark = &lt;nil&gt;
                  }

`
