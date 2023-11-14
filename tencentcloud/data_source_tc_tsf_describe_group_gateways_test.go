package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfDescribeGroupGatewaysDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDescribeGroupGatewaysDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_describe_group_gateways.describe_group_gateways")),
			},
		},
	})
}

const testAccTsfDescribeGroupGatewaysDataSource = `

data "tencentcloud_tsf_describe_group_gateways" "describe_group_gateways" {
  gateway_deploy_group_id = ""
  search_word = ""
  }

`
