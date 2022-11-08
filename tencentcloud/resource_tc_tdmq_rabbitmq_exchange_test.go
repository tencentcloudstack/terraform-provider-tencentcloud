package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTdmqRabbitmqExchange_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqExchange,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_exchange.rabbitmq_exchange", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rabbitmq_exchange.rabbitmqExchange",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqRabbitmqExchange = `

resource "tencentcloud_tdmq_rabbitmq_exchange" "rabbitmq_exchange" {
  exchange = ""
  v_host_id = ""
  type = ""
  cluster_id = ""
  remark = ""
  alternate_exchange = ""
  delayed_type = ""
}

`
