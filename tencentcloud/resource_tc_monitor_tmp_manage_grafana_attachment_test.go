package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorTmpManageGrafanaAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorTmpManageGrafanaAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_manage_grafana_attachment.tmp_manage_grafana_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_manage_grafana_attachment.tmp_manage_grafana_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorTmpManageGrafanaAttachment = `

resource "tencentcloud_monitor_tmp_manage_grafana_attachment" "tmp_manage_grafana_attachment" {
  instance_id = "prom-test"
  grafana_id = "grafana-test"
}

`
