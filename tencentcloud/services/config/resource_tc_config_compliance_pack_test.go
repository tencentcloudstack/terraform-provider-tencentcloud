package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudConfigCompliancePackResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigCompliancePack,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_config_compliance_pack.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_config_compliance_pack.example", "compliance_pack_name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_config_compliance_pack.example", "risk_level", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_config_compliance_pack.example", "compliance_pack_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_config_compliance_pack.example", "create_time"),
				),
			},
			{
				Config: testAccConfigCompliancePackUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_config_compliance_pack.example", "compliance_pack_name", "tf-example-update"),
					resource.TestCheckResourceAttr("tencentcloud_config_compliance_pack.example", "risk_level", "1"),
					resource.TestCheckResourceAttr("tencentcloud_config_compliance_pack.example", "description", "updated description"),
				),
			},
			{
				ResourceName:      "tencentcloud_config_compliance_pack.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccConfigCompliancePack = `
resource "tencentcloud_config_compliance_pack" "example" {
  compliance_pack_name = "tf-example"
  risk_level           = 2
  description          = "tf example compliance pack"

  config_rules {
    identifier = "cam-user-mfa-check"
    rule_name  = "CAM子用户开启MFA"
    risk_level = 2
  }
}
`

const testAccConfigCompliancePackUpdate = `
resource "tencentcloud_config_compliance_pack" "example" {
  compliance_pack_name = "tf-example-update"
  risk_level           = 1
  description          = "updated description"

  config_rules {
    identifier = "cam-user-mfa-check"
    rule_name  = "CAM子用户开启MFA"
    risk_level = 2
  }

  config_rules {
    identifier = "cam-user-group-bound"
    rule_name  = "CAM访问管理子用户必须关联用户组"
    risk_level = 3
  }
}
`
