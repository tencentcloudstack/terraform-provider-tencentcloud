package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTdmqRabbitmqVirtualHostListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqVirtualHostListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_rabbitmq_virtual_host_list.rabbitmq_virtual_host_list")),
			},
		},
	})
}

const testAccTdmqRabbitmqVirtualHostListDataSource = `

data "tencentcloud_tdmq_rabbitmq_virtual_host_list" "rabbitmq_virtual_host_list" {
  instance_id = ""
  }

`
