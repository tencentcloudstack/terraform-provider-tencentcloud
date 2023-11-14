package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudClsShipperTasksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsShipperTasksDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cls_shipper_tasks.shipper_tasks")),
			},
		},
	})
}

const testAccClsShipperTasksDataSource = `

data "tencentcloud_cls_shipper_tasks" "shipper_tasks" {
  shipper_id = ""
  start_time = 
  end_time = 
  }

`
