package tmp_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMonitorTmpManageGrafanaAttachmentResource_basic -v
func TestAccTencentCloudMonitorTmpManageGrafanaAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorTmpManageGrafanaAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_manage_grafana_attachment.manage_grafana_attachment", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_manage_grafana_attachment.manage_grafana_attachment", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_manage_grafana_attachment.manage_grafana_attachment", "grafana_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_manage_grafana_attachment.manage_grafana_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testManageGrafanaAttachmentVar = `
variable "prometheus_id" {
  default = "` + tcacctest.DefaultPrometheusId + `"
}
variable "grafana_id" {
  default = "` + tcacctest.DefaultGrafanaInstanceId + `"
}
`

const testAccMonitorTmpManageGrafanaAttachment = testManageGrafanaAttachmentVar + `

resource "tencentcloud_monitor_tmp_manage_grafana_attachment" "manage_grafana_attachment" {
    grafana_id  = var.grafana_id
    instance_id = var.prometheus_id
}

`
