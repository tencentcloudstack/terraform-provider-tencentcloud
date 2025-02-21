package emr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudEmrYarnResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccEmrYarn,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("tencentcloud_emr_yarn.emr_yarn", "id"),
			),
		}, {
			ResourceName:      "tencentcloud_emr_yarn.emr_yarn",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccEmrYarn = `
resource "tencentcloud_emr_yarn" "emr_yarn" {
  instance_id = "emr-rzrochgp"
  enable_resource_schedule = true
  scheduler = "fair"
  fair_global_config {
    user_max_apps_default = 1000
  }
}
`
