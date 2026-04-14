package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudConfigRemediationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigRemediation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_config_remediation.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_config_remediation.example", "remediation_id"),
					resource.TestCheckResourceAttr("tencentcloud_config_remediation.example", "remediation_type", "SCF"),
					resource.TestCheckResourceAttr("tencentcloud_config_remediation.example", "invoke_type", "MANUAL_EXECUTION"),
				),
			},
			{
				Config: testAccConfigRemediationUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_config_remediation.example", "invoke_type", "NON_EXECUTION"),
				),
			},
			{
				ResourceName:      "tencentcloud_config_remediation.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccConfigRemediation = `
resource "tencentcloud_config_compliance_pack" "example" {
  compliance_pack_name = "tf-example-remediation"
  risk_level           = 2
  description          = "tf example for remediation test"

  config_rules {
    identifier = "cam-user-mfa-check"
    rule_name  = "CAM子用户开启MFA"
    risk_level = 2
  }
}

resource "tencentcloud_config_remediation" "example" {
  rule_id                 = tencentcloud_config_compliance_pack.example.compliance_pack_id
  remediation_type        = "SCF"
  remediation_template_id = "qcs::scf:ap-guangzhou:uin/100000005287:namespace/test/functions/my-remediation-func"
  invoke_type             = "MANUAL_EXECUTION"
  source_type             = "CUSTOM"
}
`

const testAccConfigRemediationUpdate = `
resource "tencentcloud_config_compliance_pack" "example" {
  compliance_pack_name = "tf-example-remediation"
  risk_level           = 2
  description          = "tf example for remediation test"

  config_rules {
    identifier = "cam-user-mfa-check"
    rule_name  = "CAM子用户开启MFA"
    risk_level = 2
  }
}

resource "tencentcloud_config_remediation" "example" {
  rule_id                 = tencentcloud_config_compliance_pack.example.compliance_pack_id
  remediation_type        = "SCF"
  remediation_template_id = "qcs::scf:ap-guangzhou:uin/100000005287:namespace/test/functions/my-remediation-func"
  invoke_type             = "NON_EXECUTION"
  source_type             = "CUSTOM"
}
`
