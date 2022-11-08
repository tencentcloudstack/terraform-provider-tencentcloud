package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqRabbitmqCluster_basic -v
func TestAccTencentCloudTdmqRabbitmqCluster_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqCluster,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_cluster.rabbitmq_cluster", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tdmq_rabbitmq_cluster.rabbitmq_cluster", "name", "rabbit-test-123"),
					resource.TestCheckResourceAttr("tencentcloud_tdmq_rabbitmq_cluster.rabbitmq_cluster", "remark", "Description"),
				),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rabbitmq_cluster.rabbitmqCluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdmqRabbitmqCluster = `

resource "tencentcloud_tdmq_rabbitmq_cluster" "rabbitmq_cluster" {
  name = "iac-rabbit-test"
  remark = "Description"
}

`
