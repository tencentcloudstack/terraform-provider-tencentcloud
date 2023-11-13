package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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

data "tencentcloud_cls_machine_group_configs" "machine_group_configs" {
  group_id = ""
  }

`
