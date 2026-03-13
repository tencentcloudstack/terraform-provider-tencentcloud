package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMonitorExternalClusterRegisterCommandDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorExternalClusterRegisterCommandDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_external_cluster_register_command.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_external_cluster_register_command.example", "command"),
				),
			},
		},
	})
}

const testAccMonitorExternalClusterRegisterCommandDataSource = `
data "tencentcloud_monitor_external_cluster_register_command" "example" {
  instance_id = var.instance_id
  cluster_id  = var.cluster_id
}
`
