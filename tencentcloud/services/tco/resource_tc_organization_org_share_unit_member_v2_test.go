package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOrganizationOrgShareUnitMemberV2Resource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgShareUnitMemberV2,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_share_unit_member_v2.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_share_unit_member_v2.example", "unit_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_share_unit_member_v2.example", "area"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_share_unit_member_v2.example", "members.#"),
				),
			},
			{
				Config: testAccOrganizationOrgShareUnitMemberV2Update,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_org_share_unit_member_v2.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_share_unit_member_v2.example", "unit_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_share_unit_member_v2.example", "area"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_org_share_unit_member_v2.example", "members.#"),
				),
			},
			{
				ResourceName:      "tencentcloud_organization_org_share_unit_member_v2.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationOrgShareUnitMemberV2 = `
resource "tencentcloud_organization_org_share_unit" "example" {
  name        = "tf-example"
  area        = "ap-guangzhou"
  description = "description."
}

resource "tencentcloud_organization_org_share_unit_member_v2" "example" {
  unit_id = tencentcloud_organization_org_share_unit.example.unit_id
  area    = tencentcloud_organization_org_share_unit.example.area
  members {
    share_member_uin = 100040906282
  }

  members {
    share_member_uin = 100043984945
  }

  members {
    share_member_uin = 100043985088
  }

  members {
    share_member_uin = 100042287843
  }

  members {
    share_member_uin = 100042287853
  }
}
`

const testAccOrganizationOrgShareUnitMemberV2Update = `
resource "tencentcloud_organization_org_share_unit" "example" {
  name        = "tf-example"
  area        = "ap-guangzhou"
  description = "description."
}

resource "tencentcloud_organization_org_share_unit_member_v2" "example" {
  unit_id = tencentcloud_organization_org_share_unit.example.unit_id
  area    = tencentcloud_organization_org_share_unit.example.area
  members {
    share_member_uin = 100038833157
  }

  members {
    share_member_uin = 100043984945
  }

  members {
    share_member_uin = 100042287843
  }

  members {
    share_member_uin = 100042287853
  }
}
`
