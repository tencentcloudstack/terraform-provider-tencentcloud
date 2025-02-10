package tmp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudNeedFixMonitorTmpMultipleWritesListResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorTmpMultipleWritesList,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_multiple_writes_list.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_multiple_writes_list.example", "instance_id"),
				),
			},
			{
				Config: testAccMonitorTmpMultipleWritesListUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_multiple_writes_list.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_multiple_writes_list.example", "instance_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_multiple_writes_list.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorTmpMultipleWritesList = `
resource "tencentcloud_monitor_tmp_multiple_writes_list" "example" {
  instance_id = "prom-gzg3f1em"

  remote_writes {
    url = "http://172.16.0.111:9090/api/v1/prom/write"
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
    headers {
      key   = "Key"
      value = "Value"
    }
  }

  remote_writes {
    url                = "http://172.16.0.111:8080/api/v1/prom/write"
    url_relabel_config = "# 添加 label\n#- target_label: key\n#  replacement: value\n# 丢弃指标\n#- source_labels: [__name__]\n#  regex: kubelet_.+;\n#  action: drop"
    headers {
      key   = "Key"
      value = "Value"
    }
  }
}
`

const testAccMonitorTmpMultipleWritesListUpdate = `
resource "tencentcloud_monitor_tmp_multiple_writes_list" "example" {
  instance_id = "prom-gzg3f1em"

  remote_writes {
    url = "http://172.16.0.111:9090/api/v1/prom/write"
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
    headers {
      key   = "Key"
      value = "Value"
    }
  }
}
`
