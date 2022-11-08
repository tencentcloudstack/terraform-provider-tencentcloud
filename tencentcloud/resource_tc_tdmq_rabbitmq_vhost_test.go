package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqRabbitmqVhost_basic -v
func TestAccTencentCloudTdmqRabbitmqVhost_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqVhost,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vhost.rabbitmq_vhost", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tdmq_rabbitmq_vhost.rabbitmq_vhost", "vhost_id", "vhost-name"),
					resource.TestCheckResourceAttr("tencentcloud_tdmq_rabbitmq_vhost.rabbitmq_vhost", "msg_ttl", "86400000"),
					resource.TestCheckResourceAttr("tencentcloud_tdmq_rabbitmq_vhost.rabbitmq_vhost", "remark", "vhost-desc"),
				),
			},
			// {
			// 	Config: testAccTdmqRabbitmqVhostUpdet,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vhost.rabbitmq_vhost", "id"),
			// 		resource.TestCheckResourceAttr("tencentcloud_tdmq_rabbitmq_vhost.rabbitmq_vhost", "vhost_id", "vhost-name"),
			// 		resource.TestCheckResourceAttr("tencentcloud_tdmq_rabbitmq_vhost.rabbitmq_vhost", "msg_ttl", "3600000"),
			// 		resource.TestCheckResourceAttr("tencentcloud_tdmq_rabbitmq_vhost.rabbitmq_vhost", "remark", "vhost-desc-1"),
			// 	),
			// },
			{
				ResourceName:      "tencentcloud_tdmq_rabbitmq_vhost.rabbitmq_vhost",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqRabbitmqVhost = `

resource "tencentcloud_tdmq_rabbitmq_vhost" "rabbitmq_vhost" {
	cluster_id = "amqp-jz55mrq2exd9"
	vhost_id = "vhost-name"
	msg_ttl = 86400000
	remark = "vhost-desc"
}

`

const testAccTdmqRabbitmqVhostUpdet = `

resource "tencentcloud_tdmq_rabbitmq_vhost" "rabbitmq_vhost" {
	cluster_id = "amqp-jz55mrq2exd9"
	vhost_id = "vhost-name"
	msg_ttl = 3600000
	remark = "vhost-desc-1"
}

`
