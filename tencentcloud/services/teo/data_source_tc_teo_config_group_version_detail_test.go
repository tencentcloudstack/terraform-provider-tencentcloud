package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoConfigGroupVersionDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoConfigGroupVersionDetailDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_config_group_version_detail.teo_config_group_version_detail"),
				resource.TestCheckResourceAttr("data.tencentcloud_teo_config_group_version_detail.teo_config_group_version_detail", "zone_id", "zone-2xkazzl8yf6k"),
				resource.TestCheckResourceAttr("data.tencentcloud_teo_config_group_version_detail.teo_config_group_version_detail", "version_id", "ver-3lchxizh2mqn"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_config_group_version_detail.teo_config_group_version_detail", "config_group_version_info.#"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_config_group_version_detail.teo_config_group_version_detail", "content"),
			),
		}},
	})
}

const testAccTeoConfigGroupVersionDetailDataSource = `

data "tencentcloud_teo_config_group_version_detail" "teo_config_group_version_detail" {
  zone_id = "zone-2xkazzl8yf6k"
  version_id = "ver-3lchxizh2mqn"
}
`
