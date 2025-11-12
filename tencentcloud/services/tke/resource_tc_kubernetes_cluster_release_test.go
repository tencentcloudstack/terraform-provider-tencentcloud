package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesClusterReleaseResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterRelease,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_release.example", "id"),
				),
			},
		},
	})
}

const testAccKubernetesClusterRelease = `
resource "tencentcloud_kubernetes_cluster_release" "example" {
  cluster_id      = "cls-fdy7hm1q"
  name            = "tf-example"
  namespace       = "default"
  chart           = "nginx-ingress"
  chart_from      = "tke-market"
  chart_version   = "4.9.0"
  chart_namespace = "opensource-stable"
  cluster_type    = "tke"
  values {
    raw_original = <<-EOF
## nginx configuration\n##......ndhParam: \"\"\n
EOF
    values_type  = "yaml"
  }
}
`
