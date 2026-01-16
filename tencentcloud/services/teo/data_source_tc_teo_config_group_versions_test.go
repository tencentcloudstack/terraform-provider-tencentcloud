package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoConfigGroupVersionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoConfigGroupVersionsDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_config_group_versions.teo_config_group_versions"),
				resource.TestCheckResourceAttr("data.tencentcloud_teo_config_group_versions.teo_config_group_versions", "zone_id", "zone-2xkazzl8yf6k"),
				resource.TestCheckResourceAttr("data.tencentcloud_teo_config_group_versions.teo_config_group_versions", "group_id", "cg-3lchxitnb5pb"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_config_group_versions.teo_config_group_versions", "config_group_version_infos.#"),
			),
		}},
	})
}

const testAccTeoConfigGroupVersionsDataSource = `

data "tencentcloud_teo_config_group_versions" "teo_config_group_versions" {
  zone_id = "zone-2xkazzl8yf6k"
  group_id = "cg-3lchxitnb5pb"
}
`
