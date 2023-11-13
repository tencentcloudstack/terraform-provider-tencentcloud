package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqRabbitmqVirtualHostResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqVirtualHost,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_virtual_host.rabbitmq_virtual_host", "id")),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rabbitmq_virtual_host.rabbitmq_virtual_host",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqRabbitmqVirtualHost = `

resource "tencentcloud_tdmq_rabbitmq_virtual_host" "rabbitmq_virtual_host" {
  instance_id = ""
  virtual_host = ""
  description = ""
  trace_flag = 
}

`
