package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqRabbitmqVipInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqVipInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.rabbitmq_vip_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rabbitmq_vip_instance.rabbitmq_vip_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqRabbitmqVipInstance = `

resource "tencentcloud_tdmq_rabbitmq_vip_instance" "rabbitmq_vip_instance" {
  zone_ids = 
  vpc_id = ""
  subnet_id = ""
  cluster_name = ""
  node_spec = ""
  node_num = 
  storage_size = 
  enable_create_default_ha_mirror_queue = 
  auto_renew_flag = 
  time_span = 
}

`
