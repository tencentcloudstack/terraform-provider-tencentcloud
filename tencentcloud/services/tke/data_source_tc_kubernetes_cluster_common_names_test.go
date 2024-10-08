package tke_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudKubernetesCommonNamesDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesCommonNamesBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_cluster_common_names.foo", "cluster_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_kubernetes_cluster_common_names.foo", "role_ids.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_cluster_common_names.foo", "list.#"),
				),
			},
		},
	})
}

// const KeepTkeCNRoleName = `
// variable "keep_tke_cn" {
//   default = "keep-for-tke-cn"
// }
// `

const testAccKubernetesCommonNamesBasic = testAccTkeCluster + `
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
  uin    = data.tencentcloud_user_info.info.uin
}

// data "tencentcloud_kubernetes_clusters" "cls" {
//   cluster_name = "` + tcacctest.DefaultTkeClusterName + `"
// }

data "tencentcloud_cam_roles" "role_basic" {
//   name          = var.keep_tke_cn
    name = "TKE_QCSRole"
}

data "tencentcloud_kubernetes_cluster_common_names" "foo" {
  cluster_id = tencentcloud_kubernetes_cluster.managed_cluster.id
  role_ids = [data.tencentcloud_cam_roles.role_basic.role_list.0.role_id]
}
`
