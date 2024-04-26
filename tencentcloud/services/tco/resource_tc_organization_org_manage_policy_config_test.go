package tco

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudOrganizationOrgManagePolicyConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgManagePolicyConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_manage_policy_config.org_manage_policy_config", "id")),
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
