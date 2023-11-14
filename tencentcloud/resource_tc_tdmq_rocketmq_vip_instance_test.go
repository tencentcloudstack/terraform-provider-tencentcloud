package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqRocketmqVipInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRocketmqVipInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rocketmq_vip_instance.rocketmq_vip_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rocketmq_vip_instance.rocketmq_vip_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqRocketmqVipInstance = `

resource "tencentcloud_tdmq_rocketmq_vip_instance" "rocketmq_vip_instance" {
  name = ""
  spec = ""
  node_count = 
  storage_size = 
  zone_ids = 
  vpc_info {
		vpc_id = ""
		subnet_id = ""

  }
  time_span = 
}

`
