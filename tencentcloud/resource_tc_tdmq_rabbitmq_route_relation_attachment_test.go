package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTdmqRabbitmqRouteRelationAttachment_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqRouteRelationAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_route_relation_attachment.rabbitmq_route_relation_attachment", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rabbitmq_route_relation_attachment.rabbitmqRouteRelationAttachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqRabbitmqRouteRelationAttachment = `

resource "tencentcloud_tdmq_rabbitmq_route_relation_attachment" "rabbitmq_route_relation_attachment" {
  cluster_id = ""
  v_host_id = ""
  source_exchange = ""
  dest_type = ""
  dest_value = ""
  remark = ""
  routing_key = ""
}

`
