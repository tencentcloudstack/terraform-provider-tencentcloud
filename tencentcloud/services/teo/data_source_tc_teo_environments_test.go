package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoEnvironmentsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoEnvironmentsDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_teo_environments.teo_environments"),
				resource.TestCheckResourceAttr("data.tencentcloud_teo_environments.teo_environments", "zone_id", "zone-2xkazzl8yf6k"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_teo_environments.teo_environments", "env_infos.#"),
			),
		}},
	})
}

const testAccTeoEnvironmentsDataSource = `

data "tencentcloud_teo_environments" "teo_environments" {
  zone_id = "zone-2xkazzl8yf6k"
}
`
