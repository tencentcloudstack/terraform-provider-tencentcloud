package tco_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudOrganizationOrgManagePolicyConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_ORGANIZATION)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgManagePolicyConfig,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_manage_policy_config.org_manage_policy_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_manage_policy_config.org_manage_policy_config", "organization_id", "45155"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_manage_policy_config.org_manage_policy_config", "policy_type", "SERVICE_CONTROL_POLICY"),
				),
			},
			{
				ResourceName:      "tencentcloud_organization_org_manage_policy_config.org_manage_policy_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationOrgManagePolicyConfig = `

resource "tencentcloud_organization_org_manage_policy_config" "org_manage_policy_config" {
  organization_id = 45155
  policy_type = "SERVICE_CONTROL_POLICY"
}

`
