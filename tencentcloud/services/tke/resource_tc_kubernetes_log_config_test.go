package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesLogConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccKubernetesLogConfig_cls,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_log_config.kubernetes_log_config", "id"),
			),
		}},
	})
}

const testAccKubernetesLogConfig_cls = `

resource "tencentcloud_kubernetes_log_config" "kubernetes_log_config" {
  log_config_name = "xxx"
  cluster_id = ""
  logset_id = ""
  log_config = ""
}
`
