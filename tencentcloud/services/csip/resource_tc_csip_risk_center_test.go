package csip

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
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "scan_item"),
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "scan_plan_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "assets"),
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "scan_plan_content"),
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "self_defining_assets"),
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "scan_from"),
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "task_advance_cfg"),
					resource.TestCheckResourceAttrSet("tencentcloud_csip_risk_center.example", "task_mode"),
				),
			},
		},
	})
}

const testAccCsipRiskCenter = `
resource "tencentcloud_csip_risk_center" "example" {
  task_name       = "tf_example"
  scan_asset_type = 0
  scan_item       = []
  scan_plan_type  = 1
  assets {
    asset_name    = ""
    instance_type = ""
    asset_type    = ""
    asset         = ""
    region        = ""
    arn           = ""
  }
  scan_plan_content    = ""
  self_defining_assets = []
  scan_from            = ""
  task_advance_cfg {
    vul_risk {
      risk_id = ""
      enable  = 0
    }
    weak_pwd_risk {
      check_item_id = ""
      enable        = 0
    }
    cfg_risk {
      item_id       = ""
      enable        = 0
      resource_type = ""
    }
  }
  task_mode = 0
}
`
