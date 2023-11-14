package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfDescribePublicConfigSummaryDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDescribePublicConfigSummaryDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_describe_public_config_summary.describe_public_config_summary")),
			},
		},
	})
}

const testAccTsfDescribePublicConfigSummaryDataSource = `

data "tencentcloud_tsf_describe_public_config_summary" "describe_public_config_summary" {
  search_word = ""
  order_by = ""
  order_type = 
  config_tag_list = 
  disable_program_auth_check = 
  config_id_list = 
  }

`
