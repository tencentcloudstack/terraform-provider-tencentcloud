package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfBusinessLogConfigsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfBusinessLogConfigsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_business_log_configs.business_log_configs")),
			},
		},
	})
}

const testAccTsfBusinessLogConfigsDataSource = `

data "tencentcloud_tsf_business_log_configs" "business_log_configs" {
  search_word = ""
  disable_program_auth_check = 
  config_id_list = 
  }

`
