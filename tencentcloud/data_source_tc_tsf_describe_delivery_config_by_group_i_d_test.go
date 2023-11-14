package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfDescribeDeliveryConfigByGroupIDDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDescribeDeliveryConfigByGroupIDDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_describe_delivery_config_by_group_i_d.describe_delivery_config_by_group_i_d")),
			},
		},
	})
}

const testAccTsfDescribeDeliveryConfigByGroupIDDataSource = `

data "tencentcloud_tsf_describe_delivery_config_by_group_i_d" "describe_delivery_config_by_group_i_d" {
  group_id = "group-yrjkln9v"
  }

`
