package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApplicationFileConfigDataSource_basic -v
func TestAccTencentCloudTsfApplicationFileConfigDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationFileConfigDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_application_file_config.application_file_config"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_file_config.application_file_config", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_file_config.application_file_config", "result.0.total_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_file_config.application_file_config", "result.0.content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_file_config.application_file_config", "result.0.content.0.application_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_file_config.application_file_config", "result.0.content.0.application_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_file_config.application_file_config", "result.0.content.0.config_file_code"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_file_config.application_file_config", "result.0.content.0.config_file_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_file_config.application_file_config", "result.0.content.0.config_file_path"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_file_config.application_file_config", "result.0.content.0.config_file_value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_file_config.application_file_config", "result.0.content.0.config_file_value_length"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_file_config.application_file_config", "result.0.content.0.config_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_file_config.application_file_config", "result.0.content.0.config_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_file_config.application_file_config", "result.0.content.0.config_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_file_config.application_file_config", "result.0.content.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application_file_config.application_file_config", "result.0.content.0.delete_flag"),
				),
			},
		},
	})
}

const testAccTsfApplicationFileConfigDataSource = `

data "tencentcloud_tsf_application_file_config" "application_file_config" {
	# config_id = "dcfg-f-4y4ekzqv"
	# config_id_list = [""]
	# config_name = "file-log1"
	# application_id = "application-2vzk6n3v"
	# config_version = "1.2"
}

`
