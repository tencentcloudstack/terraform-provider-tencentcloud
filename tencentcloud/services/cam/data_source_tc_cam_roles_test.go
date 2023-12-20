package cam_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCamRolesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamRolesDatasourceBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.tencentcloud_cam_roles.roles", "role_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_roles.roles", "role_list.0.role_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_roles.roles", "role_list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_roles.roles", "role_list.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_roles.roles", "role_list.0.console_login"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_roles.roles", "role_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_roles.roles", "role_list.0.update_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_roles.roles", "role_list.0.document"),
				),
			},
		},
	})
}

const testAccCamRolesDatasourceBasic = tcacctest.DefaultCamVariables + `
data "tencentcloud_cam_roles" "roles" {
	name = var.cam_role_basic
}
`
