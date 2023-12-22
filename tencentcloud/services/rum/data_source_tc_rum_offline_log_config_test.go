package rum_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRumOfflineLogConfigDataSource -v
func TestAccTencentCloudRumOfflineLogConfigDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRumOfflineLogConfig,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_rum_offline_log_config.offlineLogConfig"),
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
