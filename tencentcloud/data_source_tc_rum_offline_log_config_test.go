package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudRumOfflineLogConfigDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRumOfflineLogConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_offline_log_config.offline_log_config"),
				),
			},
		},
	})
}

const testAccDataSourceRumOfflineLogConfig = `

data "tencentcloud_rum_offline_log_config" "offline_log_config" {
  project_key = ""
}

`
