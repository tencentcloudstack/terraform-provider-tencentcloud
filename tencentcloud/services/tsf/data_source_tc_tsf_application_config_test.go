package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApplicationConfigDataSource_basic -v
func TestAccTencentCloudTsfApplicationConfigDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationConfigDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_application_config.application_config"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_config.application_config", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_config.application_config", "result.0.content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_config.application_config", "result.0.content.0.application_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_config.application_config", "result.0.content.0.application_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_config.application_config", "result.0.content.0.config_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_config.application_config", "result.0.content.0.config_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_config.application_config", "result.0.content.0.config_value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_config.application_config", "result.0.content.0.config_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_config.application_config", "result.0.content.0.config_version_desc"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_config.application_config", "result.0.content.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_config.application_config", "result.0.content.0.delete_flag"),
				),
			},
		},
	})
}

const testAccTsfApplicationConfigDataSourceVar = `
variable "application_id" {
	default = "` + tcacctest.DefaultTsfApplicationId + `"
}

variable "config_id" {
	default = "` + tcacctest.DefaultTsfConfigId + `"
}

`

const testAccTsfApplicationConfigDataSource = testAccTsfApplicationConfigDataSourceVar + `

data "tencentcloud_tsf_application_config" "application_config" {
	application_id = var.application_id
	config_id = var.config_id
	# config_id_list =
	config_name = "keep-terraform-testing"
	config_version = "v1"
}

`
