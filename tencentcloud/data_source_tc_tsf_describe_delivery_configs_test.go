package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfDescribeDeliveryConfigsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDescribeDeliveryConfigsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_describe_delivery_configs.describe_delivery_configs")),
			},
		},
	})
}

const testAccTsfDescribeDeliveryConfigsDataSource = `

data "tencentcloud_tsf_describe_delivery_configs" "describe_delivery_configs" {
  search_word = ""
  }

`
