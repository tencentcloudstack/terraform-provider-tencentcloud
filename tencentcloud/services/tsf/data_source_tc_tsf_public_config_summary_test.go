package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfPublicConfigSummaryDataSource_basic -v
func TestAccTencentCloudTsfPublicConfigSummaryDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfPublicConfigSummaryDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_public_config_summary.public_config_summary"),
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
