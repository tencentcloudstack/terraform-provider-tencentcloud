package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesUpgradeTaskDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccKubernetesUpgradeTaskDetailDataSource,
			Check:  resource.ComposeTestCheckFunc(resource.AccCheckTencentCloudDataSourceID("data.tencentcloud_kubernetes_upgrade_task_detail.kubernetes_upgrade_task_detail")),
		}},
	})
}

const testAccKubernetesUpgradeTaskDetailDataSource = `

data "tencentcloud_kubernetes_upgrade_task_detail" "kubernetes_upgrade_task_detail" {
}
`
