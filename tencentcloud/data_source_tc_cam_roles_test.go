package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCamRolesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamRolesDataSource_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCamRoleExists("tencentcloud_cam_role.role"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_roles.roles", "role_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_roles.roles", "role_list.0.role_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_roles.roles", "role_list.0.name", "cam-role-test11"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_roles.roles", "role_list.0.description", "test"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_roles.roles", "role_list.0.console_login", "true"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_roles.roles", "role_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_roles.roles", "role_list.0.update_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_roles.roles", "role_list.0.document"),
				),
			},
		},
	})
}

const testAccCamRolesDataSource_basic = `
resource "tencentcloud_cam_role" "role" {
  name          = "cam-role-test11"
  document      = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"qcs\":[\"qcs::cam::uin/100009461222:uin/100009461222\"]}}]}"
  description   = "test"
  console_login = true
}
  
data "tencentcloud_cam_roles" "roles" {
  role_id = tencentcloud_cam_role.role.id
}
`
