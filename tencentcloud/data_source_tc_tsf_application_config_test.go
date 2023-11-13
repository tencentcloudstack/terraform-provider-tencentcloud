package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfApplicationConfigDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationConfigDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_application_config.application_config")),
			},
		},
	})
}

const testAccTsfApplicationConfigDataSource = `

data "tencentcloud_tsf_application_config" "application_config" {
  application_id = "app-123456"
  config_id = "config-123456"
  config_id_list = 
  config_name = "test-config"
  config_version = "1.0"
  }

`
