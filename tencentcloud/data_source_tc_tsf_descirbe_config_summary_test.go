package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfDescirbeConfigSummaryDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDescirbeConfigSummaryDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_descirbe_config_summary.descirbe_config_summary")),
			},
		},
	})
}

const testAccTsfDescirbeConfigSummaryDataSource = `

data "tencentcloud_tsf_descirbe_config_summary" "descirbe_config_summary" {
  application_id = ""
  search_word = ""
  order_by = ""
  order_type = 
  config_tag_list = 
  disable_program_auth_check = 
  config_id_list = 
  }

`
