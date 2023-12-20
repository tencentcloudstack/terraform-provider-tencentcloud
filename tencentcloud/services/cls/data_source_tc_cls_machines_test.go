package cls_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClsMachinesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsMachinesDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cls_machines.machines")),
			},
		},
	})
}

const testAccClsMachinesDataSource = `

resource "tencentcloud_cls_machine_group" "group" {
  group_name        = "tf-describe-mg-test"
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

data "tencentcloud_cls_machines" "machines" {
  group_id = tencentcloud_cls_machine_group.group.id
}

`
