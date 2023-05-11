package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfPublicConfigSummaryDataSource_basic -v
func TestAccTencentCloudTsfPublicConfigSummaryDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfPublicConfigSummaryDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_public_config_summary.public_config_summary"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_public_config_summary.public_config_summary", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_public_config_summary.public_config_summary", "result.0.content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_public_config_summary.public_config_summary", "result.0.content.0.config_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_public_config_summary.public_config_summary", "result.0.content.0.config_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_public_config_summary.public_config_summary", "result.0.content.0.delete_flag"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_public_config_summary.public_config_summary", "result.0.content.0.last_update_time"),
				),
			},
		},
	})
}

const testAccTsfPublicConfigSummaryDataSource = `

data "tencentcloud_tsf_public_config_summary" "public_config_summary" {
	search_word = "test"
	order_by = "last_update_time"
	order_type = 0
	disable_program_auth_check = true
}

`
