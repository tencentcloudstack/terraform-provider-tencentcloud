package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqRocketmqGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRocketmqGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_group.group", "id")),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rocketmq_group.group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqRocketmqGroup = `

resource "tencentcloud_tdmq_rocketmq_group" "group" {
  group_id = &lt;nil&gt;
  namespaces = &lt;nil&gt;
  read_enable = &lt;nil&gt;
  broadcast_enable = &lt;nil&gt;
  cluster_id = &lt;nil&gt;
  remark = &lt;nil&gt;
                    }

`
