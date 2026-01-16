package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoDeployConfigVersionHistoryDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoDeployConfigVersionHistoryDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_deploy_config_version_history.teo_deploy_config_version_history"),
				resource.TestCheckResourceAttr("data.tencentcloud_teo_deploy_config_version_history.teo_deploy_config_version_history", "zone_id", "zone-2xkazzl8yf6k"),
				resource.TestCheckResourceAttr("data.tencentcloud_teo_deploy_config_version_history.teo_deploy_config_version_history", "env_id", "env-3lchxiq1h855"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_deploy_config_version_history.teo_deploy_config_version_history", "records.#"),
			),
		}},
	})
}

const testAccTeoDeployConfigVersionHistoryDataSource = `

data "tencentcloud_teo_deploy_config_version_history" "teo_deploy_config_version_history" {
  zone_id = "zone-2xkazzl8yf6k"
  env_id = "env-3lchxiq1h855"
}
`
