package tmp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMonitorTmpMultipleWritesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorTmpMultipleWrites,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_multiple_writes.monitor_tmp_multiple_writes", "id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_multiple_writes.monitor_tmp_multiple_writes", "instance_id", "prom-l9cl1ptk"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_multiple_writes.monitor_tmp_multiple_writes", "remote_writes.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_multiple_writes.monitor_tmp_multiple_writes", "remote_writes.0.url", "http://172.16.0.111:9090/api/v1/prom/write"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_multiple_writes.monitor_tmp_multiple_writes", "remote_writes.0.url_relabel_config", "# 添加 label\n# - target_label: key\n#  replacement: value\n# 丢弃指标\n#- source_labels: [__name__]\n#  regex: kubelet_.+;\n#  action: drop"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_multiple_writes.monitor_tmp_multiple_writes",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorTmpMultipleWrites = `

resource "tencentcloud_monitor_tmp_multiple_writes" "monitor_tmp_multiple_writes" {
    instance_id = "prom-l9cl1ptk"

    remote_writes {
        label              = null
        max_block_size     = null
        url                = "http://172.16.0.111:9090/api/v1/prom/write"
        url_relabel_config = trimspace(<<-EOT
            # 添加 label
            # - target_label: key
            #  replacement: value
            # 丢弃指标
            #- source_labels: [__name__]
            #  regex: kubelet_.+;
            #  action: drop
        EOT
        )
    }
}
`
