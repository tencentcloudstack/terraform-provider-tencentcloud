package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRumOfflineLogConfigDataSource -v
func TestAccTencentCloudRumOfflineLogConfigDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRumOfflineLogConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_offline_log_config.offlineLogConfig"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_offline_log_config.offlineLogConfig", "project_key", "ZEYrYfvaYQ30jRdmPx"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_rum_offline_log_config.offlineLogConfig", "unique_id_set.#"),
				),
			},
		},
	})
}

const testAccDataSourceRumOfflineLogConfig = `

data "tencentcloud_rum_offline_log_config" "offlineLogConfig" {
	project_key = "ZEYrYfvaYQ30jRdmPx"
}

`
