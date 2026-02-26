package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesCancelUpgradePlanOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccKubernetesCancelUpgradePlanOperation,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cancel_upgrade_plan_operation.kubernetes_cancel_upgrade_plan_operation", "id")),
		}, {
			ResourceName:      "tencentcloud_kubernetes_cancel_upgrade_plan_operation.kubernetes_cancel_upgrade_plan_operation",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccKubernetesCancelUpgradePlanOperation = `

resource "tencentcloud_kubernetes_cancel_upgrade_plan_operation" "kubernetes_cancel_upgrade_plan_operation" {
}
`
