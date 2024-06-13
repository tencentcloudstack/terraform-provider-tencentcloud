package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOrganizationOrgManagePolicyTargetResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_ORGANIZATION)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgManagePolicyTarget,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_manage_policy_target.org_manage_policy_target", "id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_manage_policy_target.org_manage_policy_target", "target_id", "100034649025"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_manage_policy_target.org_manage_policy_target", "target_type", "MEMBER"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_manage_policy_target.org_manage_policy_target", "policy_id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_manage_policy_target.org_manage_policy_target", "policy_type", "SERVICE_CONTROL_POLICY")),
			},
			{
				ResourceName:      "tencentcloud_organization_org_manage_policy_target.org_manage_policy_target",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationOrgManagePolicyTarget = `
resource "tencentcloud_organization_org_manage_policy" "org_manage_policy" {
  name = "example-service"
  content = "{\"version\":\"2.0\",\"statement\":[{\"effect\":\"allow\",\"action\":\"*\",\"resource\":\"*\"}]}"
  type = "SERVICE_CONTROL_POLICY"
  description = "Full access policy"
  depends_on = [tencentcloud_organization_org_manage_policy_config.org_manage_policy_config]
}
resource "tencentcloud_organization_org_manage_policy_config" "org_manage_policy_config" {
  organization_id = 45155
  policy_type = "SERVICE_CONTROL_POLICY"
}
resource "tencentcloud_organization_org_manage_policy_target" "org_manage_policy_target" {
  target_id = 100034649025
  target_type = "MEMBER"
  policy_id = tencentcloud_organization_org_manage_policy.org_manage_policy.policy_id
  policy_type = "SERVICE_CONTROL_POLICY"
}

`
