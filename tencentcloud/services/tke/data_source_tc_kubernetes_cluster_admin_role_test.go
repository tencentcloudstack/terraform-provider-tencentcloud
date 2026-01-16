package tke_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudKubernetesClusterAdminRoleDataSource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterAdminRoleDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_cluster_admin_role.admin_role"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_cluster_admin_role.admin_role", "cluster_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_cluster_admin_role.admin_role", "request_id"),
				),
			},
		},
	})
}

const testAccKubernetesClusterAdminRoleDataSource = `
variable "default_cluster_id" {
  default = "` + tcacctest.DefaultTkeClusterId + `"
}

data "tencentcloud_kubernetes_cluster_admin_role" "admin_role" {
  cluster_id = var.default_cluster_id
}
`
