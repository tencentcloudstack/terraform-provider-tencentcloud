package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClsMachineGroupConfigsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsMachineGroupConfigsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cls_machine_group_configs.machine_group_configs")),
			},
		},
	})
}

const testAccClsMachineGroupConfigsDataSource = `

resource "tencentcloud_cls_machine_group" "group" {
  group_name        = "tf-describe-mg-config-test"
  service_logging   = true
  auto_update       = true
  update_end_time   = "19:05:00"
  update_start_time = "17:05:00"

  machine_group_type {
    type   = "ip"
    values = [
      "192.168.1.1",
      "192.168.1.2",
    ]
  }
}

data "tencentcloud_cls_machine_group_configs" "machine_group_configs" {
  group_id = tencentcloud_cls_machine_group.group.id
}

`
