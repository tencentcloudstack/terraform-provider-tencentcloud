package csip_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudCsipRiskCenterResource_basic -v
func TestAccTencentCloudCsipRiskCenterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCsipRiskCenter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_csip_risk_center.example", "task_name", "tf_example"),
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "scan_asset_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "scan_item.#"),
					resource.TestCheckResourceAttr("tencentcloud_csip_risk_center.example", "scan_plan_content", "0 0 0 */1 * * *"),
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "task_mode"),
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "assets.#"),
				),
			},
			{
				Config: testAccCsipRiskCenterUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_csip_risk_center.example", "task_name", "tf_example_update"),
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "scan_asset_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "scan_item.#"),
					resource.TestCheckResourceAttr("tencentcloud_csip_risk_center.example", "scan_plan_content", "0 0 0 */1 * * *"),
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "task_mode"),
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "assets.#"),
				),
			},
		},
	})
}

const testAccCsipRiskCenter = `
resource "tencentcloud_csip_risk_center" "example" {
  task_name         = "tf_example"
  scan_plan_type    = 0
  scan_asset_type   = 2
  scan_item         = ["port", "poc", "weakpass"]
  scan_plan_content = "0 0 0 */1 * * *"
  task_mode         = 0

  assets {
    asset_name    = "iac-test"
    instance_type = "1"
    asset_type    = "PublicIp"
    asset         = "49.232.172.248"
    region        = "ap-beijing"
  }
}
`

const testAccCsipRiskCenterUpdate = `
resource "tencentcloud_csip_risk_center" "example" {
  task_name         = "tf_example_update"
  scan_plan_type    = 0
  scan_asset_type   = 2
  scan_item         = ["port", "poc", "weakpass"]
  scan_plan_content = "0 0 0 */1 * * *"
  task_mode         = 0

  assets {
    asset_name    = "iac-test"
    instance_type = "1"
    asset_type    = "PublicIp"
    asset         = "49.232.172.248"
    region        = "ap-beijing"
  }
}
`
