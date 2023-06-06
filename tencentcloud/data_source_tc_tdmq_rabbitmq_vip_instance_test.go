package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixTdmqRabbitmqVipInstanceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqVipInstanceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_rabbitmq_vip_instance.rabbitmq_vip_instance"),
				),
			},
		},
	})
}

const testAccTdmqRabbitmqVipInstanceDataSource = `

data "tencentcloud_tdmq_rabbitmq_vip_instance" "rabbitmq_vip_instance" {
  filters {
		name = ""
		values = 

  }
  }

`
