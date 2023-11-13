package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorTmpExporterIntegrationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorTmpExporterIntegration,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_exporter_integration.tmp_exporter_integration", "id")),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_exporter_integration.tmp_exporter_integration",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorTmpExporterIntegration = `

resource "tencentcloud_monitor_tmp_exporter_integration" "tmp_exporter_integration" {
  instance_id = "prom-dko9d0nu"
  kind = "blackbox-exporter"
  content = "blackbox-exporter"
  kube_type = 1
  cluster_id = "job_name: demo-config"
}

`
