package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfConfigSummaryDataSource_basic -v
func TestAccTencentCloudTsfConfigSummaryDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfConfigSummaryDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_config_summary.config_summary"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_config_summary.config_summary", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_config_summary.config_summary", "result.0.total_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_config_summary.config_summary", "result.0.content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_config_summary.config_summary", "result.0.content.0.application_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_config_summary.config_summary", "result.0.content.0.application_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_config_summary.config_summary", "result.0.content.0.config_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_config_summary.config_summary", "result.0.content.0.config_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_config_summary.config_summary", "result.0.content.0.config_version_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_config_summary.config_summary", "result.0.content.0.delete_flag"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_config_summary.config_summary", "result.0.content.0.last_update_time"),
				),
			},
		},
	})
}

const testAccTsfConfigSummaryDataSource = `

data "tencentcloud_tsf_config_summary" "config_summary" {
	application_id = "application-a24x29xv"
	search_word = "terraform"
	order_by = "last_update_time"
	order_type = 0
	disable_program_auth_check = true
	config_id_list = ["dcfg-y54wzk3a"]
}

`
