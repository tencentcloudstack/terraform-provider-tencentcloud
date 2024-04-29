package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOrganizationOrgManagePolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_ORGANIZATION)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgManagePolicy,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_manage_policy.org_manage_policy", "id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_manage_policy.org_manage_policy", "name", "iac-example-svc"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_manage_policy.org_manage_policy", "content"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_manage_policy.org_manage_policy", "type", "SERVICE_CONTROL_POLICY"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_manage_policy.org_manage_policy", "description", "Full access policy"),
				),
			},
			{
				Config: testAccOrganizationOrgManagePolicyUpdate,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_manage_policy.org_manage_policy", "id"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_manage_policy.org_manage_policy", "name", "iac-example-svc1"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_manage_policy.org_manage_policy", "content"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_manage_policy.org_manage_policy", "type", "SERVICE_CONTROL_POLICY"),
					resource.TestCheckResourceAttr("tencentcloud_organization_org_manage_policy.org_manage_policy", "description", "Full access policy"),
				),
			},
			{
				ResourceName:      "tencentcloud_organization_org_manage_policy.org_manage_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationOrgManagePolicy = `

resource "tencentcloud_organization_org_manage_policy" "org_manage_policy" {
  name = "iac-example-svc"
  content = "{\"version\":\"2.0\",\"statement\":[{\"effect\":\"allow\",\"action\":\"*\",\"resource\":\"*\"}]}"
  type = "SERVICE_CONTROL_POLICY"
  description = "Full access policy"
}

`
const testAccOrganizationOrgManagePolicyUpdate = `

resource "tencentcloud_organization_org_manage_policy" "org_manage_policy" {
  name = "iac-example-svc1"
  content = "{\"version\":\"2.0\",\"statement\":[{\"effect\":\"allow\",\"action\":\"*\",\"resource\":\"*\"}]}"
  type = "SERVICE_CONTROL_POLICY"
  description = "Full access policy"
}

`
