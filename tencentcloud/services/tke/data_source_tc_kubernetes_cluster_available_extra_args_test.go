package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesClusterAvailableExtraArgsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccKubernetesClusterAvailableExtraArgsDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_cluster_available_extra_args.example"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_kubernetes_cluster_available_extra_args.example", "available_extra_args.#"),
			),
		}},
	})
}

const testAccKubernetesClusterAvailableExtraArgsDataSource = `
data "tencentcloud_kubernetes_cluster_available_extra_args" "example" {
  cluster_version = "1.34.1"
  cluster_type    = "MANAGED_CLUSTER"
}
`
