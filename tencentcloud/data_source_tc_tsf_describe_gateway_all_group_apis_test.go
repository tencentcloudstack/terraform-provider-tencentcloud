package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfDescribeGatewayAllGroupApisDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDescribeGatewayAllGroupApisDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_describe_gateway_all_group_apis.describe_gateway_all_group_apis")),
			},
		},
	})
}

const testAccTsfDescribeGatewayAllGroupApisDataSource = `

data "tencentcloud_tsf_describe_gateway_all_group_apis" "describe_gateway_all_group_apis" {
  gateway_deploy_group_id = "group-aeoej4qy"
  search_word = ""
  }

`
