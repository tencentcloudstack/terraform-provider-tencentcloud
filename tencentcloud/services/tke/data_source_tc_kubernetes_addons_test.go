package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesAddonsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesAddonDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_addons.kubernetes_addon"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_addons.kubernetes_addon", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_kubernetes_addons.kubernetes_addon", "addons.#", "8"),
				),
			},
		},
	})
}

const testAccKubernetesAddonDataSource = `

data "tencentcloud_kubernetes_addons" "kubernetes_addon" {
  cluster_id = "cls-fdy7hm1q"
}
`
