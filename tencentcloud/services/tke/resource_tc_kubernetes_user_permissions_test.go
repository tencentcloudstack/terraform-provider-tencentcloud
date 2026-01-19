package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesUserPermissionsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesUserPermissions,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_user_permissions.example", "id"),
				),
			},
			{
				Config: testAccKubernetesUserPermissionsUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_user_permissions.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_kubernetes_user_permissions.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccKubernetesUserPermissions = `
resource "tencentcloud_kubernetes_user_permissions" "example" {
  target_uin = ""
  permissions {
    cluster_id = ""
    role_name  = ""
    role_type  = ""
    is_custom  = ""
    namespace  = ""
  }
}
`

const testAccKubernetesUserPermissionsUpdate = `
resource "tencentcloud_kubernetes_user_permissions" "example" {
  target_uin = ""
  permissions {
    cluster_id = ""
    role_name  = ""
    role_type  = ""
    is_custom  = ""
    namespace  = ""
  }
}
`
