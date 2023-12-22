package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApplicationPublicConfigDataSource_basic -v
func TestAccTencentCloudTsfApplicationPublicConfigDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationPublicConfigDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_application_public_config.application_public_config"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_public_config.application_public_config", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_public_config.application_public_config", "result.0.total_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_public_config.application_public_config", "result.0.content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_public_config.application_public_config", "result.0.content.0.config_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_public_config.application_public_config", "result.0.content.0.config_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_public_config.application_public_config", "result.0.content.0.config_value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_public_config.application_public_config", "result.0.content.0.config_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_public_config.application_public_config", "result.0.content.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_public_config.application_public_config", "result.0.content.0.delete_flag"),
				),
			},
		},
	})
}

const testAccTsfApplicationPublicConfigDataSource = `

data "tencentcloud_tsf_application_public_config" "application_public_config" {
	# config_id = "dcfg-p-evjrbgly"
	# # config_id_list = [""]
	# config_name = "dsadsa"
	# config_version = "123"
}

`
