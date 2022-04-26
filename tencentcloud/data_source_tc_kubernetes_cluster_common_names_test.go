package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudKubernetesCommonNames(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

const KeepTkeCNRoleName = `
variable "keep_tke_cn" {
  default = "keep-for-tke-cn"
}
`

const testAccKubernetesCommonNamesBasic = KeepTkeCNRoleName + `
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
  uin    = data.tencentcloud_user_info.info.uin
}

data "tencentcloud_kubernetes_clusters" "cls" {
  cluster_name = "` + defaultTkeClusterName + `"
}

data "tencentcloud_cam_roles" "role_basic" {
  name          = var.keep_tke_cn
}

data "tencentcloud_kubernetes_cluster_common_names" "foo" {
  cluster_id = data.tencentcloud_kubernetes_clusters.cls.list.0.cluster_id
  role_ids = [data.tencentcloud_cam_roles.role_basic.role_list.0.role_id]
}
`
